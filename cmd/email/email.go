package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/Shu-AFK/WawiER/cmd/config"
	"github.com/Shu-AFK/WawiER/cmd/defines"
	"github.com/Shu-AFK/WawiER/cmd/logger"
)

type EConf struct {
	From     string
	Password string
	Host     string
	Port     string
	Username string
}

func LoadEmailConfig() (*EConf, error) {
	from := os.Getenv(defines.WawierEmailAddrEnv)
	pass := os.Getenv(defines.WawierEmailPassEnv)
	host := os.Getenv(defines.WawierEmailSMTPHostEnv)
	port := os.Getenv(defines.WawierSMTPPortEnv)
	user := os.Getenv(defines.WawierEmailSMTPUserEnv)

	// Fallback zu config.Conf falls ENV leer
	if from == "" {
		from = config.Conf.SmtpSenderEmail
	}
	if pass == "" {
		pass = config.Conf.SmtpPassword
	}
	if host == "" {
		host = config.Conf.SmtpHost
	}
	if port == "" {
		port = config.Conf.SmtpPort
	}
	if user == "" {
		user = config.Conf.SmtpUsername
	}

	if from == "" || pass == "" || host == "" || port == "" || user == "" {
		return nil, fmt.Errorf("missing SMTP configuration values (env vars or config)")
	}

	return &EConf{
		From:     from,
		Password: pass,
		Host:     host,
		Port:     port,
		Username: user,
	}, nil
}

func buildPlainTextBody(customerName, orderID string, items []string) string {
	itemList := "- " + strings.Join(items, "\r\n- ")

	return fmt.Sprintf(
		"Sehr geehrte/r %s,\r\n\r\n"+
			"wir möchten Sie darüber informieren, dass einige Artikel Ihrer Bestellung (Bestellnummer: %s) momentan nicht vorrätig sind:\r\n\r\n"+
			"%s\r\n\r\n"+
			"Damit Ihre Bestellung dennoch schnellstmöglich bearbeitet werden kann, bieten wir Ihnen folgende Optionen an:\r\n"+
			"- Auf die Lieferung der Artikel warten.\r\n"+
			"- Alternative Produkte auswählen.\r\n"+
			"- Die nicht verfügbaren Artikel stornieren und den Rest der Bestellung erhalten.\r\n\r\n"+
			"Bitte antworten Sie auf diese E-Mail, um uns Ihre Präferenz mitzuteilen.\r\n\r\n"+
			"Wir entschuldigen uns für die Unannehmlichkeiten und danken Ihnen für Ihre Geduld.\r\n\r\n"+
			"Mit freundlichen Grüßen,\r\nIhr Shop-Team",
		customerName, orderID, itemList,
	)
}

func buildHTMLBody(customerName, orderID string, items []string) string {
	var b strings.Builder
	b.WriteString("<ul>")
	for _, item := range items {
		b.WriteString(fmt.Sprintf("<li>%s</li>", item))
	}
	b.WriteString("</ul>")

	return fmt.Sprintf(
		"<!doctype html><html lang=\"de\"><head><meta charset=\"UTF-8\">"+
			"<style>"+
			"body{font-family:Arial,sans-serif;background:#f4f6f8;color:#333;padding:20px;}"+
			".container{background:#fff;max-width:600px;margin:auto;padding:20px 30px;border-radius:8px;box-shadow:0 2px 6px rgba(0,0,0,0.1);}"+
			"h2{color:#2a7ae2;margin-top:0;}"+
			"ul{padding-left:20px;margin:0;}"+
			"li{margin-bottom:6px;}"+
			".items{background:#f6f8fa;padding:12px;border-radius:6px;margin:15px 0;}"+
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
			"<p>Bitte antworten Sie auf diese E-Mail, um uns Ihre Präferenz mitzuteilen.</p>"+
			"<p>Wir entschuldigen uns für die Unannehmlichkeiten und danken Ihnen für Ihre Geduld.</p>"+
			"<p>Mit freundlichen Grüßen,<br>Ihr Shop-Team</p>"+
			"</div></body></html>",
		customerName, orderID, b.String(),
	)
}

func SendEmail(toAddress string, items []string, customerName, orderID string) {
	cfg, err := LoadEmailConfig()
	if err != nil {
		logger.Log.Printf("[ERROR] Email config not set: %v", err)
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
		logger.Log.Printf("[ERROR] Fehler beim Senden der Email an %s: %v", toAddress, err)
		return
	}

	logger.Log.Printf("[INFO] Email erfolgreich an %s gesendet", toAddress)
}
