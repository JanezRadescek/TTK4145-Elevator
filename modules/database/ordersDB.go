package database

import (
	"fmt"

	"../commons"
)

//StartOrdersDB starts thread save data base for orders
func StartOrdersDB(
	ID string,
	order <-chan commons.Message,
	requestCopy <-chan bool,
	sendCopy chan<- map[string]commons.Order,
) {

	//key ID of contractor , value order
	orders := make(map[string]commons.Order)

	for {
		//to prevent race conditions we allways finish case before going into new loop. No go function here.
		select {
		case tempM := <-order:
			{
				tempID := tempM.SenderIP + ":" + tempM.SenderProcessID
				fmt.Println("Got order")
				orders[tempID] = tempM.Order
			}
		case <-requestCopy:
			{
				sendCopy <- orders
			}
		}
	}
}
