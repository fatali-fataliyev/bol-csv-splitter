package output

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func SavePart(outputPath string, header []string, rows [][]string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if err := writer.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write rows: %w", err)
	}

	writer.Flush()
	slog.Info("wrote part", "file", outputPath)
	return nil
}
