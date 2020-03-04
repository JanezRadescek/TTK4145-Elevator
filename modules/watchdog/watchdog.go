package watchdog

import (
	"time"

	"../commons"
)

const delay int = 500  //2Hz //
const maxTime int = 30 //30s

//StartWatchDog will periodicly check if orders are being executed as expected
func StartWatchDog(
	ID string,
	requestCopy chan<- bool,
	reciveCopy <-chan map[string]commons.OrderStruct,
	sendMessege chan<- commons.MessageStruct,
	sendOurOrders chan<- map[string]commons.OrderStruct,
) {

	for {
		//TODO instead check when it makes sanse to check
		time.Sleep(time.Duration(delay) * time.Millisecond)

		requestCopy <- true
		tempOrders := <-reciveCopy
		ourOrders := make(map[string]commons.OrderStruct)

		curentTime := time.Now()

		for _, order := range tempOrders {
			tempT := order.StartingTime.Add(time.Duration(maxTime) * time.Second)

			if order.Contractor == ID {
				ourOrders[order.ID] = order
			} else {
				//its pointless to check on ourself if we are performing to spec.
				if tempT.Before(curentTime) && order.Progress <= 3 {

					tempM1 := commons.MessageStruct{
						SenderID: order.Contractor,
						What:     commons.CSE,
						Local:    true,
						Elevator: commons.ElevatorStruct{Operational: false},
					}
					sendMessege <- tempM1

					tempM2 := commons.MessageStruct{
						SenderID: order.Contractor,
						What:     commons.Order,
						Local:    true,
						Order:    order,
					}
					sendMessege <- tempM2

				}
			}

		}
		if len(ourOrders) != 0 {
			sendOurOrders <- ourOrders
		}
	}
}

// func updateElevatorDB() {

// }

// func updateOrderDB() {

// }
