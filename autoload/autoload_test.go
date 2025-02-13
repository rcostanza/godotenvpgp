package autoload

import (
	"godotenvpgp/internal/tests"
	"testing"
)

func init() {
	load = func() error {
		tests.MockCallStack.Push("Load")
		return nil
	}
}

func TestAutoload(t *testing.T) {
	t.Run("Load", func(t *testing.T) {
		tests.InitMockCallStack(t)
		load()
		if !tests.MockCallStack.AssertCalled("Load") {
			t.Error("expected Load to be called")
		}
	})
}
