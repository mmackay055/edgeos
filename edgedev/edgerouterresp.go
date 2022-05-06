package edgedev

type EdgeRouterSaveResp struct {
	Config  map[string]string `json:"CONFIG"`
	Success bool            `json:"success"`
}
