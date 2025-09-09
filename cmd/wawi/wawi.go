package wawi

import (
	"log"

	"github.com/Shu-AFK/WawiER/cmd/structs"
	"github.com/Shu-AFK/WawiER/cmd/wawi/wawi_reqs"
)

func HandleOrderId(orderInfo structs.OrderReq) error {
	order, err := wawi_reqs.QuerySalesOrders(orderInfo.OrderId)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Working on order: %v/%v\n", order.Items[0].Id, orderInfo.OrderId)

	items, err := wawi_reqs.QuerySalesOrderItems(order.Items[0].Id)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Got %v items in sales order, checking for oversell\n", len(items))

	for _, item := range items {
		stockData, err := wawi_reqs.GetStockData(item.ItemId)
		if err != nil {
			return err
		}
	}

	return nil
}
