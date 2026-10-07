package elevator

type RequestType int

const (
	ExternalRequest RequestType = iota
	InternalRequest
)

type RequestStatus int

const (
	RequestPending RequestStatus = iota
	RequestCompleted
)

type Request struct {
	Floor     int
	Direction Direction
	Type      RequestType
	Status    RequestStatus
}

func NewRequest(floor int, direction Direction, requestType RequestType) *Request {
	return &Request{
		Floor:     floor,
		Direction: direction,
		Type:      requestType,
		Status:    RequestPending,
	}
}
