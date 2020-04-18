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
	watchdogDriver2network <-chan commons.MessageStruct,
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
	reciver := make(chan commons.MessageStruct)

	go bcast.Transmitter(port, transmiter)
	go bcast.Receiver(port, reciver)

	for {
		var message commons.MessageStruct
		select {
		case message = <-watchdogDriver2network:
			{
				//fmt.Println("Network recived message from driver or watchdog ", message)
			}
		case message = <-reciver:
			{
				// if we get our msg back from UDP broadcast just discard it, because we litteraly had the same message nanosecond ago.
				if message.SenderID == ID {
					//fmt.Println("Network discarding message")
					continue
				} else {
					fmt.Println("Network recived message from internet ", message)
				}
			}
		}
		switch message.What {
		case commons.CSE, commons.Malfunction:
			{
				sendCSE <- message
			}
		case commons.Order:
			{
				//TODO check if updated orders id is same as contractor. it would be weird/buggy if someone else was doing work. OR would it?
				sendOrder <- message.Order
			}
		default:
			{
				//something wrong with code if anything else
				fmt.Println("Semantic Bug")
			}
		}
		if !message.Local {
			transmiter <- message
		}
	}

}
