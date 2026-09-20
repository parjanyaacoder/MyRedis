package server

import (
	"MyRedis/config"
	"MyRedis/core"
	"log"
	"net"
	"syscall"
	"time"
)

var concurrent_clients int = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()

func RunAsyncTCPServer() error {
	log.Println("starting a asynchronous TCP server on", config.Host, config.Port)

	max_clients := 20000

	// Create Kqueue Event Objects array to hold triggered events - holds File discriptors
	var events []syscall.Kevent_t = make([]syscall.Kevent_t, max_clients)

	// Create a socket - IPv4, non-blocking, socket stream - does not close the connection on getting a response
	// Keeps the connection open even after a response
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}

	// close the socket once function execution ends
	defer syscall.Close(serverFD)

	// Set socket to non-blocking
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	// Set address
	ip4 := net.ParseIP(config.Host).To4()
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	if err := syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}

	// Async IO starts here

	// Creating KQUEUE instance
	KqueueFD, err := syscall.Kqueue()

	if err != nil {
		return err
	}

	defer syscall.Close(KqueueFD)

	// Specify the event we want to listen for (Read event on serverFD)
	var socketServerEvent syscall.Kevent_t
	syscall.SetKevent(&socketServerEvent, serverFD, syscall.EVFILT_READ, syscall.EV_ADD)

	// Listen to read events on the Server itself
	if _, err = syscall.Kevent(KqueueFD, []syscall.Kevent_t{socketServerEvent}, nil, nil); err != nil {
		return err
	}

	for {

		if time.Now().After(lastCronExecTime.Add(cronFrequency)) {
			core.DeleteExpiredKeys()
			lastCronExecTime = time.Now()
		}

		nevents, err := syscall.Kevent(KqueueFD, nil, events, nil)
		if err != nil {
			continue
		}

		for i := 0; i < nevents; i++ {
			fd := int(events[i].Ident)

			if fd == serverFD {
				// accept the incoming connection from a client
				nfd, _, err := syscall.Accept(serverFD)

				if err != nil {
					log.Println("err", err)
					continue
				}

				concurrent_clients += 1
				syscall.SetNonblock(nfd, true)

				var socketClientEvent syscall.Kevent_t
				syscall.SetKevent(&socketClientEvent, nfd, syscall.EVFILT_READ, syscall.EV_ADD)

				if _, err := syscall.Kevent(KqueueFD, []syscall.Kevent_t{socketClientEvent}, nil, nil); err != nil {
					log.Fatal(err)
				}
			} else {
				comm := core.FDComm{Fd: fd}
				cmds, err := readCommands(comm)

				if err != nil {
					syscall.Close(fd)
					concurrent_clients -= 1

					continue
				}
				respond(cmds, comm)
			}

		}
	}
}
