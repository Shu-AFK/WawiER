package email

import (
	"fmt"
	"log"
	"net/smtp"
)

const (
	SmtpHost = "smtp.ethereal.email"
	SmtpPort = "587"
	Username = "enoch.borer@ethereal.email"
	Password = "envmCkDS96Nzp8jjG2"
	From     = "enoch.borer@ethereal.email"

	Subject = "Info zu ihrer Bestellung"
)

func SendEmail(emailAddress string, itemString string, customerName string, orderId string) {
	to := []string{emailAddress}

	textBody := fmt.Sprintf(
		"Sehr geehrte/r %s,\r\n\r\n"+
			"vielen Dank für Ihre Bestellung (Bestellnummer: %s).\r\n\r\n"+
			"Leider sind folgende Artikel momentan nicht sofort lieferbar, da ein Überverkauf stattgefunden hat:\r\n\r\n"+
			"%s\r\n"+
			"Wir werden Sie informieren, sobald die Artikel wieder verfügbar sind oder eine Teillieferung erfolgt.\r\n\r\n"+
			"Vielen Dank für Ihr Verständnis.\r\n\r\n"+
			"Mit freundlichen Grüßen,\r\nIhr Shop-Team",
		customerName,
		orderId,
		itemString,
	)

	htmlBody := fmt.Sprintf(
		"<!doctype html><html><body style=\"font-family:Arial, sans-serif;\">\r\n"+
			"<p>Sehr geehrte/r %s,</p>\r\n"+
			"<p>vielen Dank für Ihre Bestellung (Bestellnummer: <strong>%s</strong>).</p>\r\n"+
			"<p>Leider sind folgende Artikel momentan nicht sofort lieferbar, da ein Überverkauf stattgefunden hat:</p>\r\n"+
			"<pre style=\"background:#f6f8fa;padding:12px;border-radius:6px;\">%s</pre>\r\n"+
			"<p>Wir werden Sie informieren, sobald die Artikel wieder verfügbar sind oder eine Teillieferung erfolgt.</p>\r\n"+
			"<p>Vielen Dank für Ihr Verständnis.</p>\r\n"+
			"<p>Mit freundlichen Grüßen,<br>Ihr Shop-Team</p>\r\n"+
			"</body></html>",
		customerName,
		orderId,
		itemString,
	)

	auth := smtp.PlainAuth("", Username, Password, SmtpHost)

	boundary := "boundary42- alt-part"
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/alternative; boundary=\"%s\"\r\n"+
			"\r\n"+
			"--%s\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
			"Content-Transfer-Encoding: 7bit\r\n\r\n"+
			"%s\r\n\r\n"+
			"--%s\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"Content-Transfer-Encoding: 7bit\r\n\r\n"+
			"%s\r\n\r\n"+
			"--%s--\r\n",
		From, emailAddress, Subject, boundary,
		boundary,
		textBody,
		boundary,
		htmlBody,
		boundary,
	))

	err := smtp.SendMail(SmtpHost+":"+SmtpPort, auth, From, to, message)
	if err != nil {
		log.Printf("[ERROR] Fehler beim Senden der Email an %s: %v\n", emailAddress, err)
		return
	}

	log.Printf("[INFO] Email erfolgreich an %s gesendet\n", emailAddress)
}
