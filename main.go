package main

import (
	"fmt"
	"go-pizza-api/internal/deals"
)

func main() {
	deals := deals.GetDeals("ME46EA")
	fmt.Print(deals)
}
