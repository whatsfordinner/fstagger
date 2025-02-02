package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "fstagger",
		Short: "Tag files on your filesystem and search using them",
		Long: `fstagger is a CLI for adding one-or-more arbitrary tags to files.
These tags can be used to search for files by tag.`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
