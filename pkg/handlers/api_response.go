package handlers

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

// ApiResponse is a function that creates an APIGatewayProxyResponse object with the given status code and body.
// It serializes the body to JSON and sets the appropriate headers.
func ApiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}, StatusCode: status}
	stringBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	resp.Body = string(stringBody)
	return &resp, nil 
}
