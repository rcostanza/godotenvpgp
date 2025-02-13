package core

import (
	"os"
	"testing"

	"github.com/rcostanza/godotenvpgp/internal/tests"
)

func TestParseEnv(t *testing.T) {
	runTests(t, []tests.TestCase{
		{
			Name: "Success",
			Expect: func(t *testing.T) {
				ret := ParseEnv(`key=value`)
				if len(ret) != 1 {
					t.Error("Expected 1 value, got", ret)
				}
			},
		},
		{
			Name: "Ignoring invalid keys, empty values, and comments",
			Expect: func(t *testing.T) {
				ret := ParseEnv(`valid=value # with comment
				invalid:value
				invalid==value
				invalid=
				=value
				# comment
				# invalid=value
				`)
				if len(ret) != 1 {
					t.Error("Expected 1 value, got", ret)
				}
			},
		},
		{
			Name: "Ignoring comments",
			Expect: func(t *testing.T) {
				ret := ParseEnv("# key=value")
				if len(ret) > 0 {
					t.Error("Did not expect any values, got", ret)
				}
			},
		},
		{
			Name: "Parsing environment-specific variables",
			Expect: func(t *testing.T) {
				ret := ParseEnv(`[dev]
						key=value`)

				if len(ret) != 1 {
					t.Error("Expected 1 value, got", ret)
				} else {
					if len(ret["dev"]) != 1 {
						t.Error("Expected dev env, got", ret)
					}
				}
			},
		},
	})
}

func TestGetCurrentEnvironment(t *testing.T) {
	runTests(t, []tests.TestCase{
		{
			Name: "ENVIRONMENT",
			Expect: func(t *testing.T) {
				t.Setenv("ENVIRONMENT", "dev")
				ret := GetCurrentEnvironment()
				if ret != "dev" {
					t.Error("Expected dev, got", ret)
				}
			},
		},
		{
			Name: "environment",
			Expect: func(t *testing.T) {
				t.Setenv("environment", "dev")
				ret := GetCurrentEnvironment()
				if ret != "dev" {
					t.Error("Expected dev, got", ret)
				}
			},
		},
		{
			Name: "ENV",
			Expect: func(t *testing.T) {
				t.Setenv("ENV", "dev")
				ret := GetCurrentEnvironment()
				if ret != "dev" {
					t.Error("Expected dev, got", ret)
				}
			},
		},
		{
			Name: "env",
			Expect: func(t *testing.T) {
				t.Setenv("env", "dev")
				ret := GetCurrentEnvironment()
				if ret != "dev" {
					t.Error("Expected dev, got", ret)
				}
			},
		},
		{
			Name: "None found, return default",
			Expect: func(t *testing.T) {
				ret := GetCurrentEnvironment()
				if ret != DEFAULT_ENVIRONMENT {
					t.Error("Expected", DEFAULT_ENVIRONMENT, "got", ret)
				}
			},
		},
	})
}

func TestSetEnv(t *testing.T) {
	runTests(t, []tests.TestCase{
		{
			Name: "Set default environment variables",
			Expect: func(t *testing.T) {
				envVars := map[string]map[string]string{
					DEFAULT_ENVIRONMENT: {
						"key": "value",
					},
				}
				SetEnv(envVars)
				if os.Getenv("key") != "value" {
					t.Error("Expected value, got", os.Getenv("key"))
				}
			},
		},
		{
			Name: "Set specific environment variables",
			Expect: func(t *testing.T) {
				t.Setenv("env", "dev")
				envVars := map[string]map[string]string{
					DEFAULT_ENVIRONMENT: {
						"key": "value",
					},
					"dev": {
						"key": "valueDev",
					},
				}
				SetEnv(envVars)
				if os.Getenv("key") != "valueDev" {
					t.Error("Expected value, got", os.Getenv("key"))
				}
			},
		},
		{
			Name: "Don't set environment variables for different environment",
			Expect: func(t *testing.T) {
				t.Setenv("env", "dev")
				envVars := map[string]map[string]string{
					DEFAULT_ENVIRONMENT: {
						"key": "value",
					},
					"prod": {
						"key": "valueProd",
					},
				}
				SetEnv(envVars)
				if os.Getenv("key") != "value" {
					t.Error("Didn't expected valueProd")
				}
			},
		},
	})
}
