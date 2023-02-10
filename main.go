package main

import (
	"encoding/json"
	"fmt"
	"go-pizza-api/internal/deals"
	"log"
)

// type PostCodeEvent struct {
// 	Postcode string `json:"postcode"`
// }

// func HandleRequest(ctx context.Context, postcode PostCodeEvent) ([]deals.AllDeals, error) {
// 	return deals.GetDeals(postcode.Postcode), nil
// }

// func main() {
// 	lambda.Start(HandleRequest)
// }

// Easier for quick local testing
func main() {
	deals := deals.GetDeals("me46ea")
	json, err := json.Marshal(deals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", string(json))
}
