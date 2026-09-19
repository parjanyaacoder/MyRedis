package server 

import (
	"io"
	"net"
	"log"
	"strconv"
	"MyRedis/config"
)

func readCommand(connection net.Conn) (string, error)  {
	var buf []byte = make([]byte, 512)
	n, err := connection.Read(buf[:]) // fires read system call - Blocking call
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(command string, connection net.Conn) error {
	if _, err := connection.Write([]byte(command)); err != nil {
		return err;
	}
	return nil;
}

func RunSyncTCPServer() {
	log.Println("starting a synchronous TCP server on", config.Host, config.Port)

	var connected_clients = 0; // concurrenct connected clients

	listener, err := net.Listen("tcp", config.Host +":"+strconv.Itoa(config.Port)) // Starting a tcp server
	if err != nil {
		panic(err)
	}

	for {
		connection, err := listener.Accept() // Blocking call 

		if err != nil {
			panic(err)
		}

	connected_clients += 1
	log.Println("client connected with address: ", connection.RemoteAddr(), "concurrent clients", connected_clients);

		for { 
			commd, err := readCommand(connection)

			if err != nil {
				connection.Close();
				connected_clients -= 1;

				log.Println("Client disconnected", connection.RemoteAddr(), "concurrent clients", connected_clients);

				if err == io.EOF {
					break;
				}
				log.Println("Error: ", err);
			}

			log.Println("Command: ", commd);

			if err = respond(commd, connection); err != nil {
				log.Println("err write: ", err);
			}

		}
	}
}