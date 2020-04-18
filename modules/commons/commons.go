package commons

import "time"

//Arbitrary numbers or something that makes sanse to change
const NumFloors int = 4
const PollRate time.Duration = 20 * time.Millisecond
const RecoverTime = 100 * time.Millisecond      //If we loose conection after how much time we try again
const SendUpdateDelay = 3000 * time.Millisecond //After how much time do we send curent state of our elevator. "im alive message". 500ms for production
const DoorOpenDuratation = 5 * time.Second
const CheckDoorOpen = 100 * time.Millisecond      //used to check if doors are closed so we can move elevator.
const WatchDogFrequency = 3000 * time.Millisecond //After how much time we check if orders are being executed. //500ms for production
const MaxOrderTime = 30 * time.Second             //How much time does the elevator has before we try to find another elevator to do it.
const MaxUserTime = 30 * time.Second              //time we wait for user to press button, before we delete order.
const MaxElevatorTime = -10 * time.Second         //If we dont get CSE for $s time we assume elevator is gone

var ElevatorPort string = "15657"

//Types and enums introduces for typesafety reasons

type What int

const (
	Kill  What = 1 + iota //Not implemented
	Spawn                 //Unncesary since we are not using tcp
	Kick                  //Not implemented
	CSE
	Order
	Malfunction //when we need to tell cseDB that elevator is not doing his job
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
	ID              string
	LastTimeOnline  time.Time
	LastTimeChecked time.Time //Time when we find out it does not work as it should. As such we assume it cant perform any order that are older than this time stamp. We do assume it go instantly fixed and as such it can start doing new orders
	CurentFloor     int
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
	MovingTime       time.Time //time when an elevator started doing this Order. Sometimes order can take long time simply because of high amount of trafic and not because something is wrong
	LastUpdate       time.Time //when we wait for user to press button we need to know when did we close the door.
	Contractor       string    //ID of elevator responsible for this order
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
