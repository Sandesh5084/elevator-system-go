package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"elevator-system/internal/elevator"
)

type CLI struct {
	building *elevator.Building
	input    io.Reader
	output   io.Writer
	scanner  *bufio.Scanner
}

func NewCLI(building *elevator.Building, input io.Reader, output io.Writer) *CLI {
	return &CLI{
		building: building,
		input:    input,
		output:   output,
		scanner:  bufio.NewScanner(input),
	}
}

func (c *CLI) Run() error {
	if c.building == nil {
		return fmt.Errorf("CLI requires a building")
	}
	if c.input == nil {
		return fmt.Errorf("CLI requires an input reader")
	}
	if c.output == nil {
		return fmt.Errorf("CLI requires an output writer")
	}

	for {
		if _, err := fmt.Fprint(c.output,
			"\n1. External request\n"+
				"2. Internal request\n"+
				"3. Serve next request\n"+
				"4. Show elevator status\n"+
				"5. Exit\n"+
				"Enter command: ",
		); err != nil {
			return fmt.Errorf("write CLI menu: %w", err)
		}

		if !c.scanner.Scan() {
			if err := c.scanner.Err(); err != nil {
				return fmt.Errorf("read CLI command: %w", err)
			}
			return nil
		}

		switch strings.TrimSpace(c.scanner.Text()) {
		case "1":
			if err := c.ReadExternalRequest(); err != nil {
				return err
			}
		case "2":
			if err := c.ReadInternalRequest(); err != nil {
				return err
			}
		case "3":
			if err := c.ServeNextRequest(); err != nil {
				return fmt.Errorf("serve next request: %w", err)
			}
		case "4":
			if err := c.writeStatus(); err != nil {
				return err
			}
		case "5":
			return nil
		default:
			if _, err := fmt.Fprintf(c.output, "Invalid command %q. Choose 1 through 5.\n", strings.TrimSpace(c.scanner.Text())); err != nil {
				return fmt.Errorf("write invalid-command message: %w", err)
			}
		}
	}
}

func (c *CLI) ReadExternalRequest() error {
	if c.building == nil {
		return fmt.Errorf("CLI requires a building")
	}
	if c.input == nil {
		return fmt.Errorf("CLI requires an input reader")
	}
	if c.output == nil {
		return fmt.Errorf("CLI requires an output writer")
	}

	if _, err := fmt.Fprint(c.output, "Enter floor: "); err != nil {
		return fmt.Errorf("write floor prompt: %w", err)
	}
	floorInput, err := c.readLine("floor")
	if err != nil {
		return err
	}
	floor, err := strconv.Atoi(strings.TrimSpace(floorInput))
	if err != nil {
		return fmt.Errorf("parse floor %q: %w", strings.TrimSpace(floorInput), err)
	}

	if _, err := fmt.Fprint(c.output, "Enter direction (up/down): "); err != nil {
		return fmt.Errorf("write direction prompt: %w", err)
	}
	directionInput, err := c.readLine("direction")
	if err != nil {
		return err
	}

	direction, err := parseDirection(directionInput)
	if err != nil {
		return err
	}

	if err := c.building.RequestElevator(floor, direction); err != nil {
		return fmt.Errorf("submit external request: %w", err)
	}
	return nil
}

func (c *CLI) ReadInternalRequest() error {
	if c.building == nil {
		return fmt.Errorf("CLI requires a building")
	}
	if c.input == nil {
		return fmt.Errorf("CLI requires an input reader")
	}
	if c.output == nil {
		return fmt.Errorf("CLI requires an output writer")
	}

	if _, err := fmt.Fprint(c.output, "Enter destination floor: "); err != nil {
		return fmt.Errorf("write destination prompt: %w", err)
	}
	destinationInput, err := c.readLine("destination floor")
	if err != nil {
		return err
	}
	destination, err := strconv.Atoi(strings.TrimSpace(destinationInput))
	if err != nil {
		return fmt.Errorf("parse destination floor %q: %w", strings.TrimSpace(destinationInput), err)
	}

	if err := c.building.RequestFloor(destination); err != nil {
		return fmt.Errorf("submit internal request: %w", err)
	}
	return nil
}

func (c *CLI) ServeNextRequest() error {
	if c.building == nil {
		return fmt.Errorf("CLI requires a building")
	}
	return c.building.Elevator().ServeNextRequest()
}

func (c *CLI) writeStatus() error {
	elevator := c.building.Elevator()
	if _, err := fmt.Fprintf(c.output, "Floor: %d | Direction: %s | State: %s\n",
		elevator.CurrentFloor(),
		directionName(elevator.Direction()),
		stateName(elevator.State()),
	); err != nil {
		return fmt.Errorf("write elevator status: %w", err)
	}
	return nil
}

func directionName(direction elevator.Direction) string {
	switch direction {
	case elevator.DirectionUp:
		return "UP"
	case elevator.DirectionDown:
		return "DOWN"
	case elevator.DirectionIdle:
		return "IDLE"
	default:
		return "UNKNOWN"
	}
}

func stateName(state elevator.ElevatorState) string {
	switch state {
	case elevator.ElevatorMoving:
		return "MOVING"
	case elevator.ElevatorIdle:
		return "IDLE"
	default:
		return "UNKNOWN"
	}
}

func parseDirection(input string) (elevator.Direction, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "up":
		return elevator.DirectionUp, nil
	case "down":
		return elevator.DirectionDown, nil
	default:
		return elevator.DirectionIdle, fmt.Errorf("invalid direction %q: expected up or down", strings.TrimSpace(input))
	}
}

func (c *CLI) readLine(field string) (string, error) {
	if c.scanner.Scan() {
		return c.scanner.Text(), nil
	}
	if err := c.scanner.Err(); err != nil {
		return "", fmt.Errorf("read %s: %w", field, err)
	}
	return "", fmt.Errorf("read %s: unexpected end of input", field)
}
