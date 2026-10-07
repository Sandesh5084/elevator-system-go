package elevator

import "testing"

func TestNewExternalRequest(t *testing.T) {
	request := NewRequest(3, DirectionUp, ExternalRequest)

	if request.Floor != 3 {
		t.Errorf("Floor = %d, want 3", request.Floor)
	}
	if request.Direction != DirectionUp {
		t.Errorf("Direction = %v, want %v", request.Direction, DirectionUp)
	}
	if request.Type != ExternalRequest {
		t.Errorf("Type = %v, want %v", request.Type, ExternalRequest)
	}
	if request.Status != RequestPending {
		t.Errorf("Status = %v, want %v", request.Status, RequestPending)
	}
}

func TestNewInternalRequest(t *testing.T) {
	request := NewRequest(7, DirectionIdle, InternalRequest)

	if request.Floor != 7 {
		t.Errorf("Floor = %d, want 7", request.Floor)
	}
	if request.Type != InternalRequest {
		t.Errorf("Type = %v, want %v", request.Type, InternalRequest)
	}
	if request.Status != RequestPending {
		t.Errorf("Status = %v, want %v", request.Status, RequestPending)
	}
}
