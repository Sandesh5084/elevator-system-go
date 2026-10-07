package elevator

import (
	"fmt"
	"testing"
)

func TestBuildingIsValidFloor(t *testing.T) {
	building := NewBuilding(10, 5)

	tests := []struct {
		floor int
		want  bool
	}{
		{floor: 1, want: true},
		{floor: 10, want: true},
		{floor: 0, want: false},
		{floor: 11, want: false},
		{floor: -1, want: false},
	}

	for _, tt := range tests {
		if got := building.IsValidFloor(tt.floor); got != tt.want {
			t.Errorf("IsValidFloor(%d) = %t, want %t", tt.floor, got, tt.want)
		}
	}
}

func TestBuildingOwnsElevator(t *testing.T) {
	building := NewBuilding(10, 5)

	got := building.Elevator()
	if got == nil {
		t.Fatal("Elevator() = nil, want initialized elevator")
	}
	if got.capacity != 5 {
		t.Errorf("elevator capacity = %d, want 5", got.capacity)
	}
	if got.door == nil || got.scheduler == nil {
		t.Error("building elevator is missing its door or scheduler")
	}
}

func TestBuildingRequestElevatorAddsExternalRequest(t *testing.T) {
	building := NewBuilding(10, 5)

	if err := building.RequestElevator(3, DirectionUp); err != nil {
		t.Fatalf("RequestElevator() returned error: %v", err)
	}
	if building.Elevator().currentFloor != 1 {
		t.Errorf("current floor = %d, want 1; submitting a request should not serve it", building.Elevator().currentFloor)
	}

	requests := building.Elevator().scheduler.requests
	if len(requests) != 1 {
		t.Fatalf("pending request count = %d, want 1", len(requests))
	}

	request := requests[0]
	if request.Floor != 3 {
		t.Errorf("request floor = %d, want 3", request.Floor)
	}
	if request.Direction != DirectionUp {
		t.Errorf("request direction = %v, want %v", request.Direction, DirectionUp)
	}
	if request.Type != ExternalRequest {
		t.Errorf("request type = %v, want %v", request.Type, ExternalRequest)
	}
	if request.Status != RequestPending {
		t.Errorf("request status = %v, want %v", request.Status, RequestPending)
	}
}

func TestBuildingRejectsInvalidExternalRequestFloor(t *testing.T) {
	for _, floor := range []int{0, 11, -1} {
		t.Run(fmt.Sprintf("floor_%d", floor), func(t *testing.T) {
			building := NewBuilding(10, 5)

			if err := building.RequestElevator(floor, DirectionUp); err == nil {
				t.Fatalf("RequestElevator(%d) returned no error for invalid floor", floor)
			}
			if got := len(building.Elevator().scheduler.requests); got != 0 {
				t.Errorf("pending request count = %d, want 0", got)
			}
		})
	}
}

func TestBuildingRequestFloorAddsInternalRequest(t *testing.T) {
	building := NewBuilding(10, 5)

	if err := building.RequestFloor(7); err != nil {
		t.Fatalf("RequestFloor(7) returned error: %v", err)
	}
	if building.Elevator().currentFloor != 1 {
		t.Errorf("current floor = %d, want 1; submitting a request should not serve it", building.Elevator().currentFloor)
	}

	requests := building.Elevator().scheduler.requests
	if len(requests) != 1 {
		t.Fatalf("request count = %d, want 1", len(requests))
	}

	request := requests[0]
	if request.Floor != 7 {
		t.Errorf("request floor = %d, want 7", request.Floor)
	}
	if request.Direction != DirectionIdle {
		t.Errorf("request direction = %v, want %v", request.Direction, DirectionIdle)
	}
	if request.Type != InternalRequest {
		t.Errorf("request type = %v, want %v", request.Type, InternalRequest)
	}
	if request.Status != RequestPending {
		t.Errorf("request status = %v, want %v", request.Status, RequestPending)
	}
}

func TestBuildingRequestFloorRejectsInvalidDestination(t *testing.T) {
	building := NewBuilding(10, 5)

	for _, destination := range []int{0, -1, 11} {
		if err := building.RequestFloor(destination); err == nil {
			t.Errorf("RequestFloor(%d) returned no error for invalid destination", destination)
		}
	}
	if got := len(building.Elevator().scheduler.requests); got != 0 {
		t.Errorf("request count = %d, want 0 after invalid destinations", got)
	}
}
