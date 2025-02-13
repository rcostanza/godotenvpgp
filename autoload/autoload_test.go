package autoload

import (
	"testing"

	"github.com/rcostanza/godotenvpgp/internal/tests"
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
