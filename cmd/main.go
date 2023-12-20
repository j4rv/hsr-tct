package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/j4rv/hsr-tct/pkg/hsrtct"
	"github.com/j4rv/hsr-tct/pkg/mongodb"
	"github.com/j4rv/hsr-tct/pkg/webserver"
)

func main() {
	db := mongodb.New()
	disconnectDb, err := db.Connect()
	if err != nil {
		log.Fatal("DB Connect: ", err)
	}
	defer disconnectDb()
	err = db.AddCharacter(hsrtct.Character{
		Name:    "test",
		Level:   1,
		Element: hsrtct.Fire,
	})
	if err != nil {
		log.Fatal("AddLightcone: ", err)
	}
}

func main2() {
	var port int
	var open bool
	flag.IntVar(&port, "port", 8080, "Server's port")
	flag.BoolVar(&open, "open", false, "Open in browser")
	flag.Parse()

	db := mongodb.New()
	disconnectDb, err := db.Connect()
	if err != nil {
		log.Fatal("DB Connect: ", err)
	}
	defer disconnectDb()
	webserver.Setup(db)

	if open {
		go openWebpage(fmt.Sprintf("http://localhost:%d", port))
	}

	err = webserver.Start(port)
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
