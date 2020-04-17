package headhunter

import (
	"fmt"

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
		fmt.Println("headhunter recived order", order)
		if order.Progress <= commons.OpeningDoor1 {
			requestCopy <- true
			elevators := <-reciveCopy
			contractor := ID
			vector := elevators[ID].CurentFloor - order.DestinationFloor

			for _, tempE := range elevators {
				tempV := tempE.CurentFloor - order.DestinationFloor
				if (tempV*tempV < vector*vector) && tempE.Idle && order.StartingTime.After(tempE.LastTimeChecked) {
					vector = tempV
					contractor = tempE.ID
				}
			}
			order.Contractor = contractor
		} else {
			//There is nothing we can do past PRogress commons.OpeningDoor1
		}
		sendOrder <- order
	}
}
