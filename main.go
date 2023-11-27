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
	"github.com/thangle411/golang-web3-price-watcher/email"
	"github.com/thangle411/golang-web3-price-watcher/web3"
)

var client *http.Client

func main() {
	receiverEmail, senderEmail, appPassword := setup()
	for {
		fmt.Println("------------RUNNING------------")
		web3.Start(receiverEmail, senderEmail, appPassword)
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

	errEmail := email.TestEmail(receiverEmail, senderEmail, appPassword)
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
