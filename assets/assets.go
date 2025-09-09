package assets

import (
	"embed"
	"encoding/base64"
)

//go:embed WawiER.png
var image embed.FS

func ImageToB64() string {
	data, err := image.ReadFile("WawiER.png")
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}
