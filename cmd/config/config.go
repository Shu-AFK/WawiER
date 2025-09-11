package config

import (
	"encoding/json"
	"os"

	"github.com/Shu-AFK/WawiER/cmd/defines"
)

type Config struct {
	ApiBaseURL           string   `json:"ApiBaseURL"`
	ApiVersion           string   `json:"ApiVersion"`
	ExcludedOrderIdStart []string `json:"ExcludedOrderIdStart"`
	SmtpHost             string   `json:"SmtpHost"`
	SmtpPort             string   `json:"SmtpPort"`
	SmtpUsername         string   `json:"SmtpUsername"`
	SmtpPassword         string   `json:"SmtpPassword"`
	SmtpSenderEmail      string   `json:"SmtpSenderEmail"`
}

var Conf Config

func LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	Conf = config

	defines.APIBaseURL = config.ApiBaseURL
	defines.APIVersion = config.ApiVersion

	return nil
}
