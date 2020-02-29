package main

import (
	"fmt"
	"time"

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
	network2CSE := make(chan commons.Message)
	headhunter2CSE := make(chan bool)
	CSE2headhunter := make(chan map[string]commons.Elevator)

	headhunter2orders := make(chan commons.Message)
	driver2Orders := make(chan commons.Order)
	watchdog2orders := make(chan bool)
	orders2watchdog := make(chan map[commons.Order]time.Time)

	watchdog2network := make(chan commons.Message)
	watchdog2driver := make(chan commons.Order)

	//vertices
	database.StartCSEDB(
		network2CSE,
		headhunter2CSE,
		CSE2headhunter,
	)

	database.StartOrdersDB(
		headhunter2orders,
		driver2Orders,
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
