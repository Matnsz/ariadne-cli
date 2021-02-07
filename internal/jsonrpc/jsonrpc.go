package jsonrpc

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
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