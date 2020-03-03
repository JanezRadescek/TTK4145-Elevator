package headhunter

import (
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
		if order.Progress <= 3 {
			requestCopy <- true
			elevators := <-reciveCopy
			contractor := ID
			vector := elevators[ID].CurentFloor - order.DestinationFloor

			for _, tempE := range elevators {
				tempV := tempE.CurentFloor - order.DestinationFloor
				if (tempV*tempV < vector*vector) && tempE.Idle {
					vector = tempV
					contractor = tempE.ID
				}
			}
			order.Contractor = contractor
		}
		sendOrder <- order
	}
}
