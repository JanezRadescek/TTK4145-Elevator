package main

import (
	"fmt"
	"os"
	"strconv"

	"./modules/commons"
	"./modules/database"
	"./modules/driver"
	"./modules/headhunter"
	"./modules/watchdog"
)

func main() {
	fmt.Println("Starting elevator.")

	//create graph

	//edges
	ip := make(chan string)
	network2CSE := make(chan commons.MessageStruct)
	network2headhunter := make(chan commons.OrderStruct)

	headhunter2CSE := make(chan bool)
	CSE2headhunter := make(chan map[string]commons.ElevatorStruct)

	headhunter2orders := make(chan commons.OrderStruct)

	watchdog2orders := make(chan bool)
	orders2watchdog := make(chan map[string]commons.OrderStruct)

	watchdog2network := make(chan commons.MessageStruct)
	watchdog2driver := make(chan map[string]commons.OrderStruct)

	driver2network := make(chan commons.MessageStruct)

	//vertices
	network.StartNetwork(
		ip,
		watchdog2network,
		driver2network,
		network2CSE,
		network2headhunter,
	)

	ID := <-ip
	ID = ID + ":" + strconv.Itoa(os.Getpid)

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

	headhunter.StartHeadHunter(
		ID,
		network2headhunter,
		headhunter2orders,
		headhunter2CSE,
		CSE2headhunter,
	)

	watchdog.StartWatchDog(
		ID,
		watchdog2orders,
		orders2watchdog,
		watchdog2network,
		watchdog2driver,
	)

	driver.StartDriver(
		ID,
		watchdog2driver,
		driver2network,
	)

	//wait forever
	select {}
}
