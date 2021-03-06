package delmo_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	. "github.com/bodymindarts/delmo/delmo"
	"github.com/bodymindarts/delmo/delmo/fakes"
)

func TestTestRunner_RunTest_NoSteps(t *testing.T) {
	config := TestConfig{}
	tasks := Tasks{}
	runner := NewTestRunner(config, tasks, TaskEnvironment{})
	var b bytes.Buffer
	out := TestOutput{
		Stdout: &b,
		Stderr: &b,
	}
	runtime := new(fakes.FakeRuntime)
	runner.RunTest(runtime, out)

	if want, got := 2, runtime.CleanupCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'Cleanup()'! Want: %d, got: %d", want, got)
	}

	if want, got := 1, runtime.StartAllCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'StartAll()'! Want: %d, got: %d", want, got)
	}

	if want, got := 1, runtime.StopAllCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'StopAll()'! Want: %d, got: %d", want, got)
	}
}

func TestTestRunner_RunTest_WithSteps(t *testing.T) {
	config := TestConfig{
		Name: "test",
		Spec: SpecConfig{
			StepConfig{
				Start:   []string{"service"},
				Stop:    []string{"service"},
				Destroy: []string{"service"},
			},
			StepConfig{
				Wait:    "fake_task",
				Timeout: 60 * time.Second,
			},
			StepConfig{
				Exec:   []string{"fake_task"},
				Assert: []string{"fake_task"},
			},
		},
	}
	tasks := Tasks{
		"fake_task": TaskConfig{
			Name: "fake_task",
		},
	}
	runner := NewTestRunner(config, tasks, TaskEnvironment{})
	var b bytes.Buffer
	out := TestOutput{
		Stdout: &b,
		Stderr: &b,
	}
	runtime := new(fakes.FakeRuntime)
	runner.RunTest(runtime, out)
	t.Logf(b.String())
	outputLines := strings.Split(b.String(), "\n")
	step := 0
	if want, got := "Starting 'test' Runtime...", outputLines[step]; want != got {
		t.Errorf("Bad step execution line %d!\nWant: '%s', got: '%s'", step, want, got)
	}

	step++
	if want, got := "Executing - <Stop: [service]>", outputLines[step]; want != got {
		t.Errorf("Bad step execution line %d!\nWant: '%s', got: '%s'", step, want, got)
	}

	step++
	if want, got := "Executing - <Destroy: [service]>", outputLines[step]; want != got {
		t.Errorf("Bad step execution line %d!\nWant: '%s', got: '%s'", step, want, got)
	}

	step++
	if want, got := "Executing - <Start: [service]>", outputLines[step]; want != got {
		t.Errorf("Bad step execution line %d!\nWant: '%s', got: '%s'", step, want, got)
	}

	step++
	if want, got := "Executing - <Wait: fake_task, Timeout: 60s>", outputLines[step]; want != got {
		t.Errorf("Bad step execution line %d!\nWant: '%s', got: '%s'", step, want, got)
	}

	step++
	if want, got := "Executing - <Exec: fake_task>", outputLines[step]; want != got {
		t.Errorf("Bad step execution line %d!\nWant: '%s', got: '%s'", step, want, got)
	}

	step++
	if want, got := "Executing - <Assert: fake_task>", outputLines[step]; want != got {
		t.Errorf("Bad step execution line %d!\nWant: '%s', got: '%s'", step, want, got)
	}

	step++
	if want, got := "Stopping 'test' Runtime...", outputLines[step]; want != got {
		t.Errorf("Bad step execution line %d!\nWant: '%s', got: '%s'", step, want, got)
	}

	if want, got := 1, runtime.StartServicesCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'StartServices()'! Want: %d, got: %d", want, got)
	}
	if want, got := 1, runtime.StopServicesCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'StopServices()'! Want: %d, got: %d", want, got)
	}
	if want, got := 3, runtime.ExecuteTaskCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'ExecuteTask()'! Want: %d, got: %d", want, got)
	}
	if want, got := 1, runtime.DestroyServicesCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'DestroyServices()'! Want: %d, got: %d", want, got)
	}
}

func TestTestRunner_NoCleanupOnFailure(t *testing.T) {
	config := TestConfig{
		Name: "test",
		Spec: SpecConfig{
			StepConfig{
				Fail: []string{"fake_task"},
			},
		},
	}
	tasks := Tasks{
		"fake_task": TaskConfig{
			Name: "fake_task",
		},
	}
	runner := NewTestRunner(config, tasks, TaskEnvironment{})
	var b bytes.Buffer
	out := TestOutput{
		Stdout: &b,
		Stderr: &b,
	}
	runtime := new(fakes.FakeRuntime)
	runner.RunTest(runtime, out)

	if want, got := 1, runtime.StartAllCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'StartAll()'! Want: %d, got: %d", want, got)
	}
	if want, got := 1, runtime.StopAllCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'StopAll()'! Want: %d, got: %d", want, got)
	}

	// Cleanup should only be called once at beginning
	if want, got := 1, runtime.CleanupCallCount(); want != got {
		t.Errorf("Wrong number of calls to 'Cleanup()'! Want: %d, got: %d", want, got)
	}
}
