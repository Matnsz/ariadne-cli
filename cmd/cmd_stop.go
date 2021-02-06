package main

import (

	//    "errors"
	"github.com/Matnsz/ariadne-cli/internal/jsonrpc"
	"github.com/spf13/cobra"
)

// stopCmd represents the say command
var stopCmd = &cobra.Command{
	Use:   "stop-daemon",
	Short: "Stop the ariadne indexing daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonrpc.StopDaemon()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
