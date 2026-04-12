package responses

type SuccessResponse[T any] struct {
	Code    uint16 `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
