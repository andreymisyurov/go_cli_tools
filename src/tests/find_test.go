package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestFindComparisonWithOriginal(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_find")
	defer os.RemoveAll(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.log"), []byte("logfile"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	os.Symlink("file1.txt", filepath.Join(tmpDir, "link1"))

	myCmd := exec.Command("./find", "-f", tmpDir)
	myOutput, _ := myCmd.CombinedOutput()
	myFiles := strings.Split(strings.TrimSpace(string(myOutput)), "\n")

	origCmd := exec.Command("find", tmpDir, "-type", "f")
	origOutput, _ := origCmd.CombinedOutput()
	origFiles := strings.Split(strings.TrimSpace(string(origOutput)), "\n")

	sort.Strings(myFiles)
	sort.Strings(origFiles)

	checkResult(t, myFiles, origFiles)
}

func checkResult(t *testing.T, got, want []string) {
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
