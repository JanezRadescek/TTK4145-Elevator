package main

import (
	"fmt"

	"./modules/commons"
	"./modules/database"
	"./modules/watchdog"
)

func main() {
	fmt.Println("Starting elevator.")

	//init
	ID := "TODO"

	//create graph

	//edges
	network2CSE := make(chan commons.MessageStruct)
	headhunter2CSE := make(chan bool)
	CSE2headhunter := make(chan map[string]commons.ElevatorStruct)

	headhunter2orders := make(chan commons.MessageStruct)
	driver2Orders := make(chan commons.OrderStruct)
	watchdog2orders := make(chan bool)
	orders2watchdog := make(chan map[int]commons.OrderStruct)

	watchdog2network := make(chan commons.MessageStruct)
	watchdog2driver := make(chan commons.OrderStruct)

	//vertices
	database.StartCSEDB(
		network2CSE,
		headhunter2CSE,
		CSE2headhunter,
	)

	database.StartOrdersDB(
		headhunter2orders,
		watchdog2orders,
		orders2watchdog,
	)

	watchdog.StartWatchDog(
		ID,
		watchdog2orders,
		orders2watchdog,
		watchdog2network,
		watchdog2driver,
	)

	//wait forever
	select {}
}
