package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/rcostanza/godotenvpgp/internal/tests"
)

func TestCli(t *testing.T) {
	run := func() string {
		return tests.CaptureOutput(Cli)
	}
	runTests(t, []tests.TestCase{
		{
			Name: "Help/empty command",
			Expect: func(t *testing.T) {
				os.Args = []string{"godotenvpgp"}
				output := run()
				if !strings.Contains(output, "Usage:") {
					t.Error("Expected help message")
				}
			},
		},
		{
			Name: "Encrypt",
			Expect: func(t *testing.T) {
				os.Args = []string{"godotenvpgp", "encrypt"}
				mockFindEnvFiles([]string{".env.unencrypted"})
				mockEncryptFile("", nil)
				run()
				if tests.MockCallStack.AssertCalled("Bail") {
					t.Error("Didn't expect function to be called: Bail")
				}
				if !tests.MockCallStack.AssertCalled("FindEnvFiles") {
					t.Error("Expected function to be called: FindEnvFiles")
				}
				if !tests.MockCallStack.AssertCalled("EncryptFile") {
					t.Error("Expected function to be called: EncryptFile")
				}
			},
		},
		{
			Name: "Decrypt",
			Expect: func(t *testing.T) {
				os.Args = []string{"godotenvpgp", "decrypt"}
				mockFindEnvFiles([]string{".env.encrypted"})
				mockDecryptFile("", nil)
				run()
				if tests.MockCallStack.AssertCalled("Bail") {
					t.Error("Didn't expect function to be called: Bail")
				}
				if !tests.MockCallStack.AssertCalled("FindEnvFiles") {
					t.Error("Expected function to be called: FindEnvFiles")
				}
				if !tests.MockCallStack.AssertCalled("DecryptFile") {
					t.Error("Expected function to be called: DecryptFile")
				}
			},
		},
		{
			Name: "Show",
			Expect: func(t *testing.T) {
				os.Args = []string{"godotenvpgp", "show", ".env.encrypted"}
				expected := "VAR=1"
				mockOSStat(nil, nil)
				mockDecryptFile(expected, nil)
				output := run()
				if output != expected+"\n" {
					t.Errorf("Output did not match:\nGot:'%s'\nExp:'%s'", output, expected)
				}
			},
		},
		{
			Name: "Show: file not found",
			Expect: func(t *testing.T) {
				os.Args = []string{"godotenvpgp", "show", ".env.encrypted"}
				mockOSStat(nil, os.ErrNotExist)
				output := run()
				if !strings.Contains(output, "not found") {
					t.Error("Expected error")
				}
			},
		},
		{
			Name: "Show: invalid file",
			Expect: func(t *testing.T) {
				os.Args = []string{"godotenvpgp", "show", "invalid"}
				output := run()
				if !strings.Contains(output, "is not encrypted") {
					t.Error("Expected error")
				}
			},
		},
		{
			Name: "Invalid command",
			Expect: func(t *testing.T) {
				os.Args = []string{"godotenvpgp", "invalid"}
				output := run()
				if !strings.Contains(output, "Usage:") {
					t.Error("Expected help message")
				}
			},
		},
	})
}

func TestHelp(t *testing.T) {
	runTests(t, []tests.TestCase{
		{
			Name: "Help",
			Expect: func(t *testing.T) {
				output := tests.CaptureOutput(func() {
					showHelp()
				})
				if !strings.Contains(output, "Usage:") {
					t.Error("Expected help message")
				}
			},
		},
	})
}

func TestEncryptDecrypt(t *testing.T) {
	runTests(t, []tests.TestCase{
		{
			Name: "Encrypt",
			Expect: func(t *testing.T) {
				mockFindEnvFiles([]string{".env.unencrypted"})
				mockEncryptFile("", nil)
				tests.CaptureOutput(func() {
					err := encryptDecrypt("encrypt")
					if err != nil {
						t.Error("Did not expect error:", err)
					}
				})
			},
		},
		{
			Name: "Decrypt",
			Expect: func(t *testing.T) {
				mockFindEnvFiles([]string{".env.encrypted"})
				mockDecryptFile("", nil)
				tests.CaptureOutput(func() {
					err := encryptDecrypt("decrypt")
					if err != nil {
						t.Error("Did not expect error:", err)
					}
				})
			},
		},
		{
			Name: "No files found",
			Expect: func(t *testing.T) {
				mockFindEnvFiles([]string{})
				err := encryptDecrypt("encrypt")
				if err == nil {
					t.Error("Expected error")
				}
			},
		},
		{
			Name: "Failed encryption",
			Expect: func(t *testing.T) {
				mockFindEnvFiles([]string{".env.unencrypted"})
				mockEncryptFile("", os.ErrInvalid)
				output := tests.CaptureOutput(func() {
					encryptDecrypt("encrypt")
				})
				if strings.Contains(".env.unencrypted failed", output) {
					t.Error("Expected error")
				}
			},
		},
		{
			Name: "Failed decryption",
			Expect: func(t *testing.T) {
				mockFindEnvFiles([]string{".env.encrypted"})
				mockDecryptFile("", os.ErrInvalid)
				output := tests.CaptureOutput(func() {
					encryptDecrypt("decrypt")
				})
				if strings.Contains(".env.encrypted failed", output) {
					t.Error("Expected error")
				}
			},
		},
	})
}

func TestShow(t *testing.T) {
	runTests(t, []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				expected := "VAR=1"
				mockDecryptFile(expected, nil)
				mockOSStat(nil, nil)
				output := tests.CaptureOutput(func() {
					err := listVars(".env.encrypted")
					if err != nil {
						t.Error("Did not expect error:", err)
					}
				})
				if output != expected+"\n" {
					t.Errorf("Output did not match:\nGot:'%s'\nExp:'%s'", output, expected)
				}
			},
		},
		{
			Name: "Not enough arguments",
			Expect: func(t *testing.T) {
				os.Args = []string{"godotenvpgp", "show"}
				output := tests.CaptureOutput(Cli)
				if !strings.Contains(output, "Usage:") {
					t.Error("Expected help message")
				}
			},
		},
		{
			Name: "File not encrypted",
			Expect: func(t *testing.T) {
				os.Args = []string{"godotenvpgp", "show", ".env.notencrypted"}
				output := tests.CaptureOutput(Cli)
				if !strings.Contains(output, "not encrypted") {
					t.Error("Expected error")
				}
			},
		},
		{
			Name: "File does not exist",
			Expect: func(t *testing.T) {
				mockOSStat(nil, os.ErrNotExist)
				os.Args = []string{"godotenvpgp", "show", ".env.encrypted"}
				output := tests.CaptureOutput(Cli)
				if !strings.Contains(output, "not found") {
					t.Error("Expected error")
				}
			},
		},
		{
			Name: "Error: cannot decrypt file",
			Expect: func(t *testing.T) {
				mockOSStat(nil, nil)
				mockDecryptFile("", os.ErrInvalid)
				os.Args = []string{"godotenvpgp", "show", ".env.encrypted"}
				output := tests.CaptureOutput(Cli)
				if !strings.Contains(output, "cannot list vars") {
					t.Error("Expected error")
				}
			},
		},
	})
}
