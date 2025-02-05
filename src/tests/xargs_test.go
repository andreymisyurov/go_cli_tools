package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestXargsComparisonWithOriginal(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_xargs")
	defer os.RemoveAll(tmpDir)

	testFile1 := filepath.Join(tmpDir, "test1.txt")
	os.WriteFile(testFile1, []byte("Hello Xargs"), 0644)

	myCmd := exec.Command("sh", "-c", `echo "`+testFile1+`" | ./xargs ls -l`)
	myOutput, _ := myCmd.CombinedOutput()
	myLines := strings.Split(strings.TrimSpace(string(myOutput)), "\n")

	origCmd := exec.Command("sh", "-c", `echo "`+testFile1+`" | xargs ls -l`)
	origOutput, _ := origCmd.CombinedOutput()
	origLines := strings.Split(strings.TrimSpace(string(origOutput)), "\n")

	checkResultXargs(t, myLines, origLines)
}

func checkResultXargs(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("Expected %v, got %v", want, got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Mismatch: got %s, expected %s", got[i], want[i])
		}
	}
}
