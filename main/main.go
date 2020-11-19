package main

import (
	//"errors"
	"fmt"
	//"io/ioutil"
	//"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf(request.Body)

	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("69"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
