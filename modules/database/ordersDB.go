package database

import (
	"fmt"

	"../commons"
)

//StartOrdersDB starts thread save data base for orders
func StartOrdersDB(
	reciveOrder <-chan commons.OrderStruct,
	requestedCopy <-chan bool,
	sendCopy chan<- map[string]commons.OrderStruct,
) {
	orders := make(map[string]commons.OrderStruct)

	for {
		select {
		case order := <-reciveOrder:
			{
				fmt.Println("ordersDB recived order", order)
				if _, ok := orders[order.ID]; ok {
					//order is finished when Progress is Closing Door 2. if customer doesnt press button in ~10s driver should skip to Closing2
					if order.Progress == commons.ClosingDoor2 {
						delete(orders, order.ID)
					} else {
						//TODO check if recived order is newer version of order than what we allready have
						orders[order.ID] = order
					}
				} else {
					//to prevent a mess if multiple users want to go with same elevator to the same destination
					unique := true
					for _, tempO := range orders {
						if tempO.DestinationFloor == order.DestinationFloor &&
							tempO.Progress <= commons.OpeningDoor1 &&
							order.Progress <= commons.OpeningDoor1 {
							unique = false
							break
						}
					}
					// We are assuming infinitly sized elevator.
					if unique {
						orders[order.ID] = order
					}

				}
			}
		case <-requestedCopy:
			{
				sendCopy <- orders
			}
		}
	}
}
