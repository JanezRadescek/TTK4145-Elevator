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
	reciveCopy <-chan map[int]commons.OrderStruct,
	sendMessege chan<- commons.MessageStruct,
	sendCopy chan<- map[int]commons.OrderStruct,
) {

	for {
		time.Sleep(time.Duration(frequency) * time.Millisecond)

		requestCopy <- true
		tempOrders := <-reciveCopy

		curentTime := time.Now()

		for _, order := range tempOrders {
			tempT := order.StartingTime.Add(time.Duration(maxTime) * time.Second)
			if tempT.Before(curentTime) {

				//TODO kick vote

				tempID := strings.Split(order.Contractor, ":")
				tempIP := tempID[0]
				tempProcessID := tempID[1]

				tempM1 := commons.MessageStruct{
					SenderIP:        tempIP,
					SenderProcessID: tempProcessID,
					What:            commons.CSE,
					Local:           true,
					Elevator:        commons.ElevatorStruct{Operational: false},
				}
				sendMessege <- tempM1

				tempM2 := commons.MessageStruct{
					SenderIP:        tempIP,
					SenderProcessID: tempProcessID,
					What:            commons.Order,
					Local:           true,
					Order:           order,
				}
				sendMessege <- tempM2

			} else {
				//everything OK.
			}
		}
		if oldestOurOrder.Contractor != "" {
			sendOurOrder <- tempOrders
		}
	}
}

// func updateElevatorDB() {

// }

// func updateOrderDB() {

// }
