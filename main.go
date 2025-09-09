package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Shu-AFK/WawiER/cmd/config"
	"github.com/Shu-AFK/WawiER/cmd/defines"
	"github.com/Shu-AFK/WawiER/cmd/server"
	"github.com/Shu-AFK/WawiER/cmd/wawi/registration"
)

func main() {
	defaultPath := defines.ConfigPath
	cfgFlag := flag.String("config", defaultPath, "config file path")
	flag.Parse()

	cfgPath := *cfgFlag
	if cfgPath == "" {
		cfgPath = defaultPath
	}
	defines.ConfigPath = cfgPath

	if !strings.EqualFold(filepath.Ext(cfgPath), ".json") {
		log.Printf("[ERROR] -config must point to a .json file (got %q)\n", cfgPath)
		os.Exit(2)
	}

	if _, err := os.Stat(cfgPath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("[ERROR] Config file not found: %s\n", cfgPath)
		} else {
			log.Printf("[ERROR] Cannot access config file %s: %v\n", cfgPath, err)
		}
		os.Exit(2)
	}

	log.Printf("[INFO] Loading config from %s...\n", cfgPath)
	err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Printf("[ERROR] Failed to load config: %v\n", err)
		os.Exit(1)
	}

	_, exists := os.LookupEnv(defines.APIKeyVarName)

	log.Printf("[INFO] Checking for API key in environment...\n")
	if !exists {
		log.Printf("[INFO] No API key found, registering...\n")
		apiKey, err := wawi_registration.Register()
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("setx", defines.APIKeyVarName, apiKey)
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		err = os.Setenv(defines.APIKeyVarName, apiKey)
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("[INFO] API key found in environment\n")
	}

	server.StartServer()
}
