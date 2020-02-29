package main

import (
	"fmt"
	"time"

	"./modules/commons"
)

func main() {
	fmt.Println("Starting elevator.")
	time.Sleep(time.Second)
	//init
	//newword-database
	netBase := make(chan commons.Message)

	//wait forever
	select {}
}
