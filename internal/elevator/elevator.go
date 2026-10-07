package elevator

import (
	"errors"
)

type ElevatorState int

const (
	ElevatorIdle ElevatorState = iota
	ElevatorMoving
)

type Elevator struct {
	currentFloor   int
	direction      Direction
	state          ElevatorState
	door           *Door
	capacity       int
	passengerCount int
	scheduler      *Scheduler
}

func NewElevator(capacity int) *Elevator {
	return &Elevator{
		currentFloor:   1,
		direction:      DirectionIdle,
		state:          ElevatorIdle,
		door:           NewDoor(),
		capacity:       capacity,
		passengerCount: 0,
		scheduler:      NewScheduler(),
	}
}

func (e *Elevator) CurrentFloor() int {
	return e.currentFloor
}

func (e *Elevator) Direction() Direction {
	return e.direction
}

func (e *Elevator) State() ElevatorState {
	return e.state
}

func (e *Elevator) AddRequest(request *Request) {
	e.scheduler.AddRequest(request)
}

func (e *Elevator) ServeNextRequest() error {
	request := e.scheduler.NextRequest(e)
	if request == nil {
		return nil
	}

	if err := e.MoveTo(request.Floor); err != nil {
		return err
	}

	request.Status = RequestCompleted
	return nil
}

func (e *Elevator) MoveTo(destination int) error {
	if e.door.State() == DoorOpen {
		return errors.New("cannot move while door is open")
	}

	if destination == e.currentFloor {
		return nil
	}

	e.state = ElevatorMoving
	if destination > e.currentFloor {
		e.direction = DirectionUp
		for e.currentFloor < destination {
			e.currentFloor++
		}
	} else {
		e.direction = DirectionDown
		for e.currentFloor > destination {
			e.currentFloor--
		}
	}

	e.direction = DirectionIdle
	e.state = ElevatorIdle
	return nil
}
