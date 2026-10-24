package main

import (
	"os/exec"
	"testing"
)

func TestRestartSafeBootstrap(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("Node is needed for bootstrap integration checks")
	}
	output, err := exec.Command("node", "bootstrap.test.cjs").CombinedOutput()
	if err != nil {
		t.Fatalf("bootstrap checks: %v\n%s", err, output)
	}
}
