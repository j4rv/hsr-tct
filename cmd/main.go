package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/j4rv/hsr-tct/pkg/bolt"
	"github.com/j4rv/hsr-tct/pkg/webserver"
)

func main() {
	var port int
	var open bool
	flag.IntVar(&port, "port", 8080, "Server's port")
	flag.BoolVar(&open, "open", false, "Open in browser")
	flag.Parse()

	db := bolt.New()
	disconnectDb, err := db.Init("my.db")
	if err != nil {
		log.Fatal("DB Connect: ", err)
	}
	defer disconnectDb()

	if open {
		go openWebpage(fmt.Sprintf("http://localhost:%d", port))
	}

	err = webserver.Start(port, db)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// snippet from https://github.com/icza/gowut
func openWebpage(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
