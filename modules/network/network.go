package network

import (
	"../commons"
)

//StartNetwork is responsible to make sure "im alive" is being send all the time
func StartNetwork(
	id chan<- string,
	reciveMessageWatchDog <-chan commons.MessageStruct,
	reciveMessageDriver <-chan commons.MessageStruct,
	sendCSE chan<- commons.ElevatorStruct,
	sendOrder chan<- commons.OrderStruct,
) {
	//TODO
}
