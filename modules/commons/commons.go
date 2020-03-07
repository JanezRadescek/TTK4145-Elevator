package commons

import "time"

//What is enum
type What int

//posible values of what
const (
	Kill  What = 1 + iota //Not implemented either
	Spawn                 //Unncesary since we are not using tcp
	Kick                  //Not implemented
	CSE
	Order
)

//MessageStruct is used for communication betwen nodes.
type MessageStruct struct {
	//we dont really rely on messages arriving in corect order.
	SenderID string //IP:PID
	What     What
	Local    bool //Its convient to use same "road" as if someone else send us CSE.
	Elevator ElevatorStruct
	Order    OrderStruct
}

//ElevatorStruct stores all relevant info.
type ElevatorStruct struct {
	ID                string
	LastTimeOnline    time.Time
	Operational       bool
	CurentFloor       int
	CurentDestination int
	Idle              bool
}

//OrderProgress is enum
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
	Direction        int       //acourding to pressed buttom, 1 = up, -1 = down, it doesnt make sanse for order to have "stop" direction, thats elevators problem
	DestinationFloor int       //For progress 1-5 it is where the customer is waiting, for 6-8 it is where customer whants to go
	StartingTime     time.Time //time when **buttom** what pressed
	//UpdateTime       time.Time //time of last update
	Contractor string //ID of elevator responsible for this order
}

type PickButtonStruct struct {
	Floor     int
	Direction int
}

type LampStruct struct {
	Floor int
	ON    bool
}
