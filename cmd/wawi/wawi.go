package wawi

import (
	"fmt"
	"log"

	"github.com/Shu-AFK/WawiER/cmd/email"
	"github.com/Shu-AFK/WawiER/cmd/structs"
	"github.com/Shu-AFK/WawiER/cmd/wawi/wawi_reqs"
)

func HandleOrderId(orderInfo structs.OrderReq) error {
	order, err := wawi_reqs.QuerySalesOrders(orderInfo.OrderId)
	if err != nil {
		return err
	}
	if order.TotalItems > 1 {
		return fmt.Errorf("order has more than one item")
	}

	log.Printf("[INFO] Working on order: %v/%v\n", order.Items[0].Id, orderInfo.OrderId)

	items, err := wawi_reqs.QuerySalesOrderItems(order.Items[0].Id)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Got %v items in sales order %v, checking for oversell\n", len(items), orderInfo.OrderId)

	emailItemString := ""
	for _, item := range items {
		stockData, err := wawi_reqs.GetStockData(item.ItemId)
		if err != nil {
			return err
		}

		var totalQuantity, totalUnavailable float64
		for _, s := range stockData {
			totalQuantity += s.QuantityTotal
			totalUnavailable += s.QuantityLockedForShipment + s.QuantityInPickingLists + s.QuantityLockedForAvailability
		}

		if totalUnavailable > totalQuantity {
			missing := totalUnavailable - totalQuantity
			log.Printf("[INFO] Order %v: Item %v is oversold (missing: %v)\n", orderInfo.OrderId, item.ItemId, missing)

			emailItemString += fmt.Sprintf("Artikel %s (%s): Bestellt: %v, Vorhanden: %v\n", item.Name, item.SKU, item.Quantity, missing)
		} else {
			log.Printf("[INFO] Order %v: Item %v is not oversold\n", orderInfo.OrderId, item.ItemId)
		}
	}

	if emailItemString != "" {
		log.Printf("[INFO] Sending email for order %v to %s\n", orderInfo.OrderId, order.Items[0].Shipmentaddress.EmailAddress)
		customerName := fmt.Sprintf("%s %s", order.Items[0].Shipmentaddress.FirstName, order.Items[0].Shipmentaddress.LastName)
		email.SendEmail(order.Items[0].Shipmentaddress.EmailAddress, emailItemString, customerName, orderInfo.OrderId)
	}

	return nil
}
