package main

import (
	"fmt"

	"github.com/ariadne-tools/ariadne-cli/internal/jsonrpc"

	"github.com/spf13/cobra"
)

// addCmd represents the say command
var addCmd = &cobra.Command{
	Use:   "add [flags] path(s)",
	Short: "Add new directory/directories to the index",
	RunE: func(cmd *cobra.Command, args []string) error {
		added := jsonrpc.AddDir(args)
		fmt.Println("The following dirs added successfully:")
		for _, a := range added {
			fmt.Println(a)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
