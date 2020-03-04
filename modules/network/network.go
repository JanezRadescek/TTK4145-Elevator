package network

import (
	"fmt"
	"os"

	"../commons"
	"./Network-go/network/bcast"
	"./Network-go/network/localip"
)

const port int = 16569

//StartNetwork is responsible to make sure "im alive" is being send all the time
func StartNetwork(
	id chan<- string,
	reciver <-chan commons.MessageStruct,
	sendCSE chan<- commons.MessageStruct,
	sendOrder chan<- commons.OrderStruct,
) {
	//TODO
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	ID := fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	id <- ID

	transmiter := make(chan commons.MessageStruct)

	go bcast.Transmitter(port, transmiter)
	go bcast.Receiver(port, reciver)

	for {
		select {
		case message := <-reciver:
			{
				if message.Local || message.SenderID != ID {
					switch message.What {
					case commons.CSE:
						{
							sendCSE <- message
						}
					case commons.Order:
						{
							sendOrder <- message.Order
						}
					default:
						{
							fmt.Println("Semantic Bug")
						}
					}
				} else {
					transmiter <- message
				}
			}
		}
	}

}
