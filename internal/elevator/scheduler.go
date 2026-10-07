package elevator

type Scheduler struct {
	requests []*Request
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) AddRequest(request *Request) {
	s.requests = append(s.requests, request)
}

func (s *Scheduler) NextRequest(_ *Elevator) *Request {
	for _, request := range s.requests {
		if request != nil && request.Status == RequestPending {
			return request
		}
	}

	return nil
}
