package cmd

import (
	"github.com/spf13/cobra"
)

var (
	tagCmd = &cobra.Command{
		Use:   "tag",
		Short: "Create, list and remove tags attached to files",
	}

	tagAddCmd    = &cobra.Command{}
	tagListCmd   = &cobra.Command{}
	tagRemoveCmd = &cobra.Command{}
)
