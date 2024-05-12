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
		Railcard:    request.QueryStringParameters["railcard"],
	}

	if req.Destination == "" || req.Origin == "" || req.Departure == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Missing query parameters",
			StatusCode: 400,
		}, nil
	}

	resultsReturn, resultsSingle, err := scrapes.Scrape(req)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	var jsonReq []byte
	if req.Return == "" {
		jsonReq, _ = json.Marshal(resultsSingle)
	} else {
		jsonReq, _ = json.Marshal(resultsReturn)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonReq),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
