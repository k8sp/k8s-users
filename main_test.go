package main

import (
	"testing"
)

func TestRestartDocker(t *testing.T) {

	err := RestartDocker("bash")
	if err != nil {
		t.Errorf("restart docker %p error", "bash")
	}
}
