package driver

import (
	"fmt"
	"time"

	"../commons"
	"./driver-go/elevio"
)

var destination int
var curentFloor int
var motorDirection elevio.MotorDirection

//StartDriverSlave operates elevator acording to the high level orders and report back sensor readings
func StartDriverSlave(
	newButton chan<- elevio.ButtonEvent,
	floorSensor chan<- int,
	doorSensor chan<- bool, //true if open
	getDestination <-chan int,
	setOpenDoor <-chan bool,
) {
	elevatorPort := commons.ElevatorPort
	numFloors := commons.NumFloors

	defer func() {
		if r := recover(); r != nil {
			time.Sleep(commons.RecoverTime)
			fmt.Println()
			fmt.Println("slave is restarting. probably because of the loss of connection.")
			StartDriverSlave(newButton, floorSensor, doorSensor, getDestination, setOpenDoor)
		}
	}()
	elevio.Init("localhost:"+elevatorPort, numFloors)

	//move down so if we are inbetwen floors we can get floor sensor reading and get to know where we are
	motorDirection = elevio.MD_Down
	elevio.SetMotorDirection(motorDirection)

	var doorOpen bool = false
	destination = 0
	curentFloor = 0

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	for {
		select {
		case b := <-drv_buttons:
			{
				newButton <- b
				elevio.SetButtonLamp(b.Button, b.Floor, true)
			}
		case f := <-drv_floors:
			{
				curentFloor = f
				elevio.SetFloorIndicator(f)
				elevio.SetMotorDirection(elevio.MD_Stop)
				floorSensor <- f
			}
		case o := <-drv_obstr:
			{
				if o {
					go elevio.SetMotorDirection(elevio.MD_Stop)
				} else {
					go elevio.SetMotorDirection(motorDirection)
				}
			}
		case s := <-drv_stop:
			{
				if s {
					go elevio.SetMotorDirection(elevio.MD_Stop)
				} else {
					go elevio.SetMotorDirection(motorDirection)
				}
			}
		case <-setOpenDoor:
			{
				if !doorOpen {
					doorOpen = true
					go func() {
						elevio.SetMotorDirection(elevio.MD_Stop)

						elevio.SetDoorOpenLamp(true)

						elevio.SetButtonLamp(elevio.BT_HallUp, curentFloor, false)
						elevio.SetButtonLamp(elevio.BT_HallDown, curentFloor, false)
						elevio.SetButtonLamp(elevio.BT_Cab, curentFloor, false)

						doorSensor <- true
						time.Sleep(commons.DoorOpenDuratation)

						elevio.SetDoorOpenLamp(false)
						doorSensor <- false
						doorOpen = false
					}()
				}
			}
		case d := <-getDestination:
			{
				destination = d
				calculateDirection()
				//we dont want to stop in betwen floors. we will change(if needed) motor direction when we arive at the floor.
				if destination != curentFloor {

					go func() {
						for {
							if doorOpen {
								//we wait for door to close
								time.Sleep(commons.CheckDoorOpen)
							} else {
								elevio.SetMotorDirection(motorDirection)
								break
							}
						}
					}()
				}

			}

		}
	}
}

func calculateDirection() {
	if destination > curentFloor {
		motorDirection = elevio.MD_Up
	} else if destination == curentFloor {
		motorDirection = elevio.MD_Stop
	} else {
		motorDirection = elevio.MD_Down
	}
}
