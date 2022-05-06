package edgedev

import (
	"crypto/tls"
	"encoding/json"
        "github.com/mmackay055/edgeos/edgedev/edgedeverrors"
	"io"
	"net/http"
	"net/url"
)

const (
	edgerouter_path_scheme string = `https`
	edgerouter_path_login  string = `/`
	edgerouter_path_save   string = `/api/edge/config/save.json`
	edgerouter_path_backup string = `/files/config/`
)

type EdgeRouter struct {
	cookies             []*http.Cookie
	url                 *url.URL
	TLSCertCheckDisable bool
	client              *http.Client
}

func (er *EdgeRouter) login(host, username, password string) error {
	// create new client
	er.newClient()

	// set login creds
	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)

	// build url
	er.url = &url.URL{
		Scheme: edgerouter_path_scheme,
		Host:   host,
		Path:   edgerouter_path_login,
	}

	resp, err := er.client.PostForm(er.Endpoint(), params)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusSeeOther {
                return edgedeverrors.LoginError{Message:ERR_MSG_AUTH,}
	}

	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == "PHPSESSID" {
			er.cookies = cookies
			return nil
		}
	}

        return edgedeverrors.LoginError{Message: ERR_MSG_RESP_INVALID,}
}

func (er *EdgeRouter) save() error {
	if er.url == nil || er.cookies == nil {
                return edgedeverrors.SaveError{Message: ERR_MSG_NOTLOGGEDIN,}
	}

	er.EndpointUpdate(edgerouter_path_save)

	req, err := http.NewRequest("GET", er.Endpoint(), nil)
	if err != nil {
                return edgedeverrors.SaveError{Message: ERR_MSG_REQ_INVALID}
	}

        er.SetupRequest(req)

	resp, err := er.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusForbidden {
                return edgedeverrors.SaveError{Message: ERR_MSG_AUTH,}
	} else if resp.StatusCode == http.StatusInternalServerError {
                return edgedeverrors.SaveError{Message: ERR_MSG_DEVERR,}
	} else if resp.StatusCode != http.StatusOK {
                return edgedeverrors.SaveError{Message: ERR_MSG_UNKNOWN,}
	}


	data, err := io.ReadAll(resp.Body)
	if err != nil {
                return edgedeverrors.SaveError{Message: ERR_MSG_RESP_READ}
	}

	if err := checkSaveResp(data); err != nil {
		return err
	}

	return nil
}

func (er *EdgeRouter) backup() (*io.ReadCloser, string, error) {

	if er.url == nil || er.cookies == nil {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_NOTLOGGEDIN}
	}

	// set path for backing up
	er.EndpointUpdate(edgerouter_path_backup)

        // build request
	req, err := http.NewRequest("GET", er.Endpoint(), nil)
	if err != nil {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_REQ_INVALID}
	}

        er.SetupRequest(req)	

	// make request
	resp, err := er.client.Do(req)
	if err != nil {
		return nil, "", err
	}
        
	if resp.StatusCode == http.StatusForbidden {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_AUTH}
	} else if resp.StatusCode == http.StatusInternalServerError {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_DEVERR}
	} else if resp.StatusCode != http.StatusOK {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_UNKNOWN}
	}


        // get file type from request
        fileType := GetFileType(resp)

	return &resp.Body, fileType, nil
}

func (er *EdgeRouter) disableTLSCertCheck(disable bool) {
        er.TLSCertCheckDisable = disable
}

func checkSaveResp(data []byte) error {

	var saveResp EdgeRouterSaveResp
	if err := json.Unmarshal(data, &saveResp); err != nil {
                return edgedeverrors.SaveError{Message:ERR_MSG_RESP_PARSE,}
	}

	if !saveResp.Success {
                return edgedeverrors.SaveError{Message: "saving config failed"}
	}

	return nil
}

func (er *EdgeRouter) SetupRequest(req *http.Request) {
        // add cookies
	for _, c := range er.cookies {
		req.AddCookie(c) // add sess cookie
	}
}


func (er *EdgeRouter) EndpointUpdate(path string) {
	er.url.Path = path
}

func (er *EdgeRouter) Endpoint() string {
	return er.url.String()
}

func (er *EdgeRouter) newClient() {
	var tr *http.Transport
	if er.TLSCertCheckDisable {
		tr = &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	} else {
		tr = &http.Transport{}
	}

	// make request
	er.client = &http.Client{
		Transport: tr,
		// don't follow redirects
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
