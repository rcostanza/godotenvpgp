package envfile

import (
	"os"
	"strings"
	"testing"

	"github.com/rcostanza/godotenvpgp/internal/tests"
)

func mockLoadUnencrypted(err error) {
	loadUnencrypted = func() error { return err }
}

func mockLoadEncrypted(err error) {
	loadEncrypted = func(string) error { return err }
}

func mockFindEnvFiles(fileList []string) {
	findEnvFiles = func(string) []string { return fileList }
}

func TestLoad(t *testing.T) {

	var err error
	for _, test := range []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				mockFindEnvFiles([]string{".env.encrypted"})
				err = Load()
				if err != nil {
					t.Error("Did not expect error:", err)
				}
			},
		},
		{
			Name: "Warning: no unencrypted env file",
			Expect: func(t *testing.T) {
				mockLoadUnencrypted(os.ErrNotExist)
				output := tests.CaptureOutput(func() { Load() })
				if !strings.Contains(output, "Cannot parse unencrypted env file") {
					t.Error("Expected error, got", output)
				}
			},
		},
		{
			Name: "Warning: no encrypted env files",
			Expect: func(t *testing.T) {
				mockFindEnvFiles([]string{})
				output := tests.CaptureOutput(func() { Load() })
				if !strings.Contains(output, "No encrypted env files found") {
					t.Error("Expected error, got", output)
				}
			},
		},
		{
			Name: "Error: cannot load encrypted env file",
			Expect: func(t *testing.T) {
				mockFindEnvFiles([]string{".env.encrypted"})
				mockLoadEncrypted(os.ErrNotExist)
				err = Load()
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), "Cannot load encrypted env file") {
					t.Error("Expected error, got", err)
				}
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			tests.InitMockCallStack(t)
			tests.SetupEnvironment(t)
			mockLoadUnencrypted(nil)
			mockLoadEncrypted(nil)
			test.Expect(t)
		})
	}

}
