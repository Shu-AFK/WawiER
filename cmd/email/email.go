package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/Shu-AFK/WawiER/cmd/defines"
)

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
			"wir möchten Sie darüber informieren, dass einige Artikel Ihrer Bestellung (Bestellnummer: %s) momentan nicht vorrätig sind:\r\n\r\n"+
			"%s\r\n\r\n"+
			"Damit Ihre Bestellung dennoch schnellstmöglich bearbeitet werden kann, bieten wir Ihnen folgende Optionen an:\r\n"+
			"- Auf die Lieferung der Artikel warten.\r\n"+
			"- Alternative Produkte auswählen.\r\n"+
			"- Die nicht verfügbaren Artikel stornieren und den Rest der Bestellung erhalten.\r\n\r\n"+
			"Bitte antworten Sie auf diese E-Mail oder nutzen Sie Ihr Kundenkonto, um uns Ihre Präferenz mitzuteilen.\r\n\r\n"+
			"Wir entschuldigen uns für die Unannehmlichkeiten und danken Ihnen für Ihre Geduld.\r\n\r\n"+
			"Mit freundlichen Grüßen,\r\nIhr Shop-Team",
		customerName, orderID, items,
	)
}

func buildHTMLBody(customerName, orderID, items string) string {
	return fmt.Sprintf(
		"<!doctype html><html lang=\"de\"><head><meta charset=\"UTF-8\">"+
			"<style>"+
			"body{font-family:Arial,sans-serif;background:#f4f6f8;color:#333;padding:20px;}"+
			".container{background:#fff;max-width:600px;margin:auto;padding:20px 30px;border-radius:8px;box-shadow:0 2px 6px rgba(0,0,0,0.1);}"+
			"h2{color:#2a7ae2;margin-top:0;}"+
			"ul{padding-left:20px;}"+
			"li{margin-bottom:6px;}"+
			".items{background:#f6f8fa;padding:12px;border-radius:6px;font-family:monospace;white-space:pre-line;}"+
			"</style></head><body>"+
			"<div class=\"container\">"+
			"<h2>Wichtige Information zu Ihrer Bestellung</h2>"+
			"<p>Sehr geehrte/r %s,</p>"+
			"<p>wir möchten Sie darüber informieren, dass einige Artikel Ihrer Bestellung (Bestellnummer: <strong>%s</strong>) momentan nicht vorrätig sind:</p>"+
			"<div class=\"items\">%s</div>"+
			"<p>Damit Ihre Bestellung dennoch schnellstmöglich bearbeitet werden kann, bieten wir Ihnen folgende Optionen an:</p>"+
			"<ul>"+
			"<li>Auf die Lieferung der Artikel warten.</li>"+
			"<li>Alternative Produkte auswählen.</li>"+
			"<li>Die nicht verfügbaren Artikel stornieren und den Rest der Bestellung erhalten.</li>"+
			"</ul>"+
			"<p>Bitte antworten Sie auf diese E-Mail oder nutzen Sie Ihr Kundenkonto, um uns Ihre Präferenz mitzuteilen.</p>"+
			"<p>Wir entschuldigen uns für die Unannehmlichkeiten und danken Ihnen für Ihre Geduld.</p>"+
			"<p>Mit freundlichen Grüßen,<br>Ihr Shop-Team</p>"+
			"</div></body></html>",
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
		cfg.From, toAddress, fmt.Sprintf("Wichtige Informationen zu Ihrer Bestellung %s", orderID), boundary,
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
