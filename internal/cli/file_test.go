package cli

import (
	"godotenvpgp/internal/tests"
	"os"
	"testing"
)

func TestEncryptFile(t *testing.T) {
	tests := []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				err := saveEncryptedFile(".env.unencrypted")
				if err != nil {
					t.Error("Did not expect error:", err)
				}
			},
		},
		{
			Name: "Error: cannot encrypt file",
			Expect: func(t *testing.T) {
				mockEncryptFile("", os.ErrInvalid)
				err := saveEncryptedFile(".env.unencrypted")
				if err == nil {
					t.Error("Expected error")
				}
			},
		},
		{
			Name: "Error: cannot write file",
			Expect: func(t *testing.T) {
				mockWriteFile(os.ErrPermission)
				err := saveEncryptedFile(".env.unencrypted")
				if err == nil {
					t.Error("Expected error")
				}
			},
		},
	}
	runTests(t, tests)
}

func TestDecryptFile(t *testing.T) {
	tests := []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				err := saveDecryptedFile(".env.encrypted")
				if err != nil {
					t.Error("Did not expect error:", err)
				}
			},
		},
		{
			Name: "Error: cannot decrypt file",
			Expect: func(t *testing.T) {
				mockDecryptFile("", os.ErrInvalid)
				err := saveDecryptedFile(".env.encrypted")
				if err == nil {
					t.Error("Expected error")
				}
			},
		},
		{
			Name: "Error: cannot write file",
			Expect: func(t *testing.T) {
				mockWriteFile(os.ErrPermission)
				err := saveDecryptedFile(".env.encrypted")
				if err == nil {
					t.Error("Expected error")
				}
			},
		},
	}
	runTests(t, tests)
}
