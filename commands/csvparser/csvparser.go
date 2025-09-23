package csvparser

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatali-fataliyev/bol-csv-splitter/output"
)

func SplitCSV(csvFile string, parts []string, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outDir, err)
	}

	inputFile, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	allRows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read all rows: %w", err)
	}

	totalRows := len(allRows)
	rowsLeft := totalRows
	currentRowIndex := 0

	for i, part := range parts {
		var chunk [][]string

		if part == "rest" {
			chunk = allRows[currentRowIndex:]
			rowsLeft = 0
		} else {
			targetRows, _ := strconv.Atoi(part)
			if targetRows > rowsLeft {
				targetRows = rowsLeft
				slog.Warn("not enough rows", "part", i+1)
			}

			if targetRows > 0 {
				chunk = allRows[currentRowIndex : currentRowIndex+targetRows]
			}
			currentRowIndex += targetRows
			rowsLeft -= targetRows
		}

		if len(chunk) == 0 {
			continue
		}

		baseName := strings.TrimSuffix(filepath.Base(csvFile), filepath.Ext(csvFile))
		var outputFileName string
		if part == "rest" {
			outputFileName = fmt.Sprintf("%s_part%d_rest_%d_rows.csv", baseName, i+1, len(chunk))
		} else if len(chunk) == 1 {
			outputFileName = fmt.Sprintf("%s_part%d_1row.csv", baseName, i+1)
		} else {
			outputFileName = fmt.Sprintf("%s_part%d_%drows.csv", baseName, i+1, len(chunk))
		}

		outputPath := filepath.Join(outDir, outputFileName)

		if err := output.SavePart(outputPath, header, chunk); err != nil {
			return err
		}
	}

	return nil
}
