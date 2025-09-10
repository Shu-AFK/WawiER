package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/Shu-AFK/WawiER/cmd/defines"
)

const Subject = "Info zu Ihrer Bestellung"

type EmailConfig struct {
	From     string
	Password string
	Host     string
	Port     string
	Username string
}

func LoadEmailConfig() (*EmailConfig, error) {
	requiredVars := map[string]string{
		"From":     os.Getenv(defines.WawierEmailAddrEnv),
		"Password": os.Getenv(defines.WawierEmailPassEnv),
		"Host":     os.Getenv(defines.WawierEmailSMTPHostEnv),
		"Port":     os.Getenv(defines.WawierSMTPPortEnv),
		"Username": os.Getenv(defines.WawierEmailSMTPUserEnv),
	}

	for key, val := range requiredVars {
		if val == "" {
			return nil, fmt.Errorf("missing environment variable for %s", key)
		}
	}

	return &EmailConfig{
		From:     requiredVars["From"],
		Password: requiredVars["Password"],
		Host:     requiredVars["Host"],
		Port:     requiredVars["Port"],
		Username: requiredVars["Username"],
	}, nil
}

func buildPlainTextBody(customerName, orderID, items string) string {
	return fmt.Sprintf(
		"Sehr geehrte/r %s,\r\n\r\n"+
			"vielen Dank für Ihre Bestellung (Bestellnummer: %s).\r\n\r\n"+
			"Leider sind folgende Artikel momentan nicht sofort lieferbar, da ein Überverkauf stattgefunden hat:\r\n\r\n"+
			"%s\r\n"+
			"Wir werden Sie informieren, sobald die Artikel wieder verfügbar sind oder eine Teillieferung erfolgt.\r\n\r\n"+
			"Vielen Dank für Ihr Verständnis.\r\n\r\n"+
			"Mit freundlichen Grüßen,\r\nIhr Shop-Team",
		customerName, orderID, items,
	)
}

func buildHTMLBody(customerName, orderID, items string) string {
	return fmt.Sprintf(
		"<!doctype html><html><body style=\"font-family:Arial, sans-serif;\">\r\n"+
			"<p>Sehr geehrte/r %s,</p>\r\n"+
			"<p>vielen Dank für Ihre Bestellung (Bestellnummer: <strong>%s</strong>).</p>\r\n"+
			"<p>Leider sind folgende Artikel momentan nicht sofort lieferbar, da ein Überverkauf stattgefunden hat:</p>\r\n"+
			"<pre style=\"background:#f6f8fa;padding:12px;border-radius:6px;\">%s</pre>\r\n"+
			"<p>Wir werden Sie informieren, sobald die Artikel wieder verfügbar sind oder eine Teillieferung erfolgt.</p>\r\n"+
			"<p>Vielen Dank für Ihr Verständnis.</p>\r\n"+
			"<p>Mit freundlichen Grüßen,<br>Ihr Shop-Team</p>\r\n"+
			"</body></html>",
		customerName, orderID, items,
	)
}

func SendEmail(toAddress, items, customerName, orderID string) {
	cfg, err := LoadEmailConfig()
	if err != nil {
		log.Printf("[ERROR] Email config not set: %v", err)
		return
	}

	to := []string{toAddress}
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	textBody := buildPlainTextBody(customerName, orderID, items)
	htmlBody := buildHTMLBody(customerName, orderID, items)

	boundary := "boundary42-alt-part"
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
		cfg.From, toAddress, Subject, boundary,
		boundary,
		textBody,
		boundary,
		htmlBody,
		boundary,
	))

	err = smtp.SendMail(cfg.Host+":"+cfg.Port, auth, cfg.From, to, message)
	if err != nil {
		log.Printf("[ERROR] Fehler beim Senden der Email an %s: %v", toAddress, err)
		return
	}

	log.Printf("[INFO] Email erfolgreich an %s gesendet", toAddress)
}
