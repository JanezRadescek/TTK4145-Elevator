package driver

import (
	"strconv"
	"time"

	"../commons"
)

const sendUpdateDelay = 500 * time.Millisecond

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
	time4Update := make(chan bool)
	go func() {
		for {
			time.Sleep(sendUpdateDelay)
			time4Update <- true
		}

	}()

	for {
		select {
		case orders = <-reciveCopy:
			{
				//find the oldest
				for _, order := range orders {
					if ID == order.Contractor && order.StartingTime.Before(oldestTime) {
						oldestTime = order.StartingTime
						curentOrder = order
					}
				}
				//find other order we may do on the way to the oldest first.
				if myself.Operational {
					myself.CurentDestination = curentOrder.DestinationFloor
					vector := curentOrder.DestinationFloor - myself.CurentFloor
					for _, order := range orders {
						tempV1 := order.DestinationFloor - myself.CurentFloor
						tempV2 := order.Direction
						if (tempV1*vector > 0) && (tempV2*vector > 0) && (tempV1*tempV1 < vector*vector) {
							curentOrder = order
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
					//UpdateTime:       time.Now(),
					Contractor: "",
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

		case <-time4Update:
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
