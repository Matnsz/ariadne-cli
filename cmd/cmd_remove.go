package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ariadne-tools/ariadne-cli/internal/jsonrpc"

	//    "errors"

	"github.com/spf13/cobra"
)

// removeCmd represents the say command
var removeCmd = &cobra.Command{
	Use:   "remove [flags] index/indices",
	Short: "Remove directory/directories with the id(s) <arg(s)> from the index",
	Long: `You can remove any number of watched directories from the database.
You have to use the index/indices of the directory/directories insted of it's/their names.`,
	Run: func(cmd *cobra.Command, args []string) {
		var intArgs []int
		for _, a := range args {
			if i, err := strconv.Atoi(a); err != nil {
				log.Fatal("Error: cannot convert" + a + "to integer. Please give integer as parameter!")
			} else {
				intArgs = append(intArgs, i)
			}
		}
		removed := jsonrpc.RmDir(intArgs)
		fmt.Println("The following dirs removed successfully:")
		for _, r := range removed {
			fmt.Println(r)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
