package main

import (
	"fmt"

	"./modules/commons"
	"./modules/database"
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

	headhunter2orders := make(chan commons.Order)
	watchdog2orders := make(chan bool)
	orders2watchdog := make(chan map[commons.Order]struct{})

	//vertices
	database.StartCSEDB(
		network2CSE,
		headhunter2CSE,
		CSE2headhunter)

	database.StartOrdersDB(
		ID,
		headhunter2orders,
		watchdog2orders,
		orders2watchdog)

	//wait forever
	select {}
}
