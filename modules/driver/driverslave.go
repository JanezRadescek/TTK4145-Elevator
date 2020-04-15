package driver

import (
	"time"

	"../commons"
	"./driver-go/elevio"
)

//StartDriverSlave operates elevator acording to the high level orders and report back sensor readings
func StartDriverSlave(
	newButton chan<- elevio.ButtonEvent,
	floorSensor chan<- int,
	doorSensor chan<- bool, //true if open
	setMotorDirection <-chan int,
	setOpenDoor <-chan bool,
) {
	elevatorPort := commons.ElevatorPort
	numFloors := commons.NumFloors

	elevio.Init("localhost:"+elevatorPort, numFloors)

	var motorDirection elevio.MotorDirection = elevio.MD_Stop
	var doorOpen bool = false
	curentFloor := 0

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
				go elevio.SetButtonLamp(b.Button, b.Floor, true)
			}
		case f := <-drv_floors:
			{
				floorSensor <- f
				curentFloor = f
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
					go func() {
						doorOpen = true
						elevio.SetMotorDirection(elevio.MD_Stop)

						//Is this how you open doors for customers to get in/out ??
						elevio.SetDoorOpenLamp(true)

						elevio.SetButtonLamp(elevio.BT_HallUp, curentFloor, false)
						elevio.SetButtonLamp(elevio.BT_HallDown, curentFloor, false)
						elevio.SetButtonLamp(elevio.BT_Cab, curentFloor, false)

						doorSensor <- true
						time.Sleep(commons.DoorOpenDuratation)

						//Is this how you open doors for customers to get in/out ??
						elevio.SetDoorOpenLamp(false)
						doorSensor <- false
						doorOpen = false
					}()
				}
			}
		case d := <-setMotorDirection:
			{
				//TODO
			}

		}
	}
}
