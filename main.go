package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/thangle411/golang-web3-price-watcher/email"
	"github.com/thangle411/golang-web3-price-watcher/jsonHelper"
	"github.com/thangle411/golang-web3-price-watcher/web3"

	"github.com/joho/godotenv"
)


type Coin struct {
	Usd float64
}
type CoingeckoResponse map[string]Coin

var client *http.Client

func main() {
	var lastPrice float64

	for {
		fmt.Printf("\n---------------------------\n")
		ethPrice := getTickerPrice("ethereum")
		if ethPrice == nil {
			log.Fatal("Failed getting price from coingecko")
		}

		fmt.Printf("\nEthereum price: %v\n", ethPrice["ethereum"].Usd)
		
		pool := web3.WatchPoolBalance(ethPrice["ethereum"].Usd, web3.PoolConfig{
			TokenAddress: "0x0B7f0e51Cd1739D6C96982D55aD8fA634dd43A9C",
			PoolAddress: "0xe3170D65018882a336743a9c396C52eA4B9c5563",
		})

		fmt.Println("Price:", pool.Price)

		if priceDelta(lastPrice, pool.Price) {
			fmt.Println("Emailing...")
			godotenv.Load()
			appPassword := os.Getenv("GMAIL_APP_PASSWORD")
			if appPassword == "" {
				log.Fatal("GMAIL_APP_PASSWORD not found in .env")
			}
			appEmail := os.Getenv("GMAIL_FROM_EMAIL")
			if appPassword == "" {
				log.Fatal("GMAIL_FROM_EMAIL not found in .env")
			}
			err := email.SendEmail(pool.Price, pool, email.Email{
				AppEmail: appEmail,
				AppPassword: appPassword,
				ToEmail: []string{"leet0822@gmail.com"},
			})
			if err != nil {
				fmt.Println("Failed sending email!", err)
			}
		}

		lastPrice = pool.Price
		time.Sleep(5 * time.Minute)
	}
}

func priceDelta(last float64, current float64) bool {
	sendEmail := false
	percent := (current - last) / last * 100
	if last < 15 && current > 15 {
		sendEmail = true
	}  else if last < 25 && current > 25 {
		sendEmail = true
	}  else if last < 35 && current > 35 {
		sendEmail = true
	} else if percent > 5 {
		sendEmail = true
	}
	return sendEmail
}

func getTickerPrice(ticker string) CoingeckoResponse {
	godotenv.Load()
	cgURL := os.Getenv("COINGECKO_URL")
	if cgURL == "" {
		log.Fatal("Coingecko URL not found in .env")
	}

	url := cgURL + "/simple/price?ids=" + ticker + "&vs_currencies=usd"

	var data CoingeckoResponse
	err := jsonHelper.GetJson(url, &data)
	if err != nil {
		log.Fatal("Error getting data from Coingecko")
	}

	return data
}