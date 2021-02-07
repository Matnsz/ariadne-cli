package main

import (
	"fmt"

	//    "errors"
	"github.com/ariadne-tools/ariadne-cli/internal/jsonrpc"
	"github.com/spf13/cobra"
)

// watchedDirsCmd represents the say command
var watchedDirsCmd = &cobra.Command{
	Use:   "watched-dirs",
	Short: "Get the list of the currently watched directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		ws := jsonrpc.WatchedDirs()
		for _, w := range ws {
			fmt.Println(w)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(watchedDirsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchedDirsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchedDirsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
