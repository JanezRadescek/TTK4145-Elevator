package commons

import "time"

//What is enum
type What int

//posible values of what
const (
	Kill What = 1 + iota
	Spawn
	Kick
	CSE
	Order
)

//MessageStruct is used for communication betwen modules and nodes.
//what is "enum" with possible values being kill, spawn, CSE, order. Text is empty except in case of CSE or order.
// In case of CSE it is marshalled struct Elevator
// In case of ordere the following formating applies "d1:d2",
// where d1  is floor of the pressed button and d2 is disered direction.
type MessageStruct struct {
	//we dont relly really on messages arriving in corect order.
	//time         time.Time //Time when message was created. On forwarding it should not be changed.
	//id           int //Simpler id
	SenderIP        string
	SenderProcessID string
	What            What
	Local           bool //Its convient to use same "road" as if someone else send us CSE.
	Elevator        ElevatorStruct
	Order           OrderStruct
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

//posible values of OrderProgress
const (
	ButtonPressed OrderProgress = 1 + iota
	Moving2customer
	OpeningDoor1
	ClosingDoor1
	WaitingForDestination
	Moving2destination
	OpeningDoor2
	ClosingDoor2
)

//OrderStruct stores all relevant info
type OrderStruct struct {
	ID               string //globalOrderID ElevatorID:number
	Progress         OrderProgress
	Direction        int       //acourding to pressed buttom
	DestinationFloor int       //For progress 1-2 it iswhere customer wait, for 6-8 it is where customer whants to go
	StartingTime     time.Time //time when **buttom** what pressed
	UpdateTime       time.Time //time of last update
	Contractor       string    //ID of elevator responsible for this order
}
