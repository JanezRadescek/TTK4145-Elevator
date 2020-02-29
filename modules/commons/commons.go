package commons

import "time"

const (
	//Kill enum for What
	Kill = "KILL"
	//Spawn enum for What
	Spawn = "SPAWN"
	//Cse enum for What
	Cse = "CSE"
	//Local enum for What
	Local = "LOCAL"
)

//Message is used for communication betwen modules and nodes.
//what is "enum" with possible values being kill, spawn, CSE, order. Text is empty except in case of CSE or order.
// In case of CSE it is marshalled struct Elevator
// In case of ordere the following formating applies "d1:d2",
// where d1  is floor of the pressed button and d2 is disered direction.
type Message struct {
	//we dont relly really on messages arriving in corect order.
	//time         time.Time //Time when message was created. On forwarding it should not be changed.
	//id           int //Simpler id
	SenderIP        string
	SenderProcessID string
	What            string
	Elevator        Elevator
	Order           Order
}

//Elevator stores all relevant info.
type Elevator struct {
	//IP             string
	//Proces         string
	LastTimeOnline time.Time
	Operational    bool
	Floor          int
	Door           int
	Lamp           int
}

//Order stores all relevant info
type Order struct {
	Floor      int
	Direction  int
	Contractor string //ID of elevator responsible for this order
}
