package driver

import (
	"fmt"
	"strconv"
	"time"

	"../commons"
	"./driver-go/elevio"
)

var myself commons.ElevatorStruct
var curentOrder commons.OrderStruct
var allOurOrders map[string]commons.OrderStruct
var activeOrders map[string]commons.OrderStruct

var privateSendMessege chan<- commons.MessageStruct
var privateID string
var setMotorDirection chan int

//TODO properly react to "disruptor" presing buttons at rendom times. (like what happens if somebody somehow pushes button for floor 5 before we even let him in?)

//StartDriverMaster takes next order we are asigned and and give high level instruction on what to do with it.
func StartDriverMaster(
	ID string,
	reciveCopy <-chan map[string]commons.OrderStruct,
	sendMessege chan<- commons.MessageStruct,
) {
	privateID = ID
	privateSendMessege = sendMessege

	myself = commons.ElevatorStruct{}
	myself.LastTimeChecked = time.Now()
	myself.Idle = true

	//key is ID
	allOurOrders = make(map[string]commons.OrderStruct)

	IDcounter := 1

	time4Update := make(chan bool)
	go func() {
		for {
			time.Sleep(commons.SendUpdateDelay)
			time4Update <- true
		}

	}()

	newButton := make(chan elevio.ButtonEvent)
	floorSensor := make(chan int)
	doorSensor := make(chan bool)
	//stopButton := make(chan bool) //solve this on IO level
	setMotorDirection = make(chan int)
	setOpenDoor := make(chan bool)

	go StartDriverSlave(newButton, floorSensor, doorSensor, setMotorDirection, setOpenDoor)

	for {
		select {
		case allOurOrders = <-reciveCopy:
			{
				//fmt.Println("drivermaster recived copy", allOurOrders)

				//find the oldest. or if it doesnt exist enjoy
				if len(allOurOrders) != 0 {
					oldestTime := time.Now()
					for _, order := range allOurOrders {
						oldestTime = order.StartingTime
						curentOrder = order
						myself.Idle = false
						break
					}
					for _, order := range allOurOrders {
						if order.StartingTime.Before(oldestTime) {
							oldestTime = order.StartingTime
							curentOrder = order
						}
					}
					//find other order we may do on the way to the oldest first and start doing it.
					findCurentOrder()
				} else {
					myself.Idle = true
				}

			}

		case button := <-newButton:
			{
				switch button.Button {
				case elevio.BT_Cab:
					{

						floor := button.Floor
						fmt.Println("drivermaster recived cab to floor ", floor)
						newOrder := true //two customers might wanna go to 2 diferent floor
						for _, order := range activeOrders {
							if order.DestinationFloor == floor {
								newOrder = false
							}
							//update orders
							if order.Progress == commons.WaitingForDestination {
								order.Progress = commons.Moving2destination
								order.DestinationFloor = floor
								//activeOrders[key] = order //probably donnt need it.
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

							tempT := time.Now()

							order := commons.OrderStruct{
								ID:               ID + ":" + strconv.Itoa(IDcounter),
								Progress:         commons.Moving2destination,
								Direction:        0, //should only be used for progress == commons.buttonpressed.
								DestinationFloor: floor,
								StartingTime:     tempT,
								LastUpdate:       time.Now(),
								Contractor:       "",
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
						findCurentOrder()
					}
				case elevio.BT_HallUp, elevio.BT_HallDown:
					{
						floor := button.Floor
						fmt.Println("drivermaster recived order from floor ", floor)
						direction := 1 //up
						if button.Button == elevio.BT_HallDown {
							direction = -1
						}
						order := commons.OrderStruct{
							ID:               ID + ":" + strconv.Itoa(IDcounter),
							Progress:         commons.ButtonPressed,
							Direction:        direction, //should only be used in progress button pressed. after this it should be calculated as it is relative to elevator position.
							DestinationFloor: floor,
							StartingTime:     time.Now(),
							LastUpdate:       time.Now(),
							Contractor:       "",
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
				}

			}

		case floor := <-floorSensor:
			{
				fmt.Println("drivermaster recived floor ", floor)
				myself.CurentFloor = floor
				openDoor := false
				//fmt.Println()
				//fmt.Println("drivermaster printing active orders ")
				for _, order := range activeOrders {
					//fmt.Println("	drivermaster printing order ", order)
					switch order.Progress {
					case commons.ButtonPressed:
						{
							order.Progress = commons.Moving2customer
							order.LastUpdate = time.Now()
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
				if openDoor {

					setOpenDoor <- true
				}
				if myself.Idle {
					setMotorDirection <- 0
				}

			}

		case door := <-doorSensor:
			{
				fmt.Println("drivermaster doorsensor ", door)
				for _, order := range activeOrders {
					if door {
						if order.Progress < commons.OpeningDoor1 {
							order.Progress = commons.OpeningDoor1
						} else if order.Progress == commons.Moving2destination {
							order.Progress = commons.OpeningDoor2
						}
					} else {
						if order.Progress == commons.OpeningDoor1 {
							order.Progress = commons.ClosingDoor1
							order.Progress = commons.WaitingForDestination //door closes in a instant
						} else if order.Progress == commons.OpeningDoor2 {
							order.Progress = commons.ClosingDoor2
						}
					}
					order.LastUpdate = time.Now()
					//activeOrders[key] = order //dont needed if watchdog freaqueny high enough
					message := commons.MessageStruct{
						SenderID: ID,
						What:     commons.Order,
						Local:    false,
						Order:    order,
					}
					sendMessege <- message
				}
				if !door {
					//before door closed the destination might allready be pressed.
					findCurentOrder()
				}
			}

		case <-time4Update:
			{
				//fmt.Println("drivermaster time for update ")
				myself.LastTimeOnline = time.Now()
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

func findCurentOrder() {
	//find closest order in the same direction as the "curent" order
	activeOrders = make(map[string]commons.OrderStruct)
	activeOrders[curentOrder.ID] = curentOrder
	////myself.CurentDestination = curentOrder.DestinationFloor
	vector := curentOrder.DestinationFloor - myself.CurentFloor
	//fmt.Println("drivermaster curent order ", curentOrder)
	//fmt.Println("drivermaster findcurent order")
	//fmt.Println("drivermaster printing allourorders")
	for key, order := range allOurOrders {
		fmt.Println("	order", order)
		tempV1 := order.DestinationFloor - myself.CurentFloor
		tempV2 := order.Direction
		if (tempV1*vector > 0) && (tempV2*vector > 0) {
			//fmt.Println("	order added to active")
			activeOrders[key] = order
			if tempV1*tempV1 < vector*vector {
				curentOrder = order
			}

		}
	}
	//start doing curent order
	switch curentOrder.Progress {
	case commons.ButtonPressed, commons.Moving2customer, commons.Moving2destination:
		{
			//TODO
			direction := -1
			if myself.CurentFloor < curentOrder.DestinationFloor {
				direction = 1
			}
			setMotorDirection <- direction

			if curentOrder.Progress == commons.ButtonPressed {
				curentOrder.Progress = commons.Moving2customer
			}

			message := commons.MessageStruct{
				SenderID: privateID,
				What:     commons.Order,
				Local:    false,
				Order:    curentOrder,
			}
			privateSendMessege <- message
		}
	case commons.OpeningDoor1, commons.ClosingDoor1, commons.WaitingForDestination, commons.OpeningDoor2, commons.ClosingDoor2:
		{
			//DO  nothing. slave will close door let us now about closing, and then we still cant continue because we need destination.
		}
	default:
		{
			//no curent order
			setMotorDirection <- 0
		}

	}

}
