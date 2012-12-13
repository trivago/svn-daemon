package main

import (
	"flag"
	"github.com/kless/goconfig/config"
	"web"
	"os"
	"log"
)

var configPath = flag.String("config", "config.cfg", "config file to use")

func main() {
	flag.Parse()
	log.Print("Loading configuration")

	c, err := config.ReadDefault(*configPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	pm := new(web.PageManager)
	pm.StartServer(c)
	log.Print("Server started")
}
