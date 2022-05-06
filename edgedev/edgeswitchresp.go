package edgedev

import (
        "io"
        "encoding/json"
)

type EdgeSwitchResp struct {
        StatusCode int    `json:"statusCode"`
        Error      int    `json:"error"`
        Detail     string `json:"detail"`
        Message    string `json:"message"`
}
        
func EdgeSwitchRespParse(r io.Reader) (*EdgeSwitchResp, error) {
        data, err := io.ReadAll(r)
        if err != nil {
                return nil, err
        }

        resp := EdgeSwitchResp{}
        err = json.Unmarshal(data, &resp)
        if err != nil {
                return nil, err
        }

        return &resp, nil

}
