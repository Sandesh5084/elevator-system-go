package elevator

import "fmt"

type Building struct {
	totalFloors int
	elevator    *Elevator
}

func NewBuilding(totalFloors, elevatorCapacity int) *Building {
	elevator := NewElevator(elevatorCapacity)

	return &Building{
		totalFloors: totalFloors,
		elevator:    elevator,
	}
}

func (b *Building) IsValidFloor(floor int) bool {
	return floor >= 1 && floor <= b.totalFloors
}

func (b *Building) Elevator() *Elevator {
	return b.elevator
}

func (b *Building) RequestElevator(floor int, direction Direction) error {
	if !b.IsValidFloor(floor) {
		return fmt.Errorf("invalid floor %d", floor)
	}

	request := NewRequest(floor, direction, ExternalRequest)
	b.elevator.AddRequest(request)
	return nil
}

func (b *Building) RequestFloor(destination int) error {
	if !b.IsValidFloor(destination) {
		return fmt.Errorf("invalid floor %d", destination)
	}

	request := NewRequest(destination, DirectionIdle, InternalRequest)
	b.elevator.AddRequest(request)
	return nil
}
