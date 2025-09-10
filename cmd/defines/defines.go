package defines

import "github.com/Shu-AFK/WawiER/assets"

const (
	AppID                  = "WawiER/v1"
	DisplayName            = "WawiER"
	Description            = "Automatisierung von Email versedungen bei negativem Warenbestand"
	Version                = "1.0.0"
	ProviderName           = "Floyd Göttsch"
	ProviderWebsite        = "https://www.alpa-industrievertretungen.de/"
	XChallangeCode         = "wh5x1kgdm2koqsc311apskr45ksd"
	APIKeyVarName          = "WAWIER_APIKEY"
	WawierEmailAddrEnv     = "WAWIER_EMAIL_ADDR"
	WawierEmailPassEnv     = "WAWIER_EMAIL_PASS"
	WawierEmailSMTPHostEnv = "WAWIER_EMAIL_SMTP_HOST"
	WawierSMTPPortEnv      = "WAWIER_SMTP_PORT"
	WawierEmailSMTPUserEnv = "WAWIER_EMAIL_SMTP_USER"

	ServerApiKey = "Bearer c4b55569-3d82-44a0-b9e1-79a06b79eaf1"
)

var (
	APIVersion = "1.0"
	APIBaseURL = "http://127.0.0.1:5883/api/eazybusiness/"
	ConfigPath = "config.json"
)

var MandatoryAPIScope = []string{
	"salesorder.querysalesorders",
	"salesorder.querysalesorderlineitems",
	"stock.querystocksperitem",
}

type AppData struct {
	AppId              string   `json:"AppId"`
	DisplayName        string   `json:"DisplayName"`
	Description        string   `json:"Description"`
	Version            string   `json:"Version"`
	ProviderName       string   `json:"ProviderName"`
	ProviderWebsite    string   `json:"ProviderWebsite"`
	MandatoryApiScopes []string `json:"MandatoryApiScopes"`
	AppIcon            string   `json:"AppIcon"`
}

type RegistrationResponse struct {
	AppID                 string `json:"AppId"`
	RegistrationRequestId string `json:"RegistrationRequestId"`
	Status                int    `json:"Status"`
}

type FetchRegistrationResponse struct {
	RequestStatusInfo RequestStatus `json:"RequestStatusInfo"`
	Token             TokenInfo     `json:"Token"`
	GrantedScopes     []string      `json:"GrantedScopes"`
}

type RequestStatus struct {
	AppId                 string `json:"AppId"`
	RegistrationRequestId string `json:"RegistrationRequestId"`
	Status                int    `json:"Status"`
}

type TokenInfo struct {
	ApiKey string `json:"ApiKey"`
}

func ConstructAppData() *AppData {
	return &AppData{
		AppId:              AppID,
		DisplayName:        DisplayName,
		Description:        Description,
		Version:            Version,
		ProviderName:       ProviderName,
		ProviderWebsite:    ProviderWebsite,
		MandatoryApiScopes: MandatoryAPIScope,
		AppIcon:            assets.ImageToB64(),
	}
}
