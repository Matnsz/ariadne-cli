package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	//    "errors"
	"github.com/ariadne-tools/ariadne-cli/internal/jsonrpc"

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
		for _, row := range format(readableRows(rows)) {
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

func prettySize(size int) string {

	expressions := []string{"B", "kB", "MB", "GB", "TB"}
	var ex int

	if size > 0 {
		ex = int(math.Log2(float64(size)) / 10)
	}
	return fmt.Sprintf("%4.2f", float64(size)/math.Pow(1024, float64(ex))) + " " + expressions[ex]
}

func prettyTime(nsec int) string {

	t := time.Unix(0, int64(nsec)).Format("2006-01-02 15:04:05")

	return t
}

func prettyDir(isDir bool) string {
	if isDir {
		return "dir"
	} else {
		return ""
	}
}

func readableRows(fps []jsonrpc.FileProperties) [][]string {
	rows := make([][]string, 0, len(fps))

	for _, fp := range fps {
		row := []string{fp.Path_to_file, fp.Fname, prettySize(fp.Size), prettyTime(fp.Mtime_ns), prettyDir(fp.IsDir)}
		rows = append(rows, row)
	}

	return rows
}

func format(rows [][]string) []string {

	lines := make([]string, 0, len(rows)+2)
	header := []string{"path", "name", "size", "date", "type"}

	maxLengths := make([]int, len(header))
	for _, row := range rows {
		for i, p := range row {
			if len(p) > maxLengths[i] {
				maxLengths[i] = len(p)
			}
		}
	}
	for i, v := range header {
		if len(v) > maxLengths[i] {
			maxLengths[i] = len(v)
		}
	}

	sumLength := 0
	for _, v := range maxLengths {
		sumLength += v
	}

	makeLine := func(props []string, ml []int) string {
		line := ""
		for i := 0; i != len(props); i++ {
			pad := "%" + strconv.Itoa(-ml[i]-1) + "s"
			line = line + fmt.Sprintf(pad, props[i])
		}
		return line
	}

	lines = append(lines, makeLine(header, maxLengths))
	lines = append(lines, strings.Repeat("-", sumLength+len(header)-1))

	for _, row := range rows {
		lines = append(lines, makeLine(row, maxLengths))
	}

	return lines
}
