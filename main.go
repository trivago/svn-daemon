package main

import (
	"flag"
	"github.com/kless/goconfig/config"
	"log"
	"os"
	"github.com/xenji/svn-daemon/ui"
)

var configPath = flag.String("config", "config.cfg", "config file to use")

func main() {
	flag.Parse()
	c, err := config.ReadDefault(*configPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	pm := new(ui.PageManager)
	pm.StartServer(c)
}
