package elevator

import "testing"

func TestNewElevatorInitialState(t *testing.T) {
	const capacity = 8
	elevator := NewElevator(capacity)

	if elevator.currentFloor != 1 {
		t.Errorf("currentFloor = %d, want 1", elevator.currentFloor)
	}
	if elevator.direction != DirectionIdle {
		t.Errorf("direction = %v, want %v", elevator.direction, DirectionIdle)
	}
	if elevator.state != ElevatorIdle {
		t.Errorf("state = %v, want %v", elevator.state, ElevatorIdle)
	}
	if elevator.door == nil {
		t.Fatal("door = nil, want a closed door")
	}
	if got := elevator.door.State(); got != DoorClosed {
		t.Errorf("door.State() = %v, want %v", got, DoorClosed)
	}
	if elevator.capacity != capacity {
		t.Errorf("capacity = %d, want %d", elevator.capacity, capacity)
	}
	if elevator.passengerCount != 0 {
		t.Errorf("passengerCount = %d, want 0", elevator.passengerCount)
	}
	if elevator.scheduler == nil {
		t.Fatal("scheduler = nil, want initialized scheduler")
	}
}

func TestElevatorReadOnlyGetters(t *testing.T) {
	elevator := NewElevator(8)

	if got := elevator.CurrentFloor(); got != 1 {
		t.Errorf("CurrentFloor() = %d, want 1", got)
	}
	if got := elevator.Direction(); got != DirectionIdle {
		t.Errorf("Direction() = %v, want %v", got, DirectionIdle)
	}
	if got := elevator.State(); got != ElevatorIdle {
		t.Errorf("State() = %v, want %v", got, ElevatorIdle)
	}

	if err := elevator.MoveTo(3); err != nil {
		t.Fatalf("MoveTo(3) returned error: %v", err)
	}
	if got := elevator.CurrentFloor(); got != 3 {
		t.Errorf("CurrentFloor() after movement = %d, want 3", got)
	}
	if got := elevator.Direction(); got != DirectionIdle {
		t.Errorf("Direction() after movement = %v, want %v", got, DirectionIdle)
	}
	if got := elevator.State(); got != ElevatorIdle {
		t.Errorf("State() after movement = %v, want %v", got, ElevatorIdle)
	}
}

func TestElevatorServesNextRequest(t *testing.T) {
	elevator := NewElevator(8)
	request := NewRequest(5, DirectionUp, ExternalRequest)
	elevator.AddRequest(request)

	if err := elevator.ServeNextRequest(); err != nil {
		t.Fatalf("ServeNextRequest() returned error: %v", err)
	}
	if elevator.currentFloor != 5 {
		t.Errorf("currentFloor = %d, want 5", elevator.currentFloor)
	}
	if request.Status != RequestCompleted {
		t.Errorf("request.Status = %v, want %v", request.Status, RequestCompleted)
	}
}

func TestElevatorServeNextRequestKeepsRequestPendingOnMoveError(t *testing.T) {
	elevator := NewElevator(8)
	request := NewRequest(5, DirectionUp, ExternalRequest)
	elevator.AddRequest(request)
	elevator.door.Open()

	if err := elevator.ServeNextRequest(); err == nil {
		t.Fatal("ServeNextRequest() returned no error while door was open")
	}
	if request.Status != RequestPending {
		t.Errorf("request.Status = %v, want %v after failed move", request.Status, RequestPending)
	}
}

func TestElevatorServeNextRequestDoesNothingWhenNoPendingRequests(t *testing.T) {
	elevator := NewElevator(8)

	if err := elevator.ServeNextRequest(); err != nil {
		t.Fatalf("ServeNextRequest() with no requests returned error: %v", err)
	}
	if elevator.currentFloor != 1 {
		t.Errorf("currentFloor = %d, want 1", elevator.currentFloor)
	}
}

func TestElevatorMoveToHigherFloor(t *testing.T) {
	elevator := NewElevator(8)
	elevator.currentFloor = 2

	if err := elevator.MoveTo(5); err != nil {
		t.Fatalf("MoveTo(5) returned error: %v", err)
	}
	if elevator.currentFloor != 5 {
		t.Errorf("currentFloor = %d, want 5", elevator.currentFloor)
	}
	if elevator.direction != DirectionIdle {
		t.Errorf("direction = %v, want %v after movement", elevator.direction, DirectionIdle)
	}
	if elevator.state != ElevatorIdle {
		t.Errorf("state = %v, want %v after movement", elevator.state, ElevatorIdle)
	}
}

func TestElevatorMoveToLowerFloor(t *testing.T) {
	elevator := NewElevator(8)
	elevator.currentFloor = 5

	if err := elevator.MoveTo(2); err != nil {
		t.Fatalf("MoveTo(2) returned error: %v", err)
	}
	if elevator.currentFloor != 2 {
		t.Errorf("currentFloor = %d, want 2", elevator.currentFloor)
	}
	if elevator.direction != DirectionIdle {
		t.Errorf("direction = %v, want %v after movement", elevator.direction, DirectionIdle)
	}
	if elevator.state != ElevatorIdle {
		t.Errorf("state = %v, want %v after movement", elevator.state, ElevatorIdle)
	}
}

func TestElevatorMoveToCurrentFloor(t *testing.T) {
	elevator := NewElevator(8)
	elevator.currentFloor = 2

	if err := elevator.MoveTo(2); err != nil {
		t.Fatalf("MoveTo(2) returned error: %v", err)
	}
	if elevator.currentFloor != 2 {
		t.Errorf("currentFloor = %d, want 2", elevator.currentFloor)
	}
	if elevator.direction != DirectionIdle {
		t.Errorf("direction = %v, want %v", elevator.direction, DirectionIdle)
	}
	if elevator.state != ElevatorIdle {
		t.Errorf("state = %v, want %v", elevator.state, ElevatorIdle)
	}
}

func TestElevatorMoveToWhileDoorOpen(t *testing.T) {
	elevator := NewElevator(8)
	elevator.currentFloor = 2
	elevator.door.Open()

	if err := elevator.MoveTo(5); err == nil {
		t.Fatal("MoveTo(5) returned no error while door was open")
	}
	if elevator.currentFloor != 2 {
		t.Errorf("currentFloor = %d, want 2 after rejected movement", elevator.currentFloor)
	}
	if elevator.state != ElevatorIdle {
		t.Errorf("state = %v, want %v after rejected movement", elevator.state, ElevatorIdle)
	}
	if elevator.direction != DirectionIdle {
		t.Errorf("direction = %v, want %v after rejected movement", elevator.direction, DirectionIdle)
	}
}
