package email

import (
	"fmt"
	"log"
	"net/smtp"
)

const (
	SmtpHost = ""
	SmtpPort = ""
	Username = ""
	Password = ""
	From     = ""

	Subject = "Info zu ihrer Bestellung"
)

func SendEmail(emailAddress string, itemString string, customerName string, orderId string) {
	to := []string{emailAddress}

	body := fmt.Sprintf(
		"Sehr geehrte/r %s,\n\n"+
			"vielen Dank für Ihre Bestellung (Bestellnummer: %s).\n\n"+
			"Leider sind folgende Artikel momentan nicht sofort lieferbar, da ein Überverkauf stattgefunden hat:\n\n"+
			"%s\n"+
			"Wir werden Sie informieren, sobald die Artikel wieder verfügbar sind oder eine Teillieferung erfolgt.\n\n"+
			"Vielen Dank für Ihr Verständnis.\n\n"+
			"Mit freundlichen Grüßen,\nIhr Shop-Team",
		customerName,
		orderId,
		itemString,
	)

	auth := smtp.PlainAuth("", Username, Password, SmtpHost)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", Subject, body))

	err := smtp.SendMail(SmtpHost+":"+SmtpPort, auth, From, to, message)
	if err != nil {
		log.Printf("[ERROR] Fehler beim Senden der Email an %s: %v\n", emailAddress, err)
		return
	}

	log.Printf("[INFO] Email erfolgreich an %s gesendet\n", emailAddress)
}
