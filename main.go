package main

import (
	"context"
	"go-pizza-api/internal/deals"

	"github.com/aws/aws-lambda-go/lambda"
)

type PostCodeEvent struct {
	Postcode string `json:"postcode"`
}

func HandleRequest(ctx context.Context, postcode PostCodeEvent) ([]deals.AllDeals, error) {
	return deals.GetDeals(postcode.Postcode), nil
}

func main() {
	lambda.Start(HandleRequest)
}

// Easier for quick local testing
// func main() {
// 	deals := deals.GetDeals("me46ea")
// 	fmt.Printf("%+v", deals)
// }
