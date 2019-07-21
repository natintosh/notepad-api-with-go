package utils

// SuccessMessage :
type SuccessMessage struct {
	Result string      `json:"result,omitempty"`
	Data   interface{} `json:"data"`
}

// ErrorMessage :
type ErrorMessage struct {
	Result string      `json:"result,omitempty"`
	Error  interface{} `json:"error"`
}
