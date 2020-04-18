package database

import (
	"fmt"
	"time"

	"../commons"
)

//StartCSEDB starts thread save data base for CSE
func StartCSEDB(
	reciveCSE <-chan commons.MessageStruct,
	requestCopy <-chan bool,
	sendCopy chan<- map[string]commons.ElevatorStruct,
) {

	//key elevators ID
	elevators := make(map[string]commons.ElevatorStruct)

	deleteOfflineElevators := make(chan bool)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			deleteOfflineElevators <- true
		}
	}()

	for {
		select {
		case message := <-reciveCSE:
			{

				message.Elevator.LastTimeOnline = time.Now()
				switch message.What {
				case commons.CSE:
					{
						if elevators[message.SenderID].CurentFloor != message.Elevator.CurentFloor ||
							elevators[message.SenderID].ID != message.Elevator.ID ||
							elevators[message.SenderID].Idle != message.Elevator.Idle {
							//fmt.Println("cseDB recived NEW CSE", message)
						}
						elevators[message.SenderID] = message.Elevator
					}
				case commons.Malfunction:
					{
						fmt.Println("cseDB recived malfunction", message)
						if elevator, ok := elevators[message.SenderID]; ok {
							elevator.LastTimeChecked = message.Elevator.LastTimeChecked
							elevators[message.SenderID] = elevator
						} else {
							//something is wrong with code if we get here
						}
					}
				}

			}
		case <-deleteOfflineElevators:
			{
				tempT := time.Now().Add(commons.MaxElevatorTime)
				for tempID, tempE := range elevators {
					if tempE.LastTimeOnline.Before(tempT) {
						delete(elevators, tempID)
					}
				}
			}
		case <-requestCopy:
			{
				sendCopy <- elevators
			}
		}
	}
}
