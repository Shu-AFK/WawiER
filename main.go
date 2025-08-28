package main

import (
	"os"
	"os/exec"

	"github.com/Shu-AFK/WawiER/cmd/defines"
	wawi_registration "github.com/Shu-AFK/WawiER/cmd/wawi/registration"
)

func main() {
	_, exists := os.LookupEnv(defines.APIKeyVarName)

	if !exists {
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
	}
}
