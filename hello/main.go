package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type BodyRequest struct {
	RequestName string `json:"name"`
}

type BodyResponse struct {
	ResponseName string `json:"name"`
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {

	bodyRequest := BodyRequest{
		RequestName: "",
	}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return Response{
			Body: err.Error(),
			StatusCode: 404,
		}, err
	}

	bodyResponse := BodyResponse{
		ResponseName: bodyRequest.RequestName + " LastName",
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return Response{
			Body: err.Error(),
			StatusCode: 404,
		}, nil
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(response),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
