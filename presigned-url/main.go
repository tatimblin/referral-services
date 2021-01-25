package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
		"github.com/aws/aws-sdk-go/service/s3"
		"github.com/aws/aws-lambda-go/events"
		"github.com/aws/aws-lambda-go/lambda"

		"time"
		"encoding/json"
)

type BodyRequest struct {
	Key string `json:"key"`
}

type BodyResponse struct {
	PresignedURL string `json:"presignedUrl"`
}

type Response events.APIGatewayProxyResponse

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	bodyRequest := BodyRequest{
		Key: "",
	}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return Response{
			Body: err.Error(),
			StatusCode: 404,
		}, err
	}

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String("referral-app-assets-1"),
			Key:    aws.String(bodyRequest.Key),
	})
	str, err := req.Presign(15 * time.Minute)

	bodyResponse := BodyResponse{
		PresignedURL: str,
	}

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
			"X-MyCompany-Func-Reply": "presigned-url-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
