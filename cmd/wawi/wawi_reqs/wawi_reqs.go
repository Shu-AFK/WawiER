package wawi_reqs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Shu-AFK/WawiER/cmd/defines"
)

func QuerySalesOrders(id string) (*Order, error) {
	var salesOrder Order
	url := fmt.Sprintf("%ssalesOrders?salesOrderNumber=%s", defines.APIBaseURL, id)
	resp, err := wawiCreateRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query sales order request failed (status %d), body: [%s]", resp.StatusCode, body)
	}

	err = json.Unmarshal(body, &salesOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sales order response: %v", err)
	}

	return &salesOrder, nil
}

func QuerySalesOrderItems(id int) ([]OrderItem, error) {
	url := fmt.Sprintf("%ssalesOrders/%d/lineitems", defines.APIBaseURL, id)
	resp, err := wawiCreateRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query sales order items request failed (status %d), body: [%s]", resp.StatusCode, body)
	}

	var items []OrderItem
	err = json.Unmarshal(body, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sales order items response: %v", err)
	}

	return items, nil
}

func GetStockData(itemId int) ([]StockItem, error) {
	url := fmt.Sprintf("%sstocks?itemId=%d", defines.APIBaseURL, itemId)
	resp, err := wawiCreateRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get stock data request failed (status %d), body: [%s]", resp.StatusCode, body)
	}

	var stockItems StockResponse
	err = json.Unmarshal(body, &stockItems)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stock data response: %v", err)
	}

	return stockItems.Items, nil
}

func wawiCreateRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	apiKey := os.Getenv(defines.APIKeyVarName)
	if apiKey == "" {
		return nil, fmt.Errorf("API key environment variable not set")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Wawi %v", apiKey))
	req.Header.Set("x-appid", defines.AppID)
	req.Header.Set("x-runas", defines.AppID)
	req.Header.Set("api-version", defines.APIVersion)
	req.Header.Set("x-appversion", defines.Version)

	if method == "POST" || method == "PATCH" {
		req.Header.Set("Content-Type", "application/json")
	}

	log.Printf("[INFO] Made request to: %s\n", req.URL)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
