package elevator

type DoorState int

const (
	DoorClosed DoorState = iota
	DoorOpen
)

type Door struct {
	state DoorState
}

func NewDoor() *Door {
	return &Door{state: DoorClosed}
}

func (d *Door) Open() {
	d.state = DoorOpen
}

func (d *Door) Close() {
	d.state = DoorClosed
}

func (d *Door) State() DoorState {
	return d.state
}
