package network

import (
	"flag"
	"fmt"
	"os"
	"time"

	"./network/bcast"
	"./network/localip"
	"./network/peers"

	"../elevatorHW"
	"../fsm"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type HelloMsg struct {
	ElevatorID    string // This number identifies the elevator
	CurrentState  int    // This number, says if the elevator is moving up (1) / down (-1) / idle (0)
	LastFloor     int    // The last floor the elevator visited
	Order         OrderMsg
	TimeStamp     int64
	OrderExecuted fsm.Order
}

type OrderMsg struct {
	Order                   fsm.Order
	ElevatorToTakeThisOrder string
}

var OperatingElevators int
var OperatingElevatorsPtr *int

type ElevatorStatus struct {
	ElevatorID string
	Alive      bool
}

func DeleteDeadElevator(operatingElevatorStates map[string]HelloMsg) {
	timeNow := time.Now().Unix()
	for key, value := range operatingElevatorStates {
		if timeNow > value.TimeStamp+1 {
			delete(operatingElevatorStates, key)
		}
	}
}

func UpdateElevatorStates(newMsg HelloMsg, operatingElevatorStates map[string]HelloMsg) {

	operatingElevatorStates[newMsg.ElevatorID] = newMsg
}

func GetLocalID() string {
	localIP, _ := localip.LocalIP()
	return localIP[12:]
}

func NetworkMain(messageCh chan<- HelloMsg, receivedNetworkOrderCh chan<- HelloMsg, sendOrderToPeerCh chan OrderMsg, orderCompletedCh chan fsm.Order, sendDeletedOrderCh chan fsm.Order) {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	var elevatorID string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
		elevatorID = localIP[12:]
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)

	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16167, helloTx)
	go bcast.Receiver(16167, helloRx)
	OperatingElevatorsPtr = &OperatingElevators

	
	go func() {
		initialOrderCompleted := fsm.Order{-1, -1}
		helloMsg := HelloMsg{elevatorID, 0, 5, OrderMsg{fsm.Order{-1, -1}, "Nil"}, 0, initialOrderCompleted}

		for {
			select {
			case deletedOrder := <-sendDeletedOrderCh:
				helloMsg.CurrentState = elevatorHW.GetElevatorState()
				helloMsg.LastFloor = fsm.LatestFloor
				helloMsg.Order = OrderMsg{fsm.Order{-1, -1}, "Nil"}
				helloMsg.OrderExecuted = deletedOrder
				helloTx <- helloMsg
			case order := <-sendOrderToPeerCh:
				helloMsg.CurrentState = elevatorHW.GetElevatorState()
				helloMsg.LastFloor = fsm.LatestFloor
				helloMsg.Order = order
				helloTx <- helloMsg
			case completedOrder := <-orderCompletedCh:
				helloMsg.OrderExecuted = completedOrder
				helloMsg.CurrentState = elevatorHW.GetElevatorState()
				helloMsg.Order = OrderMsg{fsm.Order{-1, -1}, "Nil"}
				helloMsg.LastFloor = fsm.LatestFloor
				helloTx <- helloMsg
			default:
				break
			}
			helloMsg.CurrentState = elevatorHW.GetElevatorState()
			helloMsg.Order = OrderMsg{fsm.Order{-1, -1}, "Nil"}
			helloMsg.LastFloor = fsm.LatestFloor
			helloMsg.OrderExecuted = fsm.Order{-1, -1}

			helloTx <- helloMsg

			time.Sleep(250 * time.Millisecond)
		}
	}()

	//fmt.Println("Started")

	for {
		select {
		case p := <-peerUpdateCh:

			*OperatingElevatorsPtr = len(p.Peers)

		case a := <-helloRx:

			a.TimeStamp = time.Now().Unix()
			messageCh <- a
		}
	}
}
