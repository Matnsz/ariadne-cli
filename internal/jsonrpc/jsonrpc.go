package jsonrpc

import (
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

func AddDir(toAdd []string) []string {

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

func RmDir(toRm []int) []int {
	var removed []int
	if err := client.Call("RemoteCall.Remove", toRm, &removed); err != nil {
		log.Fatal("RemoteCall.Remove -> Error:", err)
	}
	return removed
}

func Search(s string) []FileProperties {
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

func ReadableRows(fps []FileProperties) [][]string {
	rows := make([][]string, 0, len(fps))

	for _, fp := range fps {
		row := []string{fp.Path_to_file, fp.Fname, prettySize(fp.Size), prettyTime(fp.Mtime_ns), prettyDir(fp.IsDir)}
		rows = append(rows, row)
	}

	return rows
}

func Format(rows [][]string) []string {

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

func StopDaemon() {
	var signal struct{}

	if err := client.Call("RemoteCall.StopDaemon", signal, &signal); err != nil {
		log.Fatal("RemoteCall.StopDaemon -> Error:", err)
	} else {
		fmt.Println("daemon stopped successfully")
	}
}

func WatchedDirs() []WatchedDirsState {
	var rows []WatchedDirsState
	if err := client.Call("RemoteCall.WatchedDirs", struct{}{}, &rows); err != nil {
		log.Fatal("RemoteCall.WatchedDirs -> Error:", err)
	}
	return rows
}
