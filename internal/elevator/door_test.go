package elevator

import "testing"

func TestNewDoorStartsClosed(t *testing.T) {
	door := NewDoor()
	if got := door.State(); got != DoorClosed {
		t.Errorf("State() = %v, want %v", got, DoorClosed)
	}
}

func TestDoorOpenAndClose(t *testing.T) {
	door := NewDoor()

	door.Open()
	if got := door.State(); got != DoorOpen {
		t.Errorf("State() after Open() = %v, want %v", got, DoorOpen)
	}

	door.Close()
	if got := door.State(); got != DoorClosed {
		t.Errorf("State() after Close() = %v, want %v", got, DoorClosed)
	}
}

func TestDoorOpenWhenAlreadyOpen(t *testing.T) {
	door := NewDoor()
	door.Open()

	door.Open()
	if got := door.State(); got != DoorOpen {
		t.Errorf("State() after repeated Open() = %v, want %v", got, DoorOpen)
	}
}

func TestDoorCloseWhenAlreadyClosed(t *testing.T) {
	door := NewDoor()

	door.Close()
	if got := door.State(); got != DoorClosed {
		t.Errorf("State() after repeated Close() = %v, want %v", got, DoorClosed)
	}
}
