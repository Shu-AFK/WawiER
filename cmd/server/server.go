package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Shu-AFK/WawiER/cmd/defines"
	"github.com/Shu-AFK/WawiER/cmd/logger"
	"github.com/Shu-AFK/WawiER/cmd/structs"
	"github.com/Shu-AFK/WawiER/cmd/wawi"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authentication")
	if authHeader != defines.ServerApiKey {
		logger.Log.Printf("[ERROR] apiHandler: Unothorized Request: %s\n", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Printf("[ERROR] apiHandler: %v\n", err)
		http.Error(w, "Fehler beim Lesen des Bodys", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var jsonBody structs.OrderReq
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		logger.Log.Printf("[ERROR] apiHandler: %v\n", err)
		http.Error(w, "Fehler beim Parsen des Bodys", http.StatusInternalServerError)
		return
	}

	logger.Log.Printf("[INFO] Handling orderId: %s\n", jsonBody.OrderId)
	err = wawi.HandleOrderId(jsonBody)
	if err != nil {
		logger.Log.Printf("[ERROR] apiHandler -> HandleOrderId: %v\n", err)
		http.Error(w, "Fehler beim Verarbeiten des Auftrags", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func StartServer() {
	http.HandleFunc("/api/neuerAuftrag", apiHandler)

	logger.Log.Printf("[INFO] Server running on: http://127.0.0.1:8080\n")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
