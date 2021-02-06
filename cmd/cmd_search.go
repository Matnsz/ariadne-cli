package main

import (
	"fmt"
	//    "errors"
	"github.com/Matnsz/ariadne-cli/internal/jsonrpc"

	"github.com/spf13/cobra"
)

// searchCmd represents the say command
var searchCmd = &cobra.Command{
	Use:   "search [flags] filter",
	Short: "Lists indexed files what you can filter with [arg]",
	Long: `You can search in the list of the indexed directories and files.
You can add one filter string if you would like to filter down the result.
The filter only filters by name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		searchString := ""

		if len(args) == 1 {
			searchString = args[0]
		}
		rows := jsonrpc.Search(searchString)
		for _, row := range jsonrpc.Format(jsonrpc.ReadableRows(rows)) {
			fmt.Println(row)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
