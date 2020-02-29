package database

import (
	"fmt"

	"../commons"
)

func start(
	cse <-chan commons.Message,
	requestCopy <-chan bool,
	sendCopy chan<- []commons.Elevator) {

	elevators := make(map[string]commons.Elevator)

	for {
		//to prevent race conditions we allways finish case before going into new loop. No go function here.
		select {
		case tempM := <-cse:
			{
				fmt.Println("TODO ", tempM)
				tempID := tempM.SenderIP + tempM.SenderPROCES
				elevators[tempID] = commons.StringToElevator(tempM.Text)
			}
		case <-requestCopy:
			{
				//sendCopy
			}
		}
	}
}
