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
	appFolder := filepath.Dir(os.Args[0])
	batPath := filepath.Join(appFolder, "exclude_from_defender.bat")

	excluded, err := checkDefenderExclusion(appFolder)
	if err != nil {
		log.Printf("[ERROR] Failed to check Windows Defender exclusions: %v", err)
	} else if !excluded {
		log.Printf("[INFO] Folder not excluded. Running batch to add exclusion...")
		err = runExclusionBatch(batPath)
		if err != nil {
			log.Printf("[ERROR] Failed to run exclusion batch: %v", err)
		} else {
			log.Printf("[INFO] Exclusion batch executed successfully.")
		}
	} else {
		log.Printf("[INFO] Folder already excluded in Windows Defender.")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
	err = config.LoadConfig(cfgPath)
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

	go server.StartServer()

	<-ctx.Done()
	log.Printf("[INFO] Shutdown requested, cleaning up...")
}
