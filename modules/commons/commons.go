package commons

//Message is used for communication betwen modules and nodes.
//what is "enum" with possible values being kill, spawn, CSE, order. Text is empty except in case of CSE or order.
// In case of CSE
// In case of ordere the following formating applies "d1:d2",
// where d1  is floor of the pressed button and d2 is disered direction.
type Message struct {
	//we dont relly really on messages arriving in corect order.
	//time         time.Time //Time when message was created. On forwarding it should not be changed.
	//id           int //Simpler id
	SenderIP     string
	SenderPROCES string
	What         string
	Text         string
}

//Elevator stores all relevant info.
type Elevator struct {
	IP          string
	Proces      string
	Operational bool
	Floor       int
	Door        int
	Lamp        int
}

//ElevatorToString creates COMMON string representation of elevator.
func ElevatorToString(e Elevator) string {
	return "TODO"
}

//StringToElevator creates elevator from COMMON string representation.
func StringToElevator(s string) Elevator {
	return Elevator{}
}

//OLD COMENTS REPRESENTING OLD IDEAS

//following formating applies "s1:s2" where
// s1 is elevator change applies to
// s2 is string representation of state made with elevatorToString

// d1 is floor of elevator,
// d2 is state of door(open,closed),
// d3 is light(on,off)"
