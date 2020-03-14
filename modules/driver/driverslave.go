package driver

import (
	"../commons"
	"./driver-go/elevio"
)

//StartDriverSlave operates elevator acording to the high level orders and report back sensor readings
func StartDriverSlave(
	pickButton chan<- int,
	floorButton chan<- int,
	floorSensor chan<- int,
	DoorSensor chan<- bool,
	setMotorDirection <-chan int,
	setLamp <-chan commons.LampStruct,
	SetDoor <-chan bool,
) {
	elevatorPort := commons.ElevatorPort
	numFloors := commons.NumFloors

	elevio.Init("localhost:"+elevatorPort, numFloors)

	var d elevio.MotorDirection = elevio.MD_Stop
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
				switch b.Button {
				case elevio.BT_Cab:
					{
						pickButton <- b.Floor
					}
				case elevio.BT_HallUp:
					{
						floorButton <- (curentFloor + 1)
					}
				case elevio.BT_HallDown:
					{
						floorButton <- (curentFloor - 1)
					}
				}
			}
		case f := <-drv_floors:
			{
				floorSensor <- f
			}
		}
	}
}
