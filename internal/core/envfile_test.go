package core

import (
	"godotenvpgp/internal/tests"
	"os"
	"testing"
)

func TestLoadUnencrypted(t *testing.T) {
	tests := []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				mockReadFile("key=value", nil)
				err := LoadUnencrypted()
				if err != nil {
					t.Error("Did not expect error:", err)
				}
				if os.Getenv("key") != "value" {
					t.Error("Expected value to be set")
				}
			},
		},
		{
			Name: "Error: cannot read file",
			Expect: func(t *testing.T) {
				mockReadFile("", os.ErrInvalid)
				err := LoadUnencrypted()
				if err == nil {
					t.Error("Expected error")
				}
			},
		},
	}
	runTests(t, tests)
}

func TestLoadEncrypted(t *testing.T) {
	tests := []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				mockReadFile(mockedEncryptedFileContent(), nil)
				err := LoadEncrypted(".env.encrypted")
				if err != nil {
					t.Error("Did not expect error:", err)
				}
				if os.Getenv("key") != "value" {
					t.Error("Expected value to be set")
				}
			},
		},
		{
			Name: "Error: cannot decrypt file",
			Expect: func(t *testing.T) {
				mockReadFile("", os.ErrInvalid)
				err := LoadEncrypted(".env.encrypted")
				if err == nil {
					t.Error("Expected error")
				}
			},
		},
	}
	runTests(t, tests)
}

func TestFindEnvFiles(t *testing.T) {
	tests := []tests.TestCase{
		{
			Name: "Success: encrypted",
			Expect: func(t *testing.T) {
				mockReadDir([]string{".env.encrypted", ".env.unencrypted"}, nil)
				files := FindEnvFiles("encrypted")
				if len(files) != 1 {
					t.Error("Expected 1 file")
				}
			},
		},
		{
			Name: "Success: unencrypted",
			Expect: func(t *testing.T) {
				mockReadDir([]string{".env.encrypted", ".env.unencrypted"}, nil)
				files := FindEnvFiles("unencrypted")
				if len(files) != 1 {
					t.Error("Expected 1 file")
				}
			},
		},
		{
			Name: "No files",
			Expect: func(t *testing.T) {
				mockReadDir([]string{}, nil)
				files := FindEnvFiles("encrypted")
				if len(files) != 0 {
					t.Error("Expected 0 files")
				}
			},
		},
		{
			Name: "Sorting order: default first",
			Expect: func(t *testing.T) {
				mockReadDir([]string{".env.encrypted", ".env.dev.encrypted"}, nil)
				files := FindEnvFiles("encrypted")
				if len(files) != 2 {
					t.Error("Expected 2 files")
				} else if files[1] != ".env.dev.encrypted" {
					t.Error("Expected .env.dev.encrypted to be last")
				}
			},
		},
		{
			Name: "Error: unknown file type argument",
			Expect: func(t *testing.T) {
				mockReadDir([]string{".env.encrypted"}, nil)
				files := FindEnvFiles("unknown")
				if len(files) != 0 {
					t.Error("Expected 0 files")
				}
			},
		},
		{
			Name: "Error: cannot read directory",
			Expect: func(t *testing.T) {
				mockReadDir([]string{}, os.ErrInvalid)
				output := tests.CaptureOutput(func() {
					FindEnvFiles("encrypted")
				})
				if output == "" {
					t.Error("Expected error")
				}
			},
		},
	}
	runTests(t, tests)
}
