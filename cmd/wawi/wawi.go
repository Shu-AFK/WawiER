package wawi

import (
	"fmt"
	"log"
	"strings"

	"github.com/Shu-AFK/WawiER/cmd/config"
	"github.com/Shu-AFK/WawiER/cmd/email"
	"github.com/Shu-AFK/WawiER/cmd/structs"
	"github.com/Shu-AFK/WawiER/cmd/wawi/wawi_reqs"
)

func HandleOrderId(orderInfo structs.OrderReq) error {
	if !CheckIfNotExcluded(orderInfo.OrderId) {
		log.Printf("[ERROR] Order number %v is in exlusion list\n", orderInfo.OrderId)
		return nil
	}

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

	var emailItems []string // slice statt string

	for _, item := range items {
		stockData, err := wawi_reqs.GetStockData(item.ItemId)
		if err != nil {
			return err
		}

		var totalQuantity, totalUnavailable float64
		for _, s := range stockData {
			totalQuantity += s.QuantityTotal
			totalUnavailable += s.QuantityLockedForShipment +
				s.QuantityInPickingLists +
				s.QuantityLockedForAvailability +
				item.Quantity
		}

		if totalUnavailable > totalQuantity {
			log.Printf("[INFO] Order %v: Item %v is oversold (missing: %v)\n",
				orderInfo.OrderId, item.ItemId, totalUnavailable-totalQuantity)

			emailItems = append(emailItems,
				fmt.Sprintf("Artikel %s (%s): Bestellt: %v, Vorhanden: %v",
					item.Name, item.SKU, item.Quantity, totalQuantity))
		} else {
			log.Printf("[INFO] Order %v: Item %v is not oversold\n",
				orderInfo.OrderId, item.ItemId)
		}
	}

	if len(emailItems) > 0 {
		log.Printf("[INFO] Sending email for order %v to %s\n",
			orderInfo.OrderId, order.Items[0].Shipmentaddress.EmailAddress)

		customerName := fmt.Sprintf("%s %s",
			order.Items[0].Shipmentaddress.FirstName,
			order.Items[0].Shipmentaddress.LastName)

		// send slice instead of single string
		email.SendEmail(
			order.Items[0].Shipmentaddress.EmailAddress,
			emailItems,
			customerName,
			orderInfo.OrderId,
		)
	}

	return nil
}

func CheckIfNotExcluded(id string) bool {
	for _, excl := range config.Conf.ExcludedOrderIdStart {
		if strings.HasPrefix(id, excl) {
			return false
		}
	}

	return true
}
