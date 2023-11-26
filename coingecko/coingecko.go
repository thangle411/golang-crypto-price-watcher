package coingecko

import (
	"log"

	"github.com/thangle411/golang-web3-price-watcher/jsonHelper"
)

type Coin struct {
	Usd float64
}
type CoingeckoResponse map[string]Coin

func GetTickerPrice(ticker string) CoingeckoResponse {
	url :=  "https://api.coingecko.com/api/v3/simple/price?ids=" + ticker + "&vs_currencies=usd"

	var data CoingeckoResponse
	err := jsonHelper.GetJson(url, &data)
	if err != nil {
		log.Fatal("Error getting data from Coingecko")
	}

	return data
}