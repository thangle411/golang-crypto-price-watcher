package web2

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/thangle411/golang-web3-price-watcher/jsonHelper"
)

type Response struct {
	Symbol string
	Price  float64
	Volume int
}

type StockData struct {
	Symbol     string
	LastPrice  float64
	PriceDelta func(last float64, current float64) bool
}

type Info struct {
	Price     float64
	Symbol    string
	SendEmail bool
}

var Tickers = []StockData{
	{
		Symbol:    "COIN",
		LastPrice: 0,
		PriceDelta: func(last float64, current float64) bool {
			if last == 0 {
				return false
			}
			sendEmail := false
			percent := (current - last) / last * 100
			if last < 90 && current > 90 {
				sendEmail = true
			} else if last < 100 && current > 100 {
				sendEmail = true
			} else if last < 110 && current > 110 {
				sendEmail = true
			} else if last < 120 && current > 120 {
				sendEmail = true
			} else if percent > 1 {
				sendEmail = true
			}
			return sendEmail
		},
	},
}

func Start(receiverEmail string, senderEmail string, appPassword string) (html string, subject string) {
	currentHour := time.Now().Local().Hour()
	htmlString := ""
	subjectString := ""

	if currentHour >= 5 && currentHour < 13 {
		infoChannel := make(chan Info)
		var wg sync.WaitGroup
		for i := 0; i < len(Tickers); i++ {
			data := &Tickers[i]
			wg.Add(1)
			go getPrice(data, infoChannel, &wg)
		}

		go func() {
			wg.Wait()
			close(infoChannel)
		}()

		for data := range infoChannel {
			fmt.Println(data)
			if data.SendEmail {
				htmlString += fmt.Sprintf(`
			<br></br>
			<li>%s is $%f</li>
			`, data.Symbol, data.Price)
				subjectString += data.Symbol + "(s)" + " - "
			}
			fmt.Println("--------------------------------------")
		}
	}

	return htmlString, subjectString
}

func getPrice(data *StockData, infoChannel chan Info, wg *sync.WaitGroup) {
	defer wg.Done()

	godotenv.Load()
	apiKey := os.Getenv("FMP_API_KEY")
	if apiKey == "" {
		log.Fatal("FMP_API_KEY not found in .env")
	}

	url := "https://financialmodelingprep.com/api/v3/quote-short/" + data.Symbol + "?apikey=" + apiKey

	var response []Response
	err := jsonHelper.GetJson(url, &response)
	if err != nil {
		fmt.Println("Error getting data from FMP")
		response[0] = Response{
			Symbol: data.Symbol,
			Price:  0,
			Volume: 0,
		}
	}
	lastPrice := data.LastPrice
	data.LastPrice = response[0].Price

	infoChannel <- Info{
		Price:     response[0].Price,
		Symbol:    response[0].Symbol,
		SendEmail: data.PriceDelta(lastPrice, response[0].Price),
	}
}
