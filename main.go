package main

import (
	"MyRedis/config"
	"MyRedis/server"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", config.Host, "Host to listen on")
	flag.IntVar(&config.Port, "port", config.Port, "Port to listen on")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Starting MyRedis server on", config.Host, ":", config.Port)

	var sigs chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	go server.RunAsyncTCPServer(&wg)
	go server.WaitForSignal(&wg, sigs)

	wg.Wait()

}
