package driver

import (
	"fmt"
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

	defer func() {
		if r := recover(); r != nil {
			time.Sleep(commons.RecoverTime)
			fmt.Println()
			fmt.Println("slave is restarting. probably because of the loss of connection.")
			StartDriverSlave(newButton, floorSensor, doorSensor, setMotorDirection, setOpenDoor)
		}
	}()
	elevio.Init("localhost:"+elevatorPort, numFloors)

	//move down so if we are inbetwen floors we can get floor sensor reading and get to know where we are
	var motorDirection elevio.MotorDirection = elevio.MD_Down
	elevio.SetMotorDirection(motorDirection)

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
				fmt.Println("driverslave button ", b)
				newButton <- b
				elevio.SetButtonLamp(b.Button, b.Floor, true)
			}
		case f := <-drv_floors:
			{
				fmt.Println("driverslave floor ", f)
				floorSensor <- f
				curentFloor = f
				elevio.SetFloorIndicator(f)
			}
		case o := <-drv_obstr:
			{
				fmt.Println("driverslave obstacle ", o)
				if o {
					go elevio.SetMotorDirection(elevio.MD_Stop)
				} else {
					go elevio.SetMotorDirection(motorDirection)
				}
			}
		case s := <-drv_stop:
			{
				fmt.Println("driverslave stop ", s)
				if s {
					go elevio.SetMotorDirection(elevio.MD_Stop)
				} else {
					go elevio.SetMotorDirection(motorDirection)
				}
			}
		case <-setOpenDoor:
			{
				fmt.Println("driverslave setopendoor ", true)
				if !doorOpen {
					doorOpen = true
					go func() {
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
				fmt.Println("driverslave set motor direction ", d)
				switch d {
				case 1:
					{
						motorDirection = elevio.MD_Up
					}
				case -1:
					{
						motorDirection = elevio.MD_Down
					}
				case 0:
					{
						motorDirection = elevio.MD_Stop
					}
				}

				go func() {
					for {
						if doorOpen {
							time.Sleep(commons.CheckDoorOpen)
						} else if curentFloor == 0 && motorDirection == elevio.MD_Down {
							motorDirection = elevio.MD_Stop
							elevio.SetMotorDirection(motorDirection)
							break
						} else if curentFloor == 4 && motorDirection == elevio.MD_Up{
							motorDirection = elevio.MD_Stop
							elevio.SetMotorDirection(motorDirection)
							break
						}else {
							fmt.Println("driverslave sending to io direction ", d)
							elevio.SetMotorDirection(motorDirection)
							break
						}
					}
				}()
			}

		}
	}
}
