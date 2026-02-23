package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestPreprocessor(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(dir+"/f1.txt", []byte("OVER9000\nGNULINUX\nJTK0BIW9\nSHORT\n"), 0644)
	os.WriteFile(dir+"/f2.txt", []byte("OVER9000\nGNULINUX\nSIXTYOFF\nGBR9297T\n"), 0644)
	os.WriteFile(dir+"/f3.txt", []byte("OVER9000\nGBR9297T\nTOOLONGCODE12345\n"), 0644)

	outputPath := dir + "/valid_codes.txt"

	cmd := exec.Command("go", "run", ".")
	cmd.Env = append(os.Environ(),
		"COUPON_FILES="+dir+"/f1.txt,"+dir+"/f2.txt,"+dir+"/f3.txt",
		"OUTPUT_PATH="+outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("preprocessor failed: %v\noutput: %s", err, output)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	content := string(data)

	expectedCodes := []string{"OVER9000", "GNULINUX", "GBR9297T"}
	for _, code := range expectedCodes {
		if !strings.Contains(content, code) {
			t.Errorf("expected %q in output, got:\n%s", code, content)
		}
	}

	notExpectedCodes := []string{"JTK0BIW9", "SIXTYOFF", "SHORT", "TOOLONGCODE12345"}
	for _, code := range notExpectedCodes {
		if strings.Contains(content, code) {
			t.Errorf("did not expect %q in output", code)
		}
	}

	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 valid codes, got %d: %v", len(lines), lines)
	}

	for i := 1; i < len(lines); i++ {
		if lines[i-1] > lines[i] {
			t.Errorf("output not sorted: %q before %q", lines[i-1], lines[i])
		}
	}
}
