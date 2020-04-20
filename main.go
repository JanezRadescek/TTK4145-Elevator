package main

import (
	"fmt"
	"os"

	"./modules/commons"
	"./modules/database"
	"./modules/driver"
	"./modules/headhunter"
	"./modules/network"
	"./modules/watchdog"
)

func main() {
	fmt.Println("Starting elevator.")

	arguments := os.Args
	//
	if len(arguments) > 0 {
		commons.ElevatorPort = arguments[1]
	}

	//create graph

	//edges
	id := make(chan string)
	network2CSE := make(chan commons.MessageStruct)
	network2headhunter := make(chan commons.OrderStruct)

	headhunter2CSE := make(chan bool)
	CSE2headhunter := make(chan map[string]commons.ElevatorStruct)

	headhunter2orders := make(chan commons.OrderStruct)

	watchdog2orders := make(chan bool)
	orders2watchdog := make(chan map[string]commons.OrderStruct)

	watchdogDriver2network := make(chan commons.MessageStruct)
	watchdog2driver := make(chan map[string]commons.OrderStruct)

	//vertices
	go network.StartNetwork(
		id,
		watchdogDriver2network,
		network2CSE,
		network2headhunter,
	)

	ID := <-id

	go database.StartCSEDB(
		network2CSE,
		headhunter2CSE,
		CSE2headhunter,
	)

	go database.StartOrdersDB(
		headhunter2orders,
		watchdog2orders,
		orders2watchdog,
	)

	go headhunter.StartHeadHunter(
		ID,
		network2headhunter,
		headhunter2orders,
		headhunter2CSE,
		CSE2headhunter,
	)

	go watchdog.StartWatchDog(
		ID,
		watchdog2orders,
		orders2watchdog,
		watchdogDriver2network,
		watchdog2driver,
	)

	go driver.StartDriverMaster(
		ID,
		watchdog2driver,
		watchdogDriver2network,
	)

	//wait forever
	fmt.Println("waiting for CTRL + C")
	select {}
}
