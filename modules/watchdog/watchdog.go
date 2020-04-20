package watchdog

import (
	"time"

	"../commons"
)

//StartWatchDog will periodicly check if orders are being executed as expected
func StartWatchDog(
	ID string,
	requestCopy chan<- bool,
	reciveCopy <-chan map[string]commons.OrderStruct,
	sendMessege chan<- commons.MessageStruct,
	sendOurOrders chan<- map[string]commons.OrderStruct,
) {
	go func() {
		for {
			time.Sleep(commons.WatchDogFrequency)
			requestCopy <- true
		}
	}()

	for {
		tempOrders := <-reciveCopy
		//key is ID
		ourOrders := make(map[string]commons.OrderStruct)

		curentTime := time.Now()

		for _, order := range tempOrders {

			//someone can call cab but not enter. Elevator is waiting for destination button which will not be pressed. so delete order.
			tempT := order.LastUpdate.Add(commons.MaxUserTime)
			if tempT.Before(curentTime) && order.Progress == commons.WaitingForDestination {
				order.Progress = commons.ClosingDoor2
				sendMessege <- commons.MessageStruct{
					SenderID: order.Contractor,
					What:     commons.Order,
					Local:    true,
					Order:    order,
				}
			} else {
				tempT = order.StartingTime.Add(commons.MaxOrderTime)

				if order.Contractor == ID {
					//its pointless to check on ourself if we are performing to spec.
					ourOrders[order.ID] = order
				} else {
					if tempT.Before(curentTime) && order.Progress <= commons.OpeningDoor1 {
						//Once the progress is 4 or more we cant switch elevator. In case of Failure at this stage unfortunatly  customer must die:)
						sendMessege <- commons.MessageStruct{
							SenderID: order.Contractor, //we seend messege in the name of contractor.
							What:     commons.Malfunction,
							Local:    true,
							Elevator: commons.ElevatorStruct{LastTimeChecked: time.Now()},
						}
						order.Contractor = ""
						sendMessege <- commons.MessageStruct{
							SenderID: order.Contractor,
							What:     commons.Order,
							Local:    true,
							Order:    order,
						}
					}
				}
			}
		}

		sendOurOrders <- ourOrders
	}
}
