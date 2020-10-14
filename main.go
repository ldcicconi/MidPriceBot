package main

import (
	"errors"
	"fmt"
	"time"

	sfoxapi "github.com/ldcicconi/sfox-api-lib"
	"github.com/shopspring/decimal"
)

var one = decimal.NewFromFloat(1)
var OneHundred = decimal.NewFromFloat(100)
var buyThreshold = decimal.NewFromFloat(-0.03)
var sellThreshold = decimal.NewFromFloat(0)
var feeRate = decimal.NewFromFloat(0.0035)

func buy(price decimal.Decimal) {
	fmt.Println("Buying at price", price)
}

func sell(sellPrice, buyPrice decimal.Decimal) {
	profitRate := sellPrice.Sub(buyPrice).Div(buyPrice)
	// I got this formula by reducing ((1-feeRate)sellPrice - (1+feeRate)buyPrice )/ ((1+feeRate)buyPrice) in WolframAlpha
	// this takes into account having to pay a fee on each trade (the buy and the sell)
	profitRateWFees := sellPrice.Mul(one.Sub(feeRate)).Div(buyPrice.Mul(one.Add(feeRate))).Sub(one)
	fmt.Println("Selling at price", sellPrice, "bought at", buyPrice)

	onePlus := decimal.NewFromFloat(1).Add(profitRateWFees)
	a := decimal.NewFromFloat(0.01).Mul(onePlus)
	profitA := a.Mul(onePlus).Sub(a)
	profitB := one.Mul(onePlus).Sub(one)
	fmt.Printf("PROFIT STATS:\nProfit rate: %v%%\nProfit rate w/ Fees: %v%%\nProfit on 0.01BTC: %v\nProfit on 1BTC: %v\n\n", profitRate.Mul(OneHundred), profitRateWFees.Mul(OneHundred), profitA, profitB)
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
	var buyPrice decimal.Decimal

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
			buyPrice = ob.Asks[0].Price
			inPosition = true
		} else if percentDifference.GreaterThan(sellThreshold) && inPosition {
			fmt.Printf("MidPrice: %v WeightedMid: %v PercentDifference: %v\n", midPrice, weightedMidPrice, percentDifference)
			sell(ob.Bids[0].Price, buyPrice)
			inPosition = false
		}
		time.Sleep(time.Second)
	}
}
