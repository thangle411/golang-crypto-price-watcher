package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/thangle411/golang-web3-price-watcher/web3"
)

type Email struct {
	AppEmail string;
	AppPassword string;
	ToEmail []string
}

func SendEmail(tokenPrice float64, pool web3.PoolBalance, emailConfig Email) error {
	e := email.NewEmail()
	e.From = "Price Tracker <" + emailConfig.AppEmail + ">"
	e.To = emailConfig.ToEmail
	e.Subject = pool.Token.Name + " price changed"
	e.HTML = []byte(fmt.Sprintf(`
	<div>%s is $%f</div>
	<div>There is %f %s and %f %s in the pool</div>
	`, pool.Token.Name, tokenPrice, pool.Eth.Balance, pool.Eth.Name, pool.Token.Balance, pool.Token.Name))
	return e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "golanglearner411@gmail.com", emailConfig.AppPassword, "smtp.gmail.com"))
}