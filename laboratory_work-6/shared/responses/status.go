package responses

type StatusResponse struct {
	Code    uint16 `json:"code"`
	Message string `json:"message"`
}
