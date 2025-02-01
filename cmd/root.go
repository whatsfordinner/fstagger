package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/whatsfordinner/fstagger/internal/db"
)

var (
	rootCmd = &cobra.Command{
		Use:   "fstagger",
		Short: "Tag files on your filesystem and search using them",
		Long: `fstagger is a CLI for adding one-or-more arbitrary tags to files.
These tags can be used to search for files by tag.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := db.Init(
				context.Background(),
			); err != nil {
				fmt.Println(err.Error())
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			db.Close(context.Background())
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
