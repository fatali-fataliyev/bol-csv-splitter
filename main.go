package main

import (
	"fmt"
	"log/slog"

	"github.com/fatali-fataliyev/bol-csv-splitter/commands/csvparser"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "bol",
		Short: "CLI tool to split CSV file.",
	}

	rootCmd.AddCommand(csvparser.CsvCmd())

	if err := rootCmd.Execute(); err != nil {
		terminateErr := fmt.Sprintf("failed to execute root command: %s", err)
		slog.Error(terminateErr)
		return
	}

}
