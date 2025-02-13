package cli

import (
	"os"
	"testing"

	"github.com/rcostanza/godotenvpgp/internal/tests"
)

func mockWriteFile(err error) {
	writeFile = func(string, []byte, os.FileMode) error { tests.MockCallStack.Push("WriteFile"); return err }
}

func mockEncryptFile(returnString string, err error) {
	encryptFile = func(string) ([]byte, error) {
		tests.MockCallStack.Push("EncryptFile")
		return []byte(returnString), err
	}
}

func mockDecryptFile(returnString string, err error) {
	decryptFile = func(string) ([]byte, error) {
		tests.MockCallStack.Push("DecryptFile")
		return []byte(returnString), err
	}
}

func mockFindEnvFiles(fileList []string) {
	findEnvFiles = func(string) []string { tests.MockCallStack.Push("FindEnvFiles"); return fileList }
}

func mockOSStat(fileInfo os.FileInfo, err error) {
	osStat = func(string) (os.FileInfo, error) { tests.MockCallStack.Push("OSStat"); return fileInfo, err }
}

func mockBail() {
	bail = func(int) {
		tests.MockCallStack.Push("Bail")
	}
}

func runTests(t *testing.T, testList []tests.TestCase) {
	for _, test := range testList {
		t.Run(test.Name, func(t *testing.T) {
			tests.InitMockCallStack(t)
			tests.SetupEnvironment(t)
			mockBail()
			mockWriteFile(nil)
			mockEncryptFile("123", nil)
			mockDecryptFile("123", nil)
			test.Expect(t)
		})
	}
}
