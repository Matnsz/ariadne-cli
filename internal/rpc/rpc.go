package rpc

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
)

// get RPC client by dialing at `rpc.DefaultRPCPath` endpoint
var client, _ = rpc.DialHTTP("tcp", "localhost:9000")
var usage = `Usage: ariadne-cli <command> [<arg>]

<Commands>
help : to get this help message
stop-daemon : to stop the ariadne indexing daemon
watched-dirs : to get a list from the currently watched directories
add <args> : to add a new directory/directories to the index
remove <args> : to remove directory/directories with the id <arg> from the index
search <arg> : search lists files and you can filter it with 'arg'`

type Query struct {
	Base string
	Args []interface{}
}

type WatchedDirsState struct {
	Id    int
	Path  string
	State string
}

type FileProperties struct {
	Path_to_file string
	Fname        string
	Size         int
	Mtime_ns     int
	IsDir        bool
}

func isValidDir(path string) bool {
	if info, err := os.Stat(path); err != nil {
		return false
	} else {
		return info.IsDir()
	}
}

func addDir(toAdd []string) []string {

	for _, path := range toAdd {
		if !isValidDir(path) {
			log.Fatalf("ERROR: The directory provided '%s' is not valid", path)
		}
	}

	var added []string
	if err := client.Call("RemoteCall.Add", toAdd, &added); err != nil {
		log.Fatal("RemoteCall.Add -> Error:", err)
	}
	return added
}

func rmDir(toRm []int) []int {
	var removed []int
	if err := client.Call("RemoteCall.Remove", toRm, &removed); err != nil {
		log.Fatal("RemoteCall.Remove -> Error:", err)
	}
	return removed
}

func search(s string) []FileProperties {
	var rows []FileProperties
	if err := client.Call("RemoteCall.Search", s, &rows); err != nil {
		log.Fatal("RemoteCall.Search -> Error:", err)
	}
	return rows
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

func readableRows(fps []FileProperties) [][]string {
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

func stopDaemon() {
	var signal struct{}

	if err := client.Call("RemoteCall.StopDaemon", signal, &signal); err != nil {
		log.Fatal("RemoteCall.StopDaemon -> Error:", err)
	} else {
		fmt.Println("daemon stopped successfully")
	}
}

func watchedDirs() []WatchedDirsState {
	var rows []WatchedDirsState
	if err := client.Call("RemoteCall.WatchedDirs", struct{}{}, &rows); err != nil {
		log.Fatal("RemoteCall.WatchedDirs -> Error:", err)
	}
	return rows
}

func Root(args []string) error {
	if len(args) < 1 {
		return errors.New(usage)
	} else {
		switch args[0] {
		case "stop-daemon":
			stopDaemon()
			return nil
		case "help":
			fmt.Println(usage)
			return nil
		case "watched-dirs":
			ws := watchedDirs()
			for _, w := range ws {
				fmt.Println(w)
			}
			return nil
		case "search":
			searchString := ""
			if len(args) == 2 {
				searchString = args[1]
			}
			rows := search(searchString)
			for _, row := range format(readableRows(rows)) {
				fmt.Println(row)
			}
			return nil
		case "add":
			added := addDir(args[1:])
			fmt.Println("The following dirs added successfully:")
			for _, a := range added {
				fmt.Println(a)
			}
			return nil
		case "remove":
			var intArgs []int
			for _, a := range args[1:] {
				if i, err := strconv.Atoi(a); err != nil {
					log.Fatal("Error: cannot convert" + a + "to integer. Please give integer as parameter!")
				} else {
					intArgs = append(intArgs, i)
				}
			}
			removed := rmDir(intArgs)
			fmt.Println("The following dirs removed successfully:")
			for _, r := range removed {
				fmt.Println(r)
			}
			return nil
		default:
			return errors.New(usage)
		}
	}
}
