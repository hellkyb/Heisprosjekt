package elevatorHW

import (
	"fmt"

	"../io"
)

const motorSpeed = 2800
const NFloors = 4
const NButtons = 3
const NLights = 4

type MotorDirection int

const (
	DirectionDown = -1
	DirectionStop = 0
	DirectionUp   = 1
)

type ButtonType int

const (
	ButtonCallDown = iota
	ButtonCallUp
	ButtonCommand
)

var lights [NFloors][NLights]int = [NFloors][NLights]int{
	{lightDown1, lightUp1, lightCommand1},
	{lightDown2, lightUp2, lightCommand2},
	{lightDown3, lightUp3, lightCommand3},
	{lightDown4, lightUp4, lightCommand4}}

func Init() {
	status := io.Init()
	if status {
		fmt.Println("Hardware initialization executed successfully")
	} else {
		fmt.Println("Initialization error")
	}

	for floor := 0; floor < NFloors; floor++ {
		for button := 0; button < NButtons; button++ {
			io.ClearBit(lights[floor][button])
		}
	}
	io.ClearBit(lightDoorOpen)

	SetMotor(DirectionDown)
InitLoop:
	for {
		floor := GetFloorSensorSignal()
		switch floor {
		case 1:
			SetMotor(DirectionStop)
			SetFloorIndicator(1)
			fmt.Println("Arrived at floor", floor)
			break InitLoop
		case 2:
			SetMotor(DirectionStop)
			SetFloorIndicator(2)
			fmt.Println("Arrived at floor", floor)
			break InitLoop
		case 3:
			SetMotor(DirectionStop)
			SetFloorIndicator(3)
			fmt.Println("Arrived at floor", floor)
			break InitLoop
		case 4:
			SetMotor(DirectionStop)
			SetFloorIndicator(4)
			fmt.Println("Arrived at floor", floor)
			break InitLoop
		}
	}
	io.ClearBit(lightStop)
}

func SecondInit() {

	for floor := 0; floor < NFloors; floor++ {
		for button := 0; button < NButtons; button++ {
			io.ClearBit(lights[floor][button])
		}
	}

	SetMotor(DirectionDown)
InitLoop:
	for {
		floor := GetFloorSensorSignal()
		switch floor {
		case 1:
			SetMotor(DirectionStop)
			SetFloorIndicator(1)
			fmt.Println("Arrived at floor", floor)
			break InitLoop
		case 2:
			SetMotor(DirectionStop)
			SetFloorIndicator(2)
			fmt.Println("Arrived at floor", floor)
			break InitLoop
		case 3:
			SetMotor(DirectionStop)
			SetFloorIndicator(3)
			fmt.Println("Arrived at floor", floor)
			break InitLoop
		case 4:
			SetMotor(DirectionStop)
			SetFloorIndicator(4)
			fmt.Println("Arrived at floor", floor)
			break InitLoop
		}
	}
	io.ClearBit(lightStop)
}

func SetDoorLight(open bool) {
	if open {
		io.SetBit(lightDoorOpen)
	} else {
		io.ClearBit(lightDoorOpen)
	}
}

func SetMotor(dir int) {
	if dir == DirectionStop {
		io.WriteAnalog(motor, 0)
	} else if dir == DirectionDown {
		io.SetBit(Motordir)
		io.WriteAnalog(motor, motorSpeed)
	} else if dir == DirectionUp {
		io.ClearBit(Motordir)
		io.WriteAnalog(motor, motorSpeed)
	}
}

func SetFloorIndicator(floor int) {
	switch floor {
	case 1:
		io.ClearBit(lightFloorInd1)
		io.ClearBit(lightFloorInd2)
	case 2:
		io.SetBit(lightFloorInd2)
		io.ClearBit(lightFloorInd1)
	case 3:
		io.ClearBit(lightFloorInd2)
		io.SetBit(lightFloorInd1)
	case 4:
		io.SetBit(lightFloorInd1)
		io.SetBit(lightFloorInd2)
	}
}

func GetFloorSensorSignal() int {
	if io.ReadBit(sensorFloor1) == 1 {
		return 1
	} else if io.ReadBit(sensorFloor2) == 1 {
		return 2
	} else if io.ReadBit(sensorFloor3) == 1 {
		return 3
	} else if io.ReadBit(sensorFloor4) == 1 {
		return 4
	} else {
		return 0
	}
}

func GetUpButton() int {
	if io.ReadBit(buttonUp1) == 1 {
		return 1
	} else if io.ReadBit(buttonUp2) == 1 {
		return 2
	} else if io.ReadBit(buttonUp3) == 1 {
		return 3
	} else {
		return 0
	}
}

func GetDownButton() int {
	if io.ReadBit(buttonDown2) == 1 {
		return 2
	} else if io.ReadBit(buttonDown3) == 1 {
		return 3
	} else if io.ReadBit(buttonDown4) == 1 {
		return 4
	} else {
		return 0
	}
}

func GetInsideElevatorButton() int {
	if io.ReadBit(buttonCommand1) == 1 {
		return 1
	} else if io.ReadBit(buttonCommand2) == 1 {
		return 2
	} else if io.ReadBit(buttonCommand3) == 1 {
		return 3
	} else if io.ReadBit(buttonCommand4) == 1 {
		return 4
	} else {
		return 0
	}
}

func SetDownLight(floor int, onOff bool) {
	if onOff {
		switch floor {
		case 2:
			io.SetBit(lightDown2)
		case 3:
			io.SetBit(lightDown3)
		case 4:
			io.SetBit(lightDown4)
		}
	} else if !onOff {
		switch floor {
		case 2:
			io.ClearBit(lightDown2)
		case 3:
			io.ClearBit(lightDown3)
		case 4:
			io.ClearBit(lightDown4)
		}
	}
}

func SetUpLight(floor int, onOff bool) {
	if onOff {
		switch floor {
		case 1:
			io.SetBit(lightUp1)
		case 2:
			io.SetBit(lightUp2)
		case 3:
			io.SetBit(lightUp3)
		}
	} else if !onOff {
		switch floor {
		case 1:
			io.ClearBit(lightUp1)
		case 2:
			io.ClearBit(lightUp2)
		case 3:
			io.ClearBit(lightUp3)
		}
	}
}

func SetInsideLight(floor int, onOff bool) {
	if onOff {
		switch floor {
		case 1:
			io.SetBit(lightCommand1)
		case 2:
			io.SetBit(lightCommand2)
		case 3:
			io.SetBit(lightCommand3)
		case 4:
			io.SetBit(lightCommand4)
		}
	} else if !onOff {
		switch floor {
		case 1:
			io.ClearBit(lightCommand1)
		case 2:
			io.ClearBit(lightCommand2)
		case 3:
			io.ClearBit(lightCommand3)
		case 4:
			io.ClearBit(lightCommand4)
		}
	}
}

func GetElevatorDirection() int {
	return io.ReadBit(Motordir)
}

func GetDoorLight() int {
	return io.ReadBit(lightDoorOpen)
}

func GetStopButtonPressed() bool {
	if io.ReadAnalog(stop) == 1 {
		return true
	} else {
		return false
	}
}

func GetElevatorState() int {
	if io.ReadAnalog(motor) == 0 {
		return 0
	} else if io.ReadAnalog(motor) != 0 && GetElevatorDirection() == 0 {
		return 1
	} else {
		return -1
	}
}

func SetStopButton(onOff bool) {
	if onOff {
		io.SetBit(lightStop)
	} else {
		io.ClearBit(lightStop)
	}
}
