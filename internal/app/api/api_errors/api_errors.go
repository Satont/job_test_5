package api_errors

import "encoding/json"

func CreateBadRequestError(messages []string) []byte {
	response := make(map[string]any)
	response["messages"] = messages

	responseBody, _ := json.Marshal(response)
	return responseBody
}
