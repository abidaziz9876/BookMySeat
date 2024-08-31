package response

type ApiResponse struct {
	Status  int         `json:"status"`
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
