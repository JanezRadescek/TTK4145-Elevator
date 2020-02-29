package watchdog

import (
	"strings"
	"time"

	"../commons"
)

const frequency int = 500 //2Hz
const maxTime int = 30    //30s

//StartWatchDog will periodicly check if orders are being executed as expected
func StartWatchDog(
	ID string,
	requestCopy chan<- bool,
	reciveCopy <-chan map[commons.Order]time.Time,
	sendMessege chan<- commons.Message,
	sendOurOrder chan<- commons.Order,
) {

	for {
		time.Sleep(time.Duration(frequency) * time.Millisecond)

		requestCopy <- true
		tempOrders := <-reciveCopy

		tempT := time.Now()

		var oldestOurOrder commons.Order
		var oldestTime time.Time = tempT

		for order, ordersTime := range tempOrders {
			if ID == order.Contractor && ordersTime.Before(oldestTime) {
				oldestTime = ordersTime
				oldestOurOrder = order
			}

			ordersTime = ordersTime.Add(time.Duration(maxTime) * time.Second)
			if ordersTime.Before(tempT) {

				//TODO kick vote

				tempID := strings.Split(order.Contractor, ":")
				tempIP := tempID[0]
				tempProcessID := tempID[1]

				tempM1 := commons.Message{
					SenderIP:        tempIP,
					SenderProcessID: tempProcessID,
					What:            commons.LocalCSE,
					Elevator:        commons.Elevator{Operational: false},
				}
				sendMessege <- tempM1

				tempM2 := commons.Message{
					SenderIP:        tempIP,
					SenderProcessID: tempProcessID,
					What:            commons.LocalOrder,
					Order:           order,
				}
				sendMessege <- tempM2

			} else {
				//everything OK.
			}
		}
		if oldestOurOrder.Contractor != "" {
			sendOurOrder <- oldestOurOrder
		}
	}
}
