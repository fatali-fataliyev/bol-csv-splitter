package csvparser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createSampleCSV(t *testing.T, rowCount int) string {
	t.Helper()

	file, err := os.CreateTemp("", "sample.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer file.Close()

	file.WriteString("id,name\n")
	for i := 1; i <= rowCount; i++ {
		file.WriteString(fmt.Sprintf("%d,user%d\n", i, i))
	}

	return file.Name()
}

func TestSplitCSV_Normal(t *testing.T) {
	csvPath := createSampleCSV(t, 10) // 10 rows
	outDir := t.TempDir()

	parts := []string{"1", "2", "rest"}
	err := SplitCSV(csvPath, parts, outDir)
	if err != nil {
		t.Fatalf("test failed in normal mode: %v", err)
	}

	expectedFiles := []string{
		"sample_part1_1row.csv",
		"sample_part2_2rows.csv",
		"sample_part3_rest_7_rows.csv",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(outDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file not found: %s", file)
		}
	}
}

func TestSplitCSV_NotEnoughRows(t *testing.T) {
	csvPath := createSampleCSV(t, 3)
	outDir := t.TempDir()

	parts := []string{"2", "5"}
	err := SplitCSV(csvPath, parts, outDir)
	if err != nil {
		t.Fatalf("not enough test error: %v", err)
	}

	expectedFiles := []string{
		"sample_part1_2rows.csv",
		"sample_part2_1row.csv",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(outDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file not found: %s", file)
		}
	}
}

func TestSplitCSV_OnlyRest(t *testing.T) {
	csvPath := createSampleCSV(t, 5)
	outDir := t.TempDir()

	parts := []string{"rest"}
	err := SplitCSV(csvPath, parts, outDir)
	if err != nil {
		t.Fatalf("only rest test error: %v", err)
	}

	expected := "sample_part1_rest_5_rows.csv"
	path := filepath.Join(outDir, expected)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file not found: %s", expected)
	}
}
