//go:build tests
// +build tests

package tests

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
	"testing"
)

const (
	ENVFILE_PASSWORD = "123"
)

func SetupEnvironment(t *testing.T) {
	t.Setenv("ENVFILE_PASSWORD", ENVFILE_PASSWORD)
}

type TestCase struct {
	Name   string
	Expect func(t *testing.T)
}

// Store calls
type BaseMockCallStack struct {
	t     *testing.T
	stack map[string]bool
}

func (m *BaseMockCallStack) Push(fn string) {
	m.stack[fn] = true
}

func (m *BaseMockCallStack) AssertCalled(fn string) bool {
	_, found := m.stack[fn]
	return found
}

var MockCallStack BaseMockCallStack

func InitMockCallStack(t *testing.T) {
	MockCallStack = BaseMockCallStack{t: t, stack: make(map[string]bool)}
}

// Source: https://medium.com/@hau12a1/golang-capturing-log-println-and-fmt-println-output-770209c791b4
func CaptureOutput(f func()) string {
	reader, writer, _ := os.Pipe()
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
