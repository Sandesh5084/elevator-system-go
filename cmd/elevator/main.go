package main

import (
	"fmt"
	"os"

	"elevator-system/internal/cli"
	"elevator-system/internal/elevator"
)

func main() {
	building := elevator.NewBuilding(10, 5)
	app := cli.NewCLI(building, os.Stdin, os.Stdout)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "elevator: %v\n", err)
		os.Exit(1)
	}
}
