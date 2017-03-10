package main

import (
	//"time"

	"./src/elevatorHW"
	"./src/fsm"
	"./src/olasnetwork"
	//"./src/io"
	"fmt"
)

// This function returns how suitet the elevator is to handle a global call
func costFunction(dir int, lastFloor int, order fsm.Order) int {
	var distanceToTarget int
	var dirCost int
	if lastFloor > order.Floor {
		distanceToTarget = lastFloor - order.Floor
	} else {
		distanceToTarget = order.Floor - lastFloor
	}

	distCost := 2 * distanceToTarget

	if dir == 0 && order.Floor == lastFloor { // Elevator is Idle at floor being called
		return 0
	}
	if order.Button == 1 { //UpType Order
		if dir == -1 { // Moving in opposite direction
			if lastFloor < order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		} else if dir == 0 { //Elevator is idle
			dirCost = 0
		} else if dir == 1 { // Elevator is moving up
			if lastFloor > order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		}

	} else if order.Button == 0 { //DownType order
		if dir == 1 { // Oposite directioin
			if lastFloor > order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		} else if dir == 0 {
			dirCost = 0
		} else if dir == -1 {
			if lastFloor < order.Floor {
				dirCost = 8
			} else {
				dirCost = 10
			}
		}
	}
	return dirCost + distCost
}

func decitionmaker(onlineElevatorStates []olasnetwork.HelloMsg) (string, int) {
	numberOfElevatorsInNetwork := olasnetwork.OperatingElevators
	if numberOfElevatorsInNetwork == 0 || numberOfElevatorsInNetwork == 1 {
		return olasnetwork.GetLocalID(), 0
	}
	var costs []int

	lowestCost := 150
	var minPos int
	if len(onlineElevatorStates) < 1 {
		return olasnetwork.GetLocalID(), 0
	}
	for i := 0; i < len(onlineElevatorStates); i++ {
		thisCost := costFunction(onlineElevatorStates[i].CurrentState, onlineElevatorStates[i].LastFloor, onlineElevatorStates[i].Order.Order)
		costs = append(costs, thisCost)
		if lowestCost > costs[i] {
			lowestCost = costs[i]
			minPos = i
		}
	}
	return onlineElevatorStates[minPos].ElevatorID, lowestCost
}

func main() {
	//start init
	fmt.Println("Starting system")
	fmt.Print("\n\n")
	fsm.StartUpMessage()
	elevatorHW.Init()
	//finished init
	fsm.CreateQueueSlice()

	var operatingElevatorStates []olasnetwork.HelloMsg

	buttonCh := make(chan fsm.Order)
	messageCh := make(chan olasnetwork.HelloMsg)
	networkOrderCh := make(chan olasnetwork.HelloMsg)
	networkSendOrderCh := make(chan olasnetwork.OrderMsg)

	go fsm.RunElevator()
	go fsm.GetButtonsPressed(buttonCh)
	go olasnetwork.NetworkMain(messageCh, networkOrderCh, networkSendOrderCh)

	for {
		//fmt.Print("Elevator states online: ")
		fmt.Println(operatingElevatorStates)
		//time.Sleep(1 * time.Second)

		select {

		case newMsg := <-messageCh:
			operatingElevatorStates = olasnetwork.UpdateElevatorStates(newMsg, operatingElevatorStates)
			fmt.Println(operatingElevatorStates)

		case newOrder := <-buttonCh:
			fmt.Print("You made an order: ")
			fmt.Println(newOrder)
			if newOrder.Button == elevatorHW.ButtonCommand {
				fsm.PutInsideOrderInLocalQueue(newOrder)
			} else {

				elevatorToHandleThisOrder, _ := decitionmaker(operatingElevatorStates)
				networkSendOrderCh <- olasnetwork.OrderMsg{newOrder, elevatorToHandleThisOrder}
			}
			fsm.PrintQueues()

		case newNetworkOrder := <-networkOrderCh:
			if newNetworkOrder.Order.ElevatorToTakeThisOrder == olasnetwork.GetLocalID() {
				fsm.PutOrderInLocalQueue(newNetworkOrder.Order.Order)
			}
		}
	}
}
