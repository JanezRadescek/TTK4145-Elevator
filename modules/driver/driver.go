package driver

import (
	"strconv"
	"time"

	"../commons"
)

const sendUpdateDelay = 1

//StartDriver takes next order we are asigned and opens door, changes floor acordingly
func StartDriver(
	ID string,
	reciveCopy <-chan map[int]commons.OrderStruct,
	sendMessege chan<- commons.MessageStruct,
) {
	myself := commons.ElevatorStruct{}
	curentOrder := commons.OrderStruct{}
	orders := make(map[int]commons.OrderStruct)
	oldestTime := time.Now()
	IDcounter := 1
	timeUp := make(chan bool)

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
				//find other order we may do on the way
				if myself.Operational {
					myself.CurentDestination = curentOrder.DestinationFloor
					vector := curentOrder.DestinationFloor - myself.CurentFloor
					for _, tempO := range orders {
						tempV1 := tempO.DestinationFloor - myself.CurentFloor
						tempV2 := tempO.Direction
						if (tempV1*vector > 0) && (tempV2*vector > 0) && (tempV1*tempV1 < vector*vector) {
							curentOrder = tempO
						}
					}

					drive(myself, curentOrder)
				}

			}

		case tempB <- buttom:
			{
				tempD := 1
				tempF := 10

				order := commons.OrderStruct{
					ID:               ID + ":" + strconv.Itoa(IDcounter),
					Progress:         commons.ButtonPressed,
					Direction:        tempD,
					DestinationFloor: tempF,
					StartingTime:     time.Now(),
					UpdateTime:       time.Now(),
					Contractor:       "",
				}
				IDcounter++
				tempM := commons.MessageStruct{
					SenderID: ID,
					What:     commons.Order,
					Local:    false,
					Order:    order,
				}
				sendMessege <- tempM
			}
			//case tempD<-door
			//case tempF<-floor

		case <-timeUp:
			{
				tempM := commons.MessageStruct{
					SenderID: ID,
					What:     commons.CSE,
					Local:    false,
					Elevator: myself,
				}
				sendMessege <- tempM
			}

		}

	}
}

func drive(myself commons.ElevatorStruct, curentOrder commons.OrderStruct) {
	//"state machine"
	//if we are in state of waiting for the floor button and we dont get respons in
	//TODO
}
