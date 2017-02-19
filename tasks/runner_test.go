package tasks

import (
	"../models"

	"os"
	"testing"
)

const testPythonScript = `
import json, sys
d = json.dumps({"message": "hello world"})
print(d)
sys.exit(0)
`

const testRubyScript = `
require 'json'
d = {"message" => "hello world"}
puts d.to_json
exit 0
`

const testPerlScript = `
my $json = '{"message": "hello world"}';
print $json;
exit 0;
`

const testPythonErrorScript = `
import sys
print("hi")
sys.exit(1)
`

func TestMain(m *testing.M) {
	f1, _ := os.Create("testpython.py")
	f1.Write([]byte(testPythonScript))
	defer f1.Close()
	f2, _ := os.Create("testruby.rb")
	f2.Write([]byte(testRubyScript))
	defer f2.Close()
	f3, _ := os.Create("testperl.pl")
	f3.Write([]byte(testPerlScript))
	defer f3.Close()
	f4, _ := os.Create("testerror.py")
	f4.Write([]byte(testPythonErrorScript))
	defer f4.Close()
	exitCode := m.Run()
	os.Remove("testpython.py")
	os.Remove("testruby.rb")
	os.Remove("testperl.pl")
	os.Remove("testerror.py")
	os.Exit(exitCode)
}

func TestRunPython(t *testing.T) {
	t.Log("Running python script")
	monitor := models.NewMonitor(
		models.Administrator{}, models.PythonInterpreter, "testpython.py", 0, 0)
	resultOut := make(chan Result, 1)
	errorOut := make(chan error, 1)
	RunMonitorScript(monitor, resultOut, errorOut)
	select {
	case r := <-resultOut:
		{
			if r.Message != "hello world" {
				t.Errorf("expected to be able to parse JSON output from the script")
			}
		}
	case e := <-errorOut:
		{
			t.Errorf("expected not to get an error: %v", e)
		}
	}
}

func TestRunRuby(t *testing.T) {
	t.Log("Running Ruby script")
	monitor := models.NewMonitor(
		models.Administrator{}, models.RubyInterpreter, "testruby.rb", 0, 0)
	resultOut := make(chan Result, 1)
	errorOut := make(chan error, 1)
	RunMonitorScript(monitor, resultOut, errorOut)
	select {
	case r := <-resultOut:
		{
			if r.Message != "hello world" {
				t.Errorf("expected to be able to parse JSON output from the script")
			}
		}
	case e := <-errorOut:
		{
			t.Errorf("expected not to get an error: %v", e)
		}
	}
}

func TestRunPerl(t *testing.T) {
	monitor := models.NewMonitor(
		models.Administrator{}, models.PerlInterpreter, "testperl.pl", 0, 0)
	resultOut := make(chan Result, 1)
	errorOut := make(chan error, 1)
	RunMonitorScript(monitor, resultOut, errorOut)
	select {
	case r := <-resultOut:
		{
			if r.Message != "hello world" {
				t.Errorf("expected to be able to parse JSON output from the script")
			}
		}
	case e := <-errorOut:
		{
			t.Errorf("expected not to get an error: %v", e)
		}
	}
}

func TestRunUnknownFails(t *testing.T) {
	monitor := models.NewMonitor(
		models.Administrator{}, models.Interpreter("unknown"), "testunknown", 0, 0)
	resultOut := make(chan Result, 1)
	errorOut := make(chan error, 1)
	RunMonitorScript(monitor, resultOut, errorOut)
	select {
	case <-resultOut:
		{
			t.Errorf("expected not to get a result")
		}
	case e := <-errorOut:
		{
			t.Logf("got expected error %v", e)
		}
	}
}

func TestRunFailProducesError(t *testing.T) {
	monitor := models.NewMonitor(
		models.Administrator{}, models.PythonInterpreter, "testerror.py", 0, 0)
	resultOut := make(chan Result, 1)
	errorOut := make(chan error, 1)
	RunMonitorScript(monitor, resultOut, errorOut)
	select {
	case <-resultOut:
		{
			t.Errorf("expected not to get a result")
		}
	case e := <-errorOut:
		{
			t.Logf("got expected error %v", e)
		}
	}
}
