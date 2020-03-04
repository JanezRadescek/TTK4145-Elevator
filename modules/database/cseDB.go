package database

import (
	"time"

	"../commons"
)

const elevatorTimeOut int = 10 //10s

//StartCSEDB starts thread save data base for CSE
func StartCSEDB(
	reciveCSE <-chan commons.MessageStruct,
	requestCopy <-chan bool,
	sendCopy chan<- map[string]commons.ElevatorStruct,
) {

	//key elevators ID, value elvator
	elevators := make(map[string]commons.ElevatorStruct)

	deleteOfflineElevators := make(chan bool)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			deleteOfflineElevators <- true
		}
	}()

	for {
		//to prevent race conditions we allways finish case before going into new loop. No go function here.
		select {
		case message := <-reciveCSE:
			{
				elevators[message.SenderID] = message.Elevator
			}
		case <-deleteOfflineElevators:
			{
				tempT := time.Now()
				tempT.Add(time.Duration(-elevatorTimeOut))
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
