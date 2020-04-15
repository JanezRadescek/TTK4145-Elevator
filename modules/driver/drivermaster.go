package driver

import (
	"strconv"
	"time"

	"../commons"
	"./driver-go/elevio"
)

const sendUpdateDelay = 500 * time.Millisecond

var myself commons.ElevatorStruct
var curentOrder commons.OrderStruct
var allOurOrders map[string]commons.OrderStruct
var activeOrders map[string]commons.OrderStruct

//TODO properly react to "disruptor" presing buttons at rendom times. (like what happens if somebody somehow pushes button for floor 5 before we even let him in?)

//StartDriverMaster takes next order we are asigned and and give high level instruction on what to do with it.
func StartDriverMaster(
	ID string,
	reciveCopy <-chan map[string]commons.OrderStruct,
	sendMessege chan<- commons.MessageStruct,
) {

	myself = commons.ElevatorStruct{}
	curentOrder = commons.OrderStruct{}
	//key is?
	allOurOrders = make(map[string]commons.OrderStruct)
	activeOrders = make(map[string]commons.OrderStruct)
	oldestTime := time.Now()
	IDcounter := 1

	time4Update := make(chan bool)
	go func() {
		for {
			time.Sleep(sendUpdateDelay)
			time4Update <- true
		}

	}()

	newButton := make(chan elevio.ButtonEvent)
	floorSensor := make(chan int)
	doorSensor := make(chan bool)
	//stopButton := make(chan bool) //solve this on IO level
	setMotorDirection := make(chan int)
	setLamp := make(chan commons.LampStruct)
	setDoor := make(chan bool)

	go StartDriverSlave(newButton, floorSensor, doorSensor, setMotorDirection, setLamp, setDoor)

	for {
		select {
		case allOurOrders = <-reciveCopy:
			{
				//find the oldest
				for _, order := range allOurOrders {
					if order.StartingTime.Before(oldestTime) {
						oldestTime = order.StartingTime
						curentOrder = order
					}
				}
				//find other order we may do on the way to the oldest first and start doing it.
				findCurentOrder()

			}

		case button := <-newButton:
			{
				switch button.Button {
				case elevio.BT_Cab:
					{
						floor := button.Floor
						newOrder := true //two customers might wanna go to 2 diferent floor
						for _, order := range activeOrders {
							if order.DestinationFloor == floor {
								newOrder = false
							}
							//update orders
							if order.Progress == commons.WaitingForDestination {
								order.Progress = commons.Moving2destination
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

							//Hole project is design as if we have floorbuttons. We only have buttons down and up.
							//As such we are hecking here a little bit.
							//In case that its in different direction, we give it curent time so that older orders can get executed first.(preventing someone hijacking elevator)
							//If we get new order in the same direction we give it the same time as curent order since it might actually be curent order.
							tempT := time.Now()

							order := commons.OrderStruct{
								ID:               ID + ":" + strconv.Itoa(IDcounter),
								Progress:         commons.Moving2destination,
								Direction:        tempD,
								DestinationFloor: floor,
								StartingTime:     tempT,
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
				case elevio.BT_HallUp, elevio.BT_HallDown:
					{
						floor := button.Floor
						direction := 1 //up
						if button.Button == elevio.BT_HallDown {
							direction = -1
						}
						order := commons.OrderStruct{
							ID:               ID + ":" + strconv.Itoa(IDcounter),
							Progress:         commons.ButtonPressed,
							Direction:        direction,
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
						setLamp <- commons.LampStruct{Floor: floor, ON: true}
					}
				}

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
				setDoor <- openDoor
			}

		case door := <-doorSensor:
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
						//after closing door find curent order
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

func findCurentOrder() {
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
		//start doing curent order
		switch curentOrder.Progress {
		case commons.ButtonPressed, commons.Moving2customer:
			{
				//TODO
			}
		case commons.OpeningDoor1, commons.ClosingDoor1, commons.WaitingForDestination:
			{
				//DO  nothing. slaver will close door let us now about closing, and then we still cant continue because we need destination.
			}
		case commons.Moving2destination:
			{
				//TODO
			}
		case commons.OpeningDoor2, commons.ClosingDoor2:
			{
				//DO  nothing again. let slave do his work.
			}

		}
	}
}
