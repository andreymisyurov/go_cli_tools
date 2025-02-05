package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMyRotate(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_rotate")
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.log")
	expectedContent := "Log entry 1\nLog entry 2\n"
	os.WriteFile(testFile, []byte(expectedContent), 0644)

	myCmd := exec.Command("./rotate", "-a", tmpDir, testFile)
	myCmd.Run()

	files, _ := filepath.Glob(filepath.Join(tmpDir, "test_*tar.gz"))
	if len(files) == 0 {
		t.Fatal("No archive created")
	}

	extractDir := filepath.Join(tmpDir, "extract")
	os.Mkdir(extractDir, 0755)
	tarCmd := exec.Command("tar", "-xzf", files[0], "-C", extractDir)
	tarCmd.Run()

	extractedFile := filepath.Join(extractDir, "test.log")
	content, _ := os.ReadFile(extractedFile)

	if string(content) != expectedContent {
		t.Fatalf("Extracted content mismatch:\nExpected: %s\nGot: %s", expectedContent, string(content))
	}
}
