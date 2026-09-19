package main

import (
	"flag"
	"log"
	"MyRedis/config"
	"MyRedis/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", config.Host, "Host to listen on");
	flag.IntVar(&config.Port, "port", config.Port, "Port to listen on");
	flag.Parse();
}


func main() {
	setupFlags();
	log.Println("Starting MyRedis server on", config.Host, ":", config.Port);
	server.RunSyncTCPServer();

}