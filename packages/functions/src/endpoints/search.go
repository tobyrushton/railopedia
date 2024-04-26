package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tobyrushton/railopedia/packages/functions/src/scrapes"
)

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	req := scrapes.Request{
		Destination: request.QueryStringParameters["destination"],
		Origin:      request.QueryStringParameters["origin"],
		Departure:   request.QueryStringParameters["departure"],
		Return:      request.QueryStringParameters["return"],
	}

	if req.Destination == "" || req.Origin == "" || req.Departure == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Missing query parameters",
			StatusCode: 400,
		}, nil
	}

	results, err := scrapes.Scrape(req)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	jsonReq, _ := json.Marshal(results)

	return events.APIGatewayProxyResponse{
		Body:       string(jsonReq),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
