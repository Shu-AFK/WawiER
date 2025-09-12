package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Shu-AFK/WawiER/cmd/defines"
)

type Config struct {
	ApiBaseURL           string   `json:"ApiBaseURL"`
	ApiVersion           string   `json:"ApiVersion"`
	ExcludedOrderIdStart []string `json:"ExcludedOrderIdStart"`

	SmtpHost        string `json:"SmtpHost"`
	SmtpPort        string `json:"SmtpPort"`
	SmtpUsername    string `json:"SmtpUsername"`
	SmtpPassword    string `json:"SmtpPassword"`
	SmtpSenderEmail string `json:"SmtpSenderEmail"`

	LogMode string `json:"LogMode"`
	LogFile string `json:"LogFile"`
}

var Conf Config

func LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	requiredFields := map[string]string{
		"ApiBaseURL":      config.ApiBaseURL,
		"ApiVersion":      config.ApiVersion,
		"SmtpHost":        config.SmtpHost,
		"SmtpPort":        config.SmtpPort,
		"SmtpUsername":    config.SmtpUsername,
		"SmtpPassword":    config.SmtpPassword,
		"SmtpSenderEmail": config.SmtpSenderEmail,
	}

	for name, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("config validation error: %s must be set", name)
		}
	}

	switch config.LogMode {
	case "":
		config.LogMode = "console"
	case "none", "console":
	case "file", "both":
		if config.LogFile == "" {
			config.LogFile = "WawiER.log"
		}
		if absPath, err := filepath.Abs(config.LogFile); err == nil {
			config.LogFile = absPath
		}
	default:
		return errors.New("invalid LogMode: must be one of none, console, file, both")
	}

	Conf = config

	defines.APIBaseURL = config.ApiBaseURL
	defines.APIVersion = config.ApiVersion

	return nil
}
