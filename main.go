package main

import (
	"context"
	"go-pizza-api/internal/deals"

	"github.com/aws/aws-lambda-go/lambda"
)

type PostCodeEvent struct {
	Postcode string `json:"postcode"`
}

func HandleRequest(ctx context.Context, postcode PostCodeEvent) []deals.AllDeals {
	return deals.GetDeals(postcode.Postcode)
}

func main() {
	lambda.Start(HandleRequest)
}
