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
		//to prevent race conditions we allways finish case before going into new loop. No go function here.
		select {
		case order := <-reciveOrder:
			{
				fmt.Println("Got order")

				if _, ok := orders[order.ID]; ok {
					if order.Progress == commons.ClosingDoor2 {
						delete(orders, order.ID)
					} else {
						orders[order.ID] = order
					}
				} else {
					unique := true
					for _, tempO := range orders {
						if tempO.DestinationFloor == order.DestinationFloor && order.Progress <= 3 {
							unique = false
							break
						}
					}
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
