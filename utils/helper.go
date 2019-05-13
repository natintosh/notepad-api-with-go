package utils

import "reflect"

// GetStatusCode :
func GetStatusCode(result, errorMessage interface{}) (bool, int) {
	exist := false
	statusCode := 0
	if reflect.TypeOf(result) == reflect.TypeOf(ErrorMessage{}) {
		exist = true
		message := result.(ErrorMessage)
		messageErrorInterface := message.Error
		messageError := messageErrorInterface.(map[string]interface{})
		statusCode = messageError["code"].(int)
	}

	return exist, statusCode
}
