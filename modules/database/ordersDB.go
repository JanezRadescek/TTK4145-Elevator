package database

import (
	"fmt"

	"../commons"
)

//StartOrdersDB starts thread save data base for orders
func StartOrdersDB(
	order <-chan commons.MessageStruct,
	requestedCopy <-chan bool,
	sendCopy chan<- map[int]commons.OrderStruct,
) {

	//key order , value time of last update
	orders := make(map[int]commons.OrderStruct)
	IDcounter := 1

	for {
		//to prevent race conditions we allways finish case before going into new loop. No go function here.
		select {
		case tempM := <-order:
			{
				fmt.Println("Got order")

				tempO := tempM.Order

				if _, ok := orders[tempO.ID]; ok {
					if tempO.Progress == commons.ClosingDoor2 {
						delete(orders, tempO.ID)
					} else {
						orders[tempO.ID] = tempO
					}
				} else {
					tempO.ID = IDcounter
					orders[IDcounter] = tempO
					IDcounter++
				}
			}
		case <-requestedCopy:
			{
				sendCopy <- orders
			}
		}
	}
}
