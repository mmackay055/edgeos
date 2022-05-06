package edgedev

import (
	"bufio"
	"io"
	"net/http"
	"os"
)

type EdgeDev interface {
	// login requests to login to device
	login(string, string, string) error

	save() error

	// backup requests device to respond with a backup of the configuration
	// saved on the device, and file type
	backup() (*io.ReadCloser, string, error)

        newClient()

	disableTLSCertCheck(bool)

        SetupRequest(*http.Request)
}

const (
        ERR_MSG_AUTH string = "authentication invalid"
        ERR_MSG_RESP_INVALID string = "invalid response from device"
        ERR_MSG_RESP_PARSE string = "failed to parse response"
        ERR_MSG_REQ_INVALID string = "invalid request"
        ERR_MSG_NOTLOGGEDIN string = "not logged in"
        ERR_MSG_DEVERR string = "device error"
        ERR_MSG_UNKNOWN string = "unknown error"
        ERR_MSG_RESP_READ string = "failed to read response"
)

func Login(ed EdgeDev, host string, username string, password string) error {
	if err := ed.login(host, username, password); err != nil {
		return err
	}

	return nil
}

func Save(ed EdgeDev) error {
	if err := ed.save(); err != nil {
                return err
	}
	return nil
}

// BackUp requests for Edge device to back itself up
func BackUp(ed EdgeDev, destPath string) error {

	//must save device first
	if err := Save(ed); err != nil {
                return err
	}

	respBody, fileType, err := ed.backup()
	if err != nil {
		return err
	}

	bufread := bufio.NewReader(*respBody)

	//write to file if specified or stdout if not
	var destination *os.File
	if len(destPath) > 0 {
                var path string
                if len(fileType) > 0 {
                        path = destPath + "." + fileType
                } else {
                        path = destPath
                }

		dest, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		destination = dest
	} else {
		destination = os.Stdout
	}

	//write to file
	bufread.WriteTo(destination)

	// close response from device
	(*respBody).Close()

	return nil
}

func DisableTLSCertCheck(ed EdgeDev, disable bool) {
	ed.disableTLSCertCheck(disable)
}

func GetFileType(resp *http.Response) string {
        switch resp.Header.Get("content-type") {
        case "application/gzip":
                return "tar.gz"
        default:
                return ""
        }
}
