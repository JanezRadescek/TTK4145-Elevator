package headhunter

import (
	"strings"

	"../commons"
)

//StartHeadHunter finds contracotr for orders
func StartHeadHunter(
	ID string,
	reciveOrder <-chan commons.OrderStruct,
	sendOrder chan<- commons.OrderStruct,
	requestCopy chan<- bool,
	reciveCopy <-chan map[string]commons.ElevatorStruct,
) {
	for {
		order := <-reciveOrder
		//fmt.Println("headhunter recived order", order)
		if order.Progress <= commons.OpeningDoor1 {
			requestCopy <- true
			elevators := <-reciveCopy
			contractor := ID
			vector := elevators[ID].CurentFloor - order.DestinationFloor

			for _, tempE := range elevators {
				tempV := tempE.CurentFloor - order.DestinationFloor
				if (tempV*tempV <= vector*vector) && tempE.Idle && order.StartingTime.After(tempE.LastTimeChecked) {
					if tempV*tempV == vector*vector && contractor < tempE.ID {
						vector = tempV
						contractor = tempE.ID
					}

				}
			}
			order.Contractor = contractor
		} else {
			//we cant distribute orders past Progress commons.OpeningDoor1 so the elevator that has people in it has to do it
			order.Contractor = strings.Split(order.ID, ":")[0]
		}
		sendOrder <- order
	}
}
