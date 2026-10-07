package elevator

import "testing"

func TestSchedulerReturnsPendingRequestsInArrivalOrder(t *testing.T) {
	scheduler := NewScheduler()
	first := NewRequest(3, DirectionUp, ExternalRequest)
	second := NewRequest(7, DirectionIdle, InternalRequest)
	scheduler.AddRequest(first)
	scheduler.AddRequest(second)

	if got := scheduler.NextRequest(NewElevator(8)); got != first {
		t.Errorf("NextRequest() = %p, want first request %p", got, first)
	}

	first.Status = RequestCompleted
	if got := scheduler.NextRequest(NewElevator(8)); got != second {
		t.Errorf("NextRequest() after completing first = %p, want second request %p", got, second)
	}
}

func TestSchedulerDoesNotReturnCompletedRequests(t *testing.T) {
	scheduler := NewScheduler()
	completed := NewRequest(3, DirectionUp, ExternalRequest)
	completed.Status = RequestCompleted
	scheduler.AddRequest(completed)

	if got := scheduler.NextRequest(NewElevator(8)); got != nil {
		t.Errorf("NextRequest() = %p, want nil when no pending requests", got)
	}
}

func TestSchedulerReturnsNilWhenEmpty(t *testing.T) {
	scheduler := NewScheduler()

	if got := scheduler.NextRequest(NewElevator(8)); got != nil {
		t.Errorf("NextRequest() = %p, want nil", got)
	}
}
