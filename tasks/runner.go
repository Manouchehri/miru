package tasks

import (
	"encoding/json"
	"fmt"
	"io"

	"../models"

	"errors"
	"os"
	"os/exec"
)

// RunMonitorScript executes a monitor script in a subprocess and writes either
// a successful result or an error to a provided channel.
// Arguments:
// monitor: Information about the monitor script to run.
// result: A channel to write a result through. Can have a buffer of size 1.
// err: A channel to write an error through. Can have a buffer of size 1.
func RunMonitorScript(
	monitor models.Monitor, result chan<- models.Report, err chan<- error) {
	// Determine which interpreter to run the script with.
	// This is a bit verbose, but it allows us to control input to prevent
	// someone trying to supply a command that we don't actually want to run.
	cmdName := "python"
	switch monitor.Interpreter() {
	case models.PythonInterpreter:
		cmdName = "python"
	case models.RubyInterpreter:
		cmdName = "ruby"
	case models.PerlInterpreter:
		cmdName = "perl"
	default:
		err <- errors.New("unknown interpreter type")
		return
	}
	fmt.Println("Determined need to run interpreter", cmdName)
	// We have to run the command in another goroutine because the input part
	// of the pipe has to have started reading by the time the output part
	// starts writing, or else we get a deadlock.
	pipeIn, pipeOut := io.Pipe()
	go func() {
		defer pipeOut.Close()
		cmd := exec.Command(cmdName, monitor.ScriptPath())
		cmd.Stdout = pipeOut
		cmd.Stderr = os.Stderr
		fmt.Println("Running script")
		startErr := cmd.Run()
		if startErr != nil {
			err <- startErr
		}
	}()
	// Decode the input into a models.Report struct or else produce an error.
	decoder := json.NewDecoder(pipeIn)
	data := models.Report{}
	fmt.Println("decoding result")
	decodeErr := decoder.Decode(&data)
	if decodeErr != nil {
		err <- decodeErr
	} else {
		result <- data
	}
}
