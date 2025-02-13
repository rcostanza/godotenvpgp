package core

import (
	"os"
	"strings"
	"testing"

	"github.com/rcostanza/godotenvpgp/internal/tests"
)

func TestEncryptFile(t *testing.T) {
	runTests(t, []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				_, err := EncryptFile(".env.unencrypted")
				if err != nil {
					t.Error("Unexpected output:", err)
				}
			},
		},
		{
			Name: "Error: no password found",
			Expect: func(t *testing.T) {
				t.Setenv("ENVFILE_PASSWORD", "")
				_, err := EncryptFile(".env.unencrypted")
				if !strings.Contains(err.Error(), "password not found") {
					t.Error("Expected error, got\n", err)
				}
			},
		},
		{
			Name: "Error: cannot read file",
			Expect: func(t *testing.T) {
				mockReadFile("", os.ErrNotExist)
				_, err := EncryptFile(".env.unencrypted")
				if !strings.Contains(err.Error(), "Cannot load") {
					t.Error("Expected error, got\n", err)
				}
			},
		},
		{
			Name: "Error: cannot encrypt file",
			Expect: func(t *testing.T) {
				mockReadFile("", nil)
				_, err := EncryptFile(".env.unencrypted")
				if err != nil {
					t.Error("Expected error, got nil", err)
				}
			},
		},
	})
}

func TestDecryptFile(t *testing.T) {
	runTests(t, []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				mockReadFile(mockedEncryptedFileContent(), nil)
				_, err := DecryptFile(".env.encrypted")
				if err != nil {
					t.Error("Unexpected output:", err)
				}
			},
		},
		{
			Name: "Error: no password found",
			Expect: func(t *testing.T) {
				t.Setenv("ENVFILE_PASSWORD", "")
				_, err := DecryptFile(".env.encrypted")
				if !strings.Contains(err.Error(), "password not found") {
					t.Error("Expected error, got\n", err)
				}
			},
		},
		{
			Name: "Error: cannot read file",
			Expect: func(t *testing.T) {
				mockReadFile("", os.ErrNotExist)
				_, err := DecryptFile(".env.encrypted")
				if !strings.Contains(err.Error(), "Cannot load") {
					t.Error("Expected error, got\n", err)
				}
			},
		},
		{
			Name: "Error: cannot decrypt file",
			Expect: func(t *testing.T) {
				t.Setenv("ENVFILE_PASSWORD", "WRONG_PASSWORD")
				_, err := DecryptFile(".env.encrypted")
				if !strings.Contains(err.Error(), "Cannot decrypt") {
					t.Error("Expected error, got\n", err)
				}
			},
		},
	})
}
