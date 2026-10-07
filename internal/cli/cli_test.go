package cli

import (
	"errors"
	"strings"
	"testing"

	"elevator-system/internal/elevator"
)

func TestReadExternalRequestSubmitsParsedRequest(t *testing.T) {
	tests := []struct {
		name      string
		direction string
		want      elevator.Direction
	}{
		{name: "up", direction: "up", want: elevator.DirectionUp},
		{name: "down", direction: "DOWN", want: elevator.DirectionDown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			building := elevator.NewBuilding(10, 5)
			input := strings.NewReader("5\n" + tt.direction + "\n")
			var output strings.Builder

			if got, err := parseDirection(tt.direction); err != nil || got != tt.want {
				t.Fatalf("parseDirection(%q) = %v, %v; want %v, nil", tt.direction, got, err, tt.want)
			}
			if err := NewCLI(building, input, &output).ReadExternalRequest(); err != nil {
				t.Fatalf("ReadExternalRequest() returned error: %v", err)
			}

			if !strings.Contains(output.String(), "Enter floor: ") ||
				!strings.Contains(output.String(), "Enter direction (up/down): ") {
				t.Errorf("expected prompts in output, got %q", output.String())
			}
			if err := building.Elevator().ServeNextRequest(); err != nil {
				t.Fatalf("ServeNextRequest() returned error: %v", err)
			}
			if got := building.Elevator().CurrentFloor(); got != 5 {
				t.Errorf("CurrentFloor() after serving request = %d, want 5", got)
			}
		})
	}
}

func TestReadExternalRequestReturnsInputAndValidationErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "invalid floor input", input: "abc\nup\n"},
		{name: "invalid direction", input: "5\nsideways\n"},
		{name: "invalid building floor", input: "11\nup\n"},
		{name: "missing direction", input: "5\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			building := elevator.NewBuilding(10, 5)
			err := NewCLI(building, strings.NewReader(tt.input), &strings.Builder{}).ReadExternalRequest()
			if err == nil {
				t.Fatal("ReadExternalRequest() returned no error")
			}
			if err := building.Elevator().ServeNextRequest(); err != nil {
				t.Fatalf("ServeNextRequest() returned error: %v", err)
			}
			if got := building.Elevator().CurrentFloor(); got != 1 {
				t.Errorf("CurrentFloor() = %d, want 1 after rejected request", got)
			}
		})
	}
}

type failedReader struct{}

func (failedReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestReadExternalRequestReturnsIOErrors(t *testing.T) {
	building := elevator.NewBuilding(10, 5)

	t.Run("reader", func(t *testing.T) {
		err := NewCLI(building, failedReader{}, &strings.Builder{}).ReadExternalRequest()
		if err == nil || !strings.Contains(err.Error(), "read failed") {
			t.Errorf("ReadExternalRequest() error = %v, want read failure", err)
		}
	})

	t.Run("writer", func(t *testing.T) {
		err := NewCLI(building, strings.NewReader("5\nup\n"), failedWriter{}).ReadExternalRequest()
		if err == nil || !strings.Contains(err.Error(), "write failed") {
			t.Errorf("ReadExternalRequest() error = %v, want write failure", err)
		}
	})
}

func TestReadInternalRequestSubmitsAndServesDestination(t *testing.T) {
	building := elevator.NewBuilding(10, 5)
	var output strings.Builder
	cli := NewCLI(building, strings.NewReader("7\n"), &output)

	if err := cli.ReadInternalRequest(); err != nil {
		t.Fatalf("ReadInternalRequest() returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Enter destination floor: ") {
		t.Errorf("destination prompt missing from output %q", output.String())
	}

	if err := cli.ServeNextRequest(); err != nil {
		t.Fatalf("ServeNextRequest() returned error: %v", err)
	}
	if got := building.Elevator().CurrentFloor(); got != 7 {
		t.Errorf("CurrentFloor() = %d, want 7", got)
	}
}

func TestReadInternalRequestReturnsInputAndValidationErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "invalid floor input", input: "seven\n"},
		{name: "invalid building floor", input: "11\n"},
		{name: "missing floor", input: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			building := elevator.NewBuilding(10, 5)
			cli := NewCLI(building, strings.NewReader(tt.input), &strings.Builder{})
			if err := cli.ReadInternalRequest(); err == nil {
				t.Fatal("ReadInternalRequest() returned no error")
			}
			if err := cli.ServeNextRequest(); err != nil {
				t.Fatalf("ServeNextRequest() returned error: %v", err)
			}
			if got := building.Elevator().CurrentFloor(); got != 1 {
				t.Errorf("CurrentFloor() = %d, want 1 after rejected request", got)
			}
		})
	}
}

func TestRunProcessesCommandsUntilExit(t *testing.T) {
	building := elevator.NewBuilding(10, 5)
	input := strings.NewReader("1\n5\nup\n2\n7\n3\n4\n3\n4\n5\n")
	var output strings.Builder

	if err := NewCLI(building, input, &output).Run(); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	for _, want := range []string{
		"1. External request",
		"2. Internal request",
		"3. Serve next request",
		"4. Show elevator status",
		"5. Exit",
		"Floor: 5 | Direction: IDLE | State: IDLE",
		"Floor: 7 | Direction: IDLE | State: IDLE",
	} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("Run() output does not contain %q", want)
		}
	}
	if got := building.Elevator().CurrentFloor(); got != 7 {
		t.Errorf("CurrentFloor() after command loop = %d, want 7", got)
	}
}

func TestRunReportsInvalidCommandAndContinues(t *testing.T) {
	var output strings.Builder
	cli := NewCLI(elevator.NewBuilding(10, 5), strings.NewReader("unknown\n5\n"), &output)

	if err := cli.Run(); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}
	if !strings.Contains(output.String(), `Invalid command "unknown"`) {
		t.Errorf("Run() output does not report invalid command: %q", output.String())
	}
}

func TestRunReturnsActionErrors(t *testing.T) {
	cli := NewCLI(elevator.NewBuilding(10, 5), strings.NewReader("1\n11\nup\n"), &strings.Builder{})

	if err := cli.Run(); err == nil {
		t.Fatal("Run() returned no error for invalid external request floor")
	}
}

func TestRunReturnsOutputErrors(t *testing.T) {
	cli := NewCLI(elevator.NewBuilding(10, 5), strings.NewReader("5\n"), failedWriter{})

	if err := cli.Run(); err == nil || !strings.Contains(err.Error(), "write CLI menu") {
		t.Errorf("Run() error = %v, want CLI menu write error", err)
	}
}
