package cmd

import (
	"github.com/kaustuvbot/kwik-cmd/internal/export"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
	importFile  string
)

var exportCmd = &cobra.Command{
	Use:   "export [filename]",
	Short: "Export command history",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := "kwik-cmd-export.json"
		if len(args) > 0 {
			filename = args[0]
		}

		switch exportFormat {
		case "json":
			return export.ExportJSON(filename)
		case "csv":
			return export.ExportCSV(filename)
		default:
			return export.ExportJSON(filename)
		}
	},
}

var importCmd = &cobra.Command{
	Use:   "import <filename>",
	Short: "Import command history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return export.ImportJSON(args[0])
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Export format (json, csv)")
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
}
