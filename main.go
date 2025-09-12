package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Shu-AFK/WawiER/cmd/config"
	"github.com/Shu-AFK/WawiER/cmd/defines"
	"github.com/Shu-AFK/WawiER/cmd/logger"
	"github.com/Shu-AFK/WawiER/cmd/server"
	"github.com/Shu-AFK/WawiER/cmd/wawi/registration"
)

func runExclusionBatch(batPath string) error {
	cmd := exec.Command("powershell", "-Command",
		"Start-Process '"+batPath+"' -Verb runAs -WindowStyle Hidden")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	return cmd.Run()
}

func main() {
	// --- Parse CLI flags ---
	defaultPath := defines.ConfigPath
	cfgFlag := flag.String("config", defaultPath, "config file path")
	flag.Parse()

	cfgPath := *cfgFlag
	if cfgPath == "" {
		cfgPath = defaultPath
	}
	defines.ConfigPath = cfgPath

	if !strings.EqualFold(filepath.Ext(cfgPath), ".json") {
		log.Panic("-config must point to a .json file")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		if os.IsNotExist(err) {
			log.Panic("Config file not found: " + cfgPath)
		}
		log.Panic("Cannot access config file: " + err.Error())
	}

	if err := config.LoadConfig(cfgPath); err != nil {
		log.Panic("Failed to load config: " + err.Error())
	}

	// --- Init logger ---
	err := logger.InitLogger(config.Conf.LogMode, config.Conf.LogFile)
	if err != nil {
		log.Panic(err)
	}
	logger.Log.Printf("[INFO] Loaded config from %s", cfgPath)

	// --- Run Defender exclusion batch ---
	appFolder := filepath.Dir(os.Args[0])
	batPath := filepath.Join(appFolder, "exclude_from_defender.bat")

	logger.Log.Printf("[INFO] Running Defender exclusion batch (if needed)...")
	if err := runExclusionBatch(batPath); err != nil {
		logger.Log.Printf("[ERROR] Failed to run exclusion batch: %v", err)
	}

	// --- API key ---
	_, exists := os.LookupEnv(defines.APIKeyVarName)
	logger.Log.Printf("[INFO] Checking for API key in environment...")
	if !exists {
		logger.Log.Printf("[INFO] No API key found, registering...")
		apiKey, err := wawi_registration.Register()
		if err != nil {
			log.Panic(err)
		}

		if err := exec.Command("setx", defines.APIKeyVarName, apiKey).Run(); err != nil {
			log.Panic(err)
		}

		if err := os.Setenv(defines.APIKeyVarName, apiKey); err != nil {
			log.Panic(err)
		}
	} else {
		logger.Log.Printf("[INFO] API key found in environment")
	}

	// --- Context & Server ---
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go server.StartServer()

	<-ctx.Done()
	logger.Log.Printf("[INFO] Shutdown requested, cleaning up...")
}
