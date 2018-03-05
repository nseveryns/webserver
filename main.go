package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/nseveryns/webserver/configuration"
	"github.com/nseveryns/webserver/provider"
)

var (
	config configuration.Configuration
)

func main() {
	loadConfiguration()
	startNetwork()
}

func loadConfiguration() {
	path := flag.String("c", "conf.json", "Path to the configuration file")
	flag.Parse()
	config = configuration.LoadConfiguration(*path)
}

func startNetwork() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Port))
	if err != nil {
		log.Fatal("Unable to bind to port.", err)
	}
	log.Println("Now listening on ", config.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Failed to accept new connection.", err)
		}
		wrapper := provider.Create(conn, config)
		go wrapper.HandleConnection()
		defer wrapper.WrapUp()
	}
}
