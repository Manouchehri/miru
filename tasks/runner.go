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
// lastReport: The last report that the monitor generated.
// result: A channel to write a result through. Can have a buffer of size 1.
// err: A channel to write an error through. Can have a buffer of size 1.
func RunMonitorScript(
	monitor models.Monitor,
	lastReport models.Report,
	result chan<- models.Report,
	err chan<- error) {
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
	defer pipeIn.Close()
	defer pipeOut.Close()
	go func() {
		cmd := exec.Command(cmdName, monitor.ScriptPath())
		cmd.Stdout = pipeOut
		cmd.Stderr = os.Stderr
		fmt.Println("Running script with lastReport", lastReport)
		stdin, getInputErr := cmd.StdinPipe()
		// TODO - This isn't actually taking care of stdin problems.
		//        We should be timing out scripts that we can't write to soon.
		if getInputErr != nil {
			fmt.Println("Couldn't get stdin for script", getInputErr)
			err <- getInputErr
		}
		go func() {
			defer stdin.Close()
			encoded := lastReport.String()
			fmt.Println("encoded lastReport to", encoded)
			_, err := io.WriteString(stdin, encoded)
			fmt.Println("Error writing to script:", err)
		}()
		startErr := cmd.Run()
		if startErr != nil {
			fmt.Println("start error", startErr)
			err <- startErr
		}
		fmt.Println("Finished running")
	}()
	// Decode the input into a models.Report struct or else produce an error.
	data := make(map[string]interface{})
	decoder := json.NewDecoder(pipeIn)
	decodeErr := decoder.Decode(&data)
	if decodeErr != nil {
		fmt.Println("Failed to decode", decodeErr)
		err <- decodeErr
	} else {
		fmt.Println("Successfully decoded data", data)
		changeSig, found1 := data["lastChangeSignificance"]
		message, found2 := data["message"]
		checksum, found3 := data["checksum"]
		newState, found4 := data["state"]
		if !found1 || !found2 || !found3 || !found4 {
			fmt.Println("Didn't find all expected fields")
			fmt.Println(found1, found2, found3, found4)
			err <- errors.New("script output invalid data")
			return
		}
		switch changeSig.(type) {
		case uint:
			lastReport.SetChange(models.Importance(changeSig.(uint)))
		case float64:
			lastReport.SetChange(models.Importance(uint(changeSig.(float64))))
		}
		lastReport.SetMessage(message.(string))
		lastReport.SetChecksum(checksum.(string))
		lastReport.SetState(newState.(map[string]interface{}))
		fmt.Println("Put together report with message", lastReport.Message())
		result <- lastReport
	}
}
