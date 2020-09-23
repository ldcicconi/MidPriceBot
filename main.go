package main

import (
	"errors"
	"fmt"
	"time"

	sfoxapi "github.com/ldcicconi/sfox-api-lib"
	"github.com/shopspring/decimal"
)

var OneHundred = decimal.NewFromFloat(100)
var buyThreshold = decimal.NewFromFloat(-0.03)
var sellThreshold = decimal.NewFromFloat(0)

func buy(price decimal.Decimal) {
	fmt.Println("buying at price", price)
}

func sell(price decimal.Decimal) {
	fmt.Println("selling at price", price)
}

func percentDifference(v1, v2 decimal.Decimal) (decimal.Decimal, error) {
	if v1.IsZero() {
		return decimal.Zero, errors.New("v1 is zero")
	}

	return v2.Sub(v1).Div(v1).Mul(OneHundred), nil
}

func main() {
	// Program starts up
	fmt.Println("STARTING BOT")

	sfox := sfoxapi.NewSFOXAPI("")

	var inPosition bool

	for {
		// Get the weighted mid price and true mid price
		ob, err := sfox.GetOrderbook("btcusd")
		if err != nil {
			fmt.Println("ERROR GETTING ORDERBOOK:", err)
			continue
		}

		midPrice, _ := ob.MidPrice()
		weightedMidPrice, _ := ob.WeightedMidPriceSimple()

		percentDifference, _ := percentDifference(midPrice, weightedMidPrice)
		// fmt.Printf("MidPrice: %v WeightedMid: %v PercentDifference: %v\n", midPrice, weightedMidPrice, percentDifference)

		if percentDifference.LessThan(buyThreshold) && !inPosition {
			fmt.Printf("MidPrice: %v WeightedMid: %v PercentDifference: %v\n", midPrice, weightedMidPrice, percentDifference)
			buy(ob.Asks[0].Price)
			inPosition = true
		} else if percentDifference.GreaterThan(sellThreshold) && inPosition {
			fmt.Printf("MidPrice: %v WeightedMid: %v PercentDifference: %v\n", midPrice, weightedMidPrice, percentDifference)
			sell(ob.Bids[0].Price)
			inPosition = false
		}
		time.Sleep(time.Second)
	}
}
