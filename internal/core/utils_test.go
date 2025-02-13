package core

import (
	"godotenvpgp/internal/tests"
	"testing"
)

func TestGetFilePassword(t *testing.T) {
	tests := []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				password, err := getFilePassword(".env.encrypted")
				if err != nil {
					t.Error("Did not expect error:", err)
				}
				if password == "" {
					t.Error("Expected password")
				}
			},
		},
		{
			Name: "Non-default environment",
			Expect: func(t *testing.T) {
				expectedPw := "test"
				t.Setenv("ENVFILE_PASSWORD_TESTENV", expectedPw)
				password, err := getFilePassword(".env.testenv.encrypted")
				if err != nil {
					t.Error("Did not expect error:", err)
				}
				if password != expectedPw {
					t.Errorf("Expected password %s, got %s", expectedPw, password)
				}
			},
		},
		{
			Name: "Error: password not found",
			Expect: func(t *testing.T) {
				t.Setenv("ENVFILE_PASSWORD", "")
				_, err := getFilePassword(".env.unencrypted")
				if err == nil {
					t.Error("Expected error")
				}
			},
		},
		{
			Name: "Error: not a valid .env file",
			Expect: func(t *testing.T) {
				_, err := getFilePassword("file.txt")
				if err == nil {
					t.Error("Expected error")
				}
			},
		},
	}
	runTests(t, tests)
}
