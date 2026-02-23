package csvparser

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func CsvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "csv",
		Short: "This command works with CSV files",
	}

	splitCmd := &cobra.Command{
		Use:   "split <csv-file>",
		Short: "Split a CSV file into non-overlapping parts",
		Long: `Split a CSV file into multiple non-overlapping parts.
				Example:
  					bol csv split customers.csv --parts=1,10,rest --out-dir=myDir/`,
		Args: cobra.ExactArgs(1),
	}

	splitCmd.Flags().String("parts", "", "Comma-separated parts like 1,10,rest (required)")
	if err := splitCmd.MarkFlagRequired("parts"); err != nil {
		slog.Error("failed to mark --parts as required", "err", err)
		return nil
	}

	splitCmd.Flags().String("out-dir", "./", "Output directory for generated part files (optional)")

	splitCmd.RunE = func(cmd *cobra.Command, args []string) error {
		csvFile := args[0]

		partsFlag, _ := cmd.Flags().GetString("parts")
		outDir, _ := cmd.Flags().GetString("out-dir")

		partSpecs := []string{}

		for _, token := range strings.Split(partsFlag, ",") {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			if token == "rest" {
				partSpecs = append(partSpecs, "rest")
				continue
			}
			if _, err := strconv.Atoi(token); err != nil {
				return fmt.Errorf("invalid part: %s (must be an integer or 'rest')", token)
			}
			partSpecs = append(partSpecs, token)
		}

		if len(partSpecs) == 0 {
			return fmt.Errorf("--parts must contain at least one valid value")
		}

		if err := SplitCSV(csvFile, partSpecs, outDir); err != nil {
			slog.Error("failed to split CSV", slog.Any("err", err))
			return err
		}
		return nil
	}

	cmd.AddCommand(splitCmd)
	return cmd
}
