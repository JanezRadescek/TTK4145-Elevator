package database

import (
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
				if _, ok := orders[order.ID]; ok {
					//order is finished when Progress is Closing Door 2. if customer doesnt press button in ~10s driver should skip to Closing2
					if order.Progress == commons.ClosingDoor2 {
						delete(orders, order.ID)
					} else {
						//TODO check if recived order is newer version of order than what we allready have
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
					//To prevent sending multiple elevators. We are assuming infinitly sized elevator.
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
