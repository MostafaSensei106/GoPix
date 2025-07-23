package main

import (
	"testing"
	"github.com/mostafasensei106/gopix/cmd"
)

// Mock the cmd package to check if Execute is called
type mockCmd struct{}

func (m *mockCmd) Execute() {
	calledExecute = true
}

var calledExecute bool

func TestMain(t *testing.T) {
	// Replace the original cmd with the mock
	originalCmd := cmd.Execute
	cmd.Execute = new(mockCmd).Execute

	// Call the main function
	main()

	// Assert that Execute was called
	if !calledExecute {
		t.Errorf("Expected cmd.Execute() to be called, but it wasn't")
	}

	// Restore the original cmd.Execute function
	cmd.Execute = originalCmd
}
