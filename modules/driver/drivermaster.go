package driver

import (
	"strconv"
	"time"

	"../commons"
)

const sendUpdateDelay = 500 * time.Millisecond

//StartDriver takes next order we are asigned and opens door, changes floor acordingly
func StartDriverMaster(
	ID string,
	reciveCopy <-chan map[string]commons.OrderStruct,
	sendMessege chan<- commons.MessageStruct,
) {

	myself := commons.ElevatorStruct{}
	curentOrder := commons.OrderStruct{}
	//key is?
	allOurOrders := make(map[string]commons.OrderStruct)
	activeOrders := make(map[string]commons.OrderStruct)
	oldestTime := time.Now()
	IDcounter := 1

	time4Update := make(chan bool)
	go func() {
		for {
			time.Sleep(sendUpdateDelay)
			time4Update <- true
		}

	}()

	pickButton := make(chan commons.PickButtonStruct)
	floorButton := make(chan int)
	floorSensor := make(chan int)
	//stopButton := make(chan bool) //solve this on IO level
	setMotorDirection := make(chan int)
	setLamp := make(chan commons.LampStruct)
	DoorChan := make(chan bool)

	go StartDriverSlave()

	for {
		select {
		case allOurOrders = <-reciveCopy:
			{
				//find the oldest
				for _, order := range allOurOrders {
					if ID == order.Contractor && order.StartingTime.Before(oldestTime) {
						oldestTime = order.StartingTime
						curentOrder = order
					}
				}
				//find other order we may do on the way to the oldest first.
				if myself.Operational {
					myself.CurentDestination = curentOrder.DestinationFloor
					vector := curentOrder.DestinationFloor - myself.CurentFloor
					for key, order := range allOurOrders {
						tempV1 := order.DestinationFloor - myself.CurentFloor
						tempV2 := order.Direction
						if (tempV1*vector > 0) && (tempV2*vector > 0) {
							activeOrders[key] = order
							if tempV1*tempV1 < vector*vector {
								curentOrder = order
							}

						}
					}
				}

			}

		case button := <-pickButton:
			{
				order := commons.OrderStruct{
					ID:               ID + ":" + strconv.Itoa(IDcounter),
					Progress:         commons.ButtonPressed,
					Direction:        button.Direction,
					DestinationFloor: button.Floor,
					StartingTime:     time.Now(),
					//UpdateTime:       time.Now(),
					Contractor: "",
				}
				IDcounter++
				message := commons.MessageStruct{
					SenderID: ID,
					What:     commons.Order,
					Local:    false,
					Order:    order,
				}
				sendMessege <- message
				setLamp <- commons.LampStruct{button.Floor, true}

			}
		case floor := <-floorButton:
			{
				newOrder := true //two customers might wanna go to 2 diferent floor
				for _, order := range activeOrders {
					if order.DestinationFloor == floor {
						newOrder = false
					}
					//update orders
					if order.Progress == 5 {
						order.Progress = 6
						order.DestinationFloor = floor

						tempM := commons.MessageStruct{
							SenderID: ID,
							What:     commons.Order,
							Local:    false,
							Order:    order,
						}
						sendMessege <- tempM
					}
				}
				if newOrder {
					tempD := -1
					if myself.CurentFloor < floor {
						tempD = 1
					}
					order := commons.OrderStruct{
						ID:               ID + ":" + strconv.Itoa(IDcounter),
						Progress:         commons.Moving2destination,
						Direction:        tempD,
						DestinationFloor: floor,
						StartingTime:     time.Now(),
						//UpdateTime:       time.Now(),
						Contractor: "",
					}
					IDcounter++
					message := commons.MessageStruct{
						SenderID: ID,
						What:     commons.Order,
						Local:    false,
						Order:    order,
					}
					sendMessege <- message
				}
				setMotorDirection <- curentOrder.Direction
			}
		case floor := <-floorSensor:
			{
				myself.CurentFloor = floor
				openDoor := false
				for _, order := range activeOrders {

					switch order.Progress {
					case commons.ButtonPressed:
						{
							order.Progress = commons.Moving2customer
							tempM := commons.MessageStruct{
								SenderID: ID,
								What:     commons.Order,
								Local:    false,
								Order:    order,
							}
							sendMessege <- tempM
							if order.DestinationFloor == floor {
								openDoor = true
							}
						}
					case commons.Moving2customer, commons.Moving2destination:
						{
							if order.DestinationFloor == floor {
								openDoor = true
							}
						}
					}

				}
				DoorChan <- openDoor
			}

		case door := <-DoorChan:
			{
				for _, order := range activeOrders {
					if door {
						if order.Progress < commons.OpeningDoor1 {
							order.Progress = commons.OpeningDoor1
						} else if order.Progress == commons.Moving2destination {
							order.Progress = commons.OpeningDoor2
						}
						message := commons.MessageStruct{
							SenderID: ID,
							What:     commons.Order,
							Local:    false,
							Order:    order,
						}
						sendMessege <- message
					} else {
						if order.Progress == commons.OpeningDoor1 {
							order.Progress = commons.ClosingDoor1
						} else if order.Progress == commons.OpeningDoor2 {
							order.Progress = commons.ClosingDoor2
						}
						message := commons.MessageStruct{
							SenderID: ID,
							What:     commons.Order,
							Local:    false,
							Order:    order,
						}
						sendMessege <- message
					}
				}
			}

		case <-time4Update:
			{
				tempM := commons.MessageStruct{
					SenderID: ID,
					What:     commons.CSE,
					Local:    false,
					Elevator: myself,
				}
				sendMessege <- tempM
			}

		}

	}
}
