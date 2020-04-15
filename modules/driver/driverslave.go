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
	setLamp <-chan commons.LampStruct,
	setDoor <-chan bool,
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
			}
		case f := <-drv_floors:
			{
				floorSensor <- f
			}
		case door := <-setDoor:
			{
				elevio.SetMotorDirection(elevio.MD_Stop)
				motorDirection = elevio.MD_Stop

				//Is this how you open doors for customers to get in/out ??
				elevio.SetDoorOpenLamp(true)

				doorSensor <- true
				time.Sleep(commons.DoorOpenDuratation)

				elevio.SetDoorOpenLamp(false)
				doorSensor <- false
			}
		}
	}
}
