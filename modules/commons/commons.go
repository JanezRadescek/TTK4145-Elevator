package commons

import "time"

//Arbitrary numbers or something that makes sanse to change
const NumFloors int = 4
const PollRate time.Duration = 20 * time.Millisecond
const DoorOpenDuratation = 5 * time.Second
const CheckDoorOpen = 100 * time.Millisecond //used to check if doors are closed so we can move elevator.

var ElevatorPort string = "15657"

//Types and enums introduces for typesafety reasons

type What int

const (
	Kill  What = 1 + iota //Not implemented
	Spawn                 //Unncesary since we are not using tcp
	Kick                  //Not implemented
	CSE
	Order
)

//MessageStruct is used for communication betwen nodes.
type MessageStruct struct {
	SenderID string //IP:PID
	What     What
	Local    bool //Its convient to use same "road" as if someone else send us CSE.
	Elevator ElevatorStruct
	Order    OrderStruct
}

//ElevatorStruct stores all relevant info.
type ElevatorStruct struct {
	ID             string
	LastTimeOnline time.Time
	Operational    bool
	CurentFloor    int
	//CurentDestination int
	Idle bool
}

type OrderProgress int

//posible values of OrderProgress. Time intervals of orderProgress are open (_,_). behavior on borders is undefined. Undefined bahavior has measure 0 so its OK.
const (
	ButtonPressed OrderProgress = 1 + iota
	Moving2customer
	OpeningDoor1 //Last time to switch elevator. We cant use diferent elevator once doors are closed
	ClosingDoor1
	WaitingForDestination // After this Destination floor changes
	Moving2destination    //6
	OpeningDoor2
	ClosingDoor2
)

//OrderStruct stores all relevant info
type OrderStruct struct {
	ID               string //ElevatorID:number
	Progress         OrderProgress
	Direction        int       //acourding to pressed buttom, 1 = up, -1 = down. primary use to calculate if we want do this order while some other order is in progress.
	DestinationFloor int       //For progress 1-5 it is where the customer is waiting, for 6-8 it is where customer whants to go
	StartingTime     time.Time //time when **buttom** what pressed
	//UpdateTime       time.Time //time of last update
	Contractor string //ID of elevator responsible for this order
}

// type PickButtonStruct struct {
// 	Floor     int
// 	Direction int
// }

type LampStruct struct {
	Floor int
	ON    bool
}

type State int

const (
	Idle State = 1 + iota
	OpeningDoors
	ClosingDoors
	MovingUP
	MovingDown
	Waiting
)
