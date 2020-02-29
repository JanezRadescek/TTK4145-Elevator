package database

import (
	"fmt"
	"time"

	"../commons"
)

//StartCSEDB starts thread save data base for CSE
func StartCSEDB(
	cse <-chan commons.Message,
	requestCopy <-chan bool,
	sendCopy chan<- map[string]commons.Elevator,
) {

	//key elevators ID, value elvator
	elevators := make(map[string]commons.Elevator)

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
		case tempM := <-cse:
			{
				tempID := tempM.SenderIP + ":" + tempM.SenderProcessID
				fmt.Println("Got update from ", tempID)
				elevators[tempID] = tempM.Elevator
			}
		case <-deleteOfflineElevators:
			{
				//TODO delete elevators that are ofline for longer than 30s
			}
		case <-requestCopy:
			{
				sendCopy <- elevators
			}
		}
	}
}
