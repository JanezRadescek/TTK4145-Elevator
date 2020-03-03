package driver

import (
	"time"

	"../commons"
)

//StartDriver will periodicly check if orders are being executed as expected
func StartDriver(
	ID string,
	reciveCopy <-chan map[int]commons.OrderStruct,
	sendMessege chan<- commons.MessageStruct,
) {

	myself := commons.ElevatorStruct{}
	curentOrder := commons.OrderStruct{}
	orders := make(map[int]commons.OrderStruct)
	oldestTime := time.Now()

	for {
		select {
		case orders = <-reciveCopy:
			{
				//find the oldest
				for _, tempO := range orders {
					if ID == tempO.Contractor && tempO.StartingTime.Before(oldestTime) {
						oldestTime = tempO.StartingTime
						curentOrder = tempO
					}
				}
				//find all other we may do at the same time
				//TODO
			}
		case tempB := <-buttom:
			{
				//update myself
			}
		case tempF := <-floor:
			{

			}
		case tempD := <-door:
			{

			}
		case tempO := <-operational:
			{

			}

		}

	}
}
