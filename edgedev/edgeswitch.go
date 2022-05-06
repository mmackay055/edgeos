package edgedev

import (
	"bytes"
	"crypto/tls"
        "github.com/mmackay055/edgeos/edgedev/edgedeverrors"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
        "fmt"
)

const (
	edgeswitch_path_scheme      string = `https`
	edgeswitch_path_login       string = `/api/v1.0/user/login`
	edgeswitch_path_save        string = ``
	edgeswitch_path_backup      string = `/api/v1.0/system/backup`
	edgeswitch_auth_cookie_name string = `x-auth-token`
)

type EdgeSwitch struct {
	url                 *url.URL
	TLSCertCheckDisable bool
	client              *http.Client
	xAuthToken          string
}

func (es *EdgeSwitch) login(host, username, password string) error {
	// create new client
	es.newClient()

	// create login credentials object
	creds, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}

	// build url
	es.url = &url.URL{
		Scheme: edgeswitch_path_scheme,
		Host:   host,
		Path:   edgeswitch_path_login,
	}

        req, err := http.NewRequest("POST", es.url.String(), bytes.NewReader(creds))
        if err != nil {
                return edgedeverrors.LoginError{Message: ERR_MSG_REQ_INVALID}
	}

        // setup request defaults
        es.SetupRequest(req)

	// send login request
	resp, err := es.client.Do(req)
	if err != nil {
                return edgedeverrors.LoginError{Message: err.Error()}
	}

	// unmarshall response
	esResp, err := EdgeSwitchRespParse(resp.Body)
	if err != nil {
                return edgedeverrors.LoginError{Message: ERR_MSG_RESP_PARSE + ": " + err.Error()}
	}

	// check status code
	if resp.StatusCode != http.StatusOK {
                return edgedeverrors.LoginError{Message: ERR_MSG_AUTH + ": " + esResp.Detail,}
	}

	// find x-auth-token header
	xAuthTokens, ok := resp.Header["X-Auth-Token"]
	if ok && len(xAuthTokens) > 0 {
		es.xAuthToken = xAuthTokens[0]
		return nil
	}
        return edgedeverrors.LoginError{Message: ERR_MSG_RESP_INVALID + ": missing x-auth-token header",}
}

func (es *EdgeSwitch) save() error {
	// TODO implement
        /*
        if  !es.isLoggedIn() {
                return edgedeverrors.SaveError{Message: ERR_MSG_NOTLOGGEDIN}
	}

	es.EndpointUpdate(edgeswitch_path_save)

	req, err := http.NewRequest("GET", es.Endpoint(), nil)
	if err != nil {
                return edgedeverrors.SaveError{Message: ERR_MSG_REQ_INVALID}
	}

	for _, c := range es.cookies {
		req.AddCookie(c)
	}

	resp, err := es.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusForbidden {
		return errors.New("save failed: authentication failure")
	} else if resp.StatusCode == http.StatusInternalServerError {
		return errors.New("save failed: device error")
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("save failed: unknown error")
	}


	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("save failed: failed to read response")
	}

	if err := checkSaveResp(data); err != nil {
		return err
	}
        */
        return nil
}

func (es *EdgeSwitch) backup() (*io.ReadCloser, string, error) {
        if  !es.isLoggedIn() {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_NOTLOGGEDIN}
	}

	// set path for backing up
	es.EndpointUpdate(edgeswitch_path_backup)

        // build request
	req, err := http.NewRequest("GET", es.Endpoint(), nil)
	if err != nil {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_REQ_INVALID}
	}

        // add request token
        es.SetupRequest(req)

	// make request
	resp, err := es.client.Do(req)
	if err != nil {
                return nil, "", edgedeverrors.BackupError{Message: err.Error()}
	}
        
	if resp.StatusCode == http.StatusForbidden {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_AUTH,}
	} else if resp.StatusCode == http.StatusInternalServerError {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_DEVERR,}
	} else if resp.StatusCode != http.StatusOK {
                return nil, "", edgedeverrors.BackupError{Message: ERR_MSG_UNKNOWN,}
	}

        // get file type from request
        fileType := GetFileType(resp)

	return &resp.Body, fileType, nil
}

func (es *EdgeSwitch) SetupRequest(req *http.Request) {
        if len(es.xAuthToken) > 0 {
                req.Header.Add("x-auth-token", es.xAuthToken)
        }
        req.Header.Add("Referer", fmt.Sprintf("%s://%s/", es.url.Scheme, es.url.Hostname()))
}

func (es *EdgeSwitch) disableTLSCertCheck(disable bool) {
	es.TLSCertCheckDisable = disable
}

func (es *EdgeSwitch) Endpoint() string {
	return es.url.String()
}

func (es *EdgeSwitch) EndpointUpdate(path string) {
	es.url.Path = path
}

func (es *EdgeSwitch) newClient() {
	var tr *http.Transport
	if es.TLSCertCheckDisable {
		tr = &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	} else {
		tr = &http.Transport{}
	}

	// make request
	es.client = &http.Client{
		Transport: tr,
		// don't follow redirects
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (es *EdgeSwitch) isLoggedIn() bool {
        return es.url != nil && len(es.xAuthToken) > 0
}
