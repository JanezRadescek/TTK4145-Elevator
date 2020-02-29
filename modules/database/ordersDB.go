package database

import (
	"fmt"
	"time"

	"../commons"
)

//StartOrdersDB starts thread save data base for orders
func StartOrdersDB(
	newOrder <-chan commons.Message,
	deleteOrder <-chan commons.Order,
	requestedCopy <-chan bool,
	sendCopy chan<- map[commons.Order]time.Time,
) {

	//key order , value time of last update
	orders := make(map[commons.Order]time.Time)

	for {
		//to prevent race conditions we allways finish case before going into new loop. No go function here.
		select {
		case tempM := <-newOrder:
			{
				fmt.Println("Got order")
				//TODO check if it is order in progress
				orders[tempM.Order] = time.Now()
			}
		case tempO := <-deleteOrder:
			{
				delete(orders, tempO)
			}
		case <-requestedCopy:
			{
				sendCopy <- orders
			}
		}
	}
}
