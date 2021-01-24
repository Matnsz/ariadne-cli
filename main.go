package main

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
	
	"github.com/matnsz/ariadne-cli/ariadne-cli/rpc"
)

const (
	VERSION = "0.3.0"
)

func main() {

	if err := rpc.root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
