package tests

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestWCComparisonWithOriginal(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "test_wc.txt")
	defer os.Remove(tmpFile.Name())

	content := "Hello, world!\nThis is a test.\n12345 67890\n"
	tmpFile.WriteString(content)
	tmpFile.Close()

	tests := []struct {
		flag string
	}{
		{"-l"},
		{"-w"},
		{"-m"},
	}

	for _, tt := range tests {
		myCmd := exec.Command("./wc", tt.flag, tmpFile.Name())
		myOutput, _ := myCmd.CombinedOutput()
		myResult := strings.Fields(string(myOutput))[0]

		origCmd := exec.Command("wc", tt.flag, tmpFile.Name())
		origOutput, _ := origCmd.CombinedOutput()
		origResult := strings.Fields(string(origOutput))[0]

		if myResult != origResult {
			t.Errorf("Mismatch for flag %s: got %s, expected %s", tt.flag, myResult, origResult)
		}
	}
}
