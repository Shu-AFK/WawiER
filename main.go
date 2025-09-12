package main

import (
	"context"
	"flag"
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

func checkDefenderExclusion(folderPath string) (bool, error) {
	psCmd := `Get-MpPreference | Select-Object -ExpandProperty ExclusionPath`
	cmd := exec.Command("powershell", "-Command", psCmd)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.TrimSpace(strings.ToLower(line)) == strings.ToLower(folderPath) {
			return true, nil
		}
	}
	return false, nil
}

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

	// --- Config loading ---
	cfgPath := *cfgFlag
	if cfgPath == "" {
		cfgPath = defaultPath
	}
	defines.ConfigPath = cfgPath

	if !strings.EqualFold(filepath.Ext(cfgPath), ".json") {
		panic("-config must point to a .json file")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		if os.IsNotExist(err) {
			panic("Config file not found: " + cfgPath)
		}
		panic("Cannot access config file: " + err.Error())
	}

	if err := config.LoadConfig(cfgPath); err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// --- Init logger from config ---
	err := logger.InitLogger(config.Conf.LogMode, config.Conf.LogFile)
	if err != nil {
		panic(err)
	}

	logger.Log.Printf("[INFO] Loaded config from %s", cfgPath)

	// --- Defender exclusion check ---
	appFolder := filepath.Dir(os.Args[0])
	batPath := filepath.Join(appFolder, "exclude_from_defender.bat")

	excluded, err := checkDefenderExclusion(appFolder)
	if err != nil {
		logger.Log.Printf("[ERROR] Failed to check Windows Defender exclusions: %v", err)
	} else if !excluded {
		logger.Log.Printf("[INFO] Folder not excluded. Running batch to add exclusion...")
		err = runExclusionBatch(batPath)
		if err != nil {
			logger.Log.Printf("[ERROR] Failed to run exclusion batch: %v", err)
			panic(err)
		} else {
			logger.Log.Printf("[INFO] Exclusion batch executed successfully.")
		}
	} else {
		logger.Log.Printf("[INFO] Folder already excluded in Windows Defender.")
	}

	// --- API key ---
	_, exists := os.LookupEnv(defines.APIKeyVarName)
	logger.Log.Printf("[INFO] Checking for API key in environment...")
	if !exists {
		logger.Log.Printf("[INFO] No API key found, registering...")
		apiKey, err := wawi_registration.Register()
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("setx", defines.APIKeyVarName, apiKey)
		if err := cmd.Run(); err != nil {
			panic(err)
		}

		if err := os.Setenv(defines.APIKeyVarName, apiKey); err != nil {
			panic(err)
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
