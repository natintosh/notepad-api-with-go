package utils

// SuccessMessage :
type SuccessMessage struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
}

// ErrorMessage :
type ErrorMessage struct {
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
}
