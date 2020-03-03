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
	//IP             string
	//Proces         string
	LastTimeOnline    time.Time
	Operational       bool
	CurentFloor       int
	CurentDestination int
	Door              int
	Lamp              int
}

//OrderProgress is enum
type OrderProgress int

//posible values of OrderProgress
const (
	ButtomPressed OrderProgress = 1 + iota
	Moving2customer
	OpeningDoor1
	ClosingDoor1
	DestinationChosen
	Moving2destination
	OpeningDoor2
	ClosingDoor2
)

//OrderStruct stores all relevant info
type OrderStruct struct {
	ID               int //only used locally
	Progress         OrderProgress
	CustomerFloor    int
	Direction        int //acourding to pressed buttom
	DestinationFloor int
	StartingTime     time.Time //time when **buttom** what pressed
	UpdateTime       time.Time //time of last update
	Contractor       string    //ID of elevator responsible for this order
}
