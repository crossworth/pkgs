package xtest

import (
	"os"
	"os/exec"
	"testing"
)

// RunInSubprocess runs the current test in a subprocess to isolate process-wide changes.
// The function should be called in the first lines of the tested function.
// It can be used for tests that require full isolation.
func RunInSubprocess(t *testing.T) {
	marker := "go_test_func_" + t.Name()
	if os.Getenv("RUN_SUBPROCESS") == marker {
		return // already in subprocess
	}
	cmd := exec.Command(os.Args[0], "-test.run=^"+t.Name()+"$")
	cmd.Env = append(os.Environ(), "RUN_SUBPROCESS="+marker)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run %s in subprocess: %v", t.Name(), err)
	}
}
