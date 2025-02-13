package main

import (
	"testing"

	"github.com/rcostanza/godotenvpgp/internal/tests"
)

func TestMainfn(t *testing.T) {
	t.Run("Main", func(t *testing.T) {
		tests.InitMockCallStack(t)
		fnCli = func() { tests.MockCallStack.Push("Cli") }
		main()
		if !tests.MockCallStack.AssertCalled("Cli") {
			t.Error("Expected Cli to be called")
		}
	})
}
