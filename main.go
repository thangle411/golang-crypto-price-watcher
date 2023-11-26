package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/joho/godotenv"
	"github.com/thangle411/golang-web3-price-watcher/coingecko"
	"github.com/thangle411/golang-web3-price-watcher/email"
	"github.com/thangle411/golang-web3-price-watcher/web3"
)

type Pool struct {
	poolAddress string
	tokenAddress string
}
type WatcherConfig struct {
	pools []Pool
}

var client *http.Client

func main() {
	var lastPrice float64
	receiverEmail, senderEmail, appPassword := setup()

	for {
		fmt.Printf("\n---------------------------\n")
		ethPrice := coingecko.GetTickerPrice("ethereum")
		if ethPrice == nil {
			log.Fatal("Failed getting price from coingecko")
		}

		fmt.Printf("\nEthereum price: %v\n", ethPrice["ethereum"].Usd)
		
		pool := web3.WatchPoolBalance(ethPrice["ethereum"].Usd, web3.PoolConfig{
			TokenAddress: "0x0B7f0e51Cd1739D6C96982D55aD8fA634dd43A9C",
			PoolAddress: "0xe3170D65018882a336743a9c396C52eA4B9c5563",
		})

		fmt.Printf("Last price %v, current price %v:\n", lastPrice, pool.Price)

		if priceDelta(lastPrice, pool.Price) {
			fmt.Println("Emailing...")
			err := email.SendEmail(pool.Price, pool, email.Email{
				AppEmail: senderEmail,
				AppPassword: appPassword,
				ToEmail: []string{receiverEmail},
			})
			if err != nil {
				fmt.Println("Failed sending email!", err)
			}
		}

		lastPrice = pool.Price
		time.Sleep(5 * time.Minute)
	}
}

func setup() (string, string, string) {
	var receiverEmail string
	var senderEmail string
	var appPassword string

	fmt.Println("\nUse environment variables? [y/n]: ")
	r := bufio.NewReader(os.Stdin)
	res, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	hasConfig := strings.ToLower(strings.TrimSpace(res))[0] == 'y'

	if hasConfig {
		godotenv.Load()
		receiverEmail = os.Getenv("GMAIL_RECEIVER_EMAIL")
		if receiverEmail == "" {
			log.Fatal("GMAIL_RECEIVER_EMAIL not found in .env")
		}
		appPassword = os.Getenv("GMAIL_APP_PASSWORD")
		if appPassword == "" {
			log.Fatal("GMAIL_APP_PASSWORD not found in .env")
		}
		senderEmail = os.Getenv("GMAIL_FROM_EMAIL")
		if senderEmail == "" {
			log.Fatal("GMAIL_FROM_EMAIL not found in .env")
		}
	} else {
		receiverEmail, senderEmail, appPassword = credentials()
	}

	errEmail := testEmail(receiverEmail, senderEmail, appPassword)
	if errEmail != nil {
		fmt.Println("\nCould not send test email, please double check your inputs")
		return setup()
	}

	return receiverEmail, senderEmail, appPassword
}

func credentials() (string, string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter the receiver email: ") 
	receiver, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("No receiver email provider")
	}

	fmt.Println("Enter the email to send from: ") 
	sender, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("No sender email provider")
	}

	fmt.Print("Enter app password for sender email: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("No app password provider")
	}

	password := string(bytePassword)
	return strings.TrimSpace(receiver), strings.TrimSpace(sender), strings.TrimSpace(password)
}

func testEmail(receiverEmail string, senderEmail string, appPassword string) error {
	mockPool := web3.PoolBalance{
		Eth: web3.Token{Balance: 0, Name: "Test email"},
		Token: web3.Token{Balance: 0, Name: "Test email"},
		Price: 0.00,
	}
	err := email.SendEmail(0.00, mockPool, email.Email{
		AppEmail: senderEmail,
		AppPassword: appPassword,
		ToEmail: []string{receiverEmail},
	})
	if err != nil {
		return err
	}
	return nil
}

func priceDelta(last float64, current float64) bool {
	if last == 0 {
		return false
	}

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