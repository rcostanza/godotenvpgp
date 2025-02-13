package core

import (
	"os"
	"testing"

	"github.com/rcostanza/godotenvpgp/internal/tests"
)

func mockReadFile(returnString string, err error) {
	readFile = func(string) ([]byte, error) {
		tests.MockCallStack.Push("ReadFile")
		return []byte(returnString), err
	}
}

func mockedEncryptedFileContent() string {
	return `-----BEGIN PGP MESSAGE-----

wy4ECQMIfwI6AcBue9TgQSOl38b6IQD3UcKli7yp1oWdTfyLilj66dmbRBQYsq9W
0jQB5YvYkL6sVl8HxzOZ28C33BnlvM3oKB2Tbwc/UAqSrz0amZIvUjjXWM7cEZxm
59vQVHn2
=Flvc
-----END PGP MESSAGE-----`
}

// Mock methods for mockedDirEntry to fulfill the interface fs.DirEntry
type mockedDirEntry struct{ name string }

func (m mockedDirEntry) Name() string               { return m.name }
func (m mockedDirEntry) IsDir() bool                { return false }
func (m mockedDirEntry) Type() os.FileMode          { return 0 }
func (m mockedDirEntry) Info() (os.FileInfo, error) { return nil, nil }

func mockReadDir(returnStrings []string, err error) {
	readDir = func(string) ([]os.DirEntry, error) {
		tests.MockCallStack.Push("ReadDir")
		list := make([]os.DirEntry, 0)
		for _, s := range returnStrings {
			list = append(list, mockedDirEntry{name: s})
		}
		return list, err
	}
}

func runTests(t *testing.T, testList []tests.TestCase) {
	for _, test := range testList {
		t.Run(test.Name, func(t *testing.T) {
			tests.InitMockCallStack(t)
			tests.SetupEnvironment(t)
			mockReadFile("123", nil)
			test.Expect(t)
		})
	}
}
