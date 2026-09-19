package server

import (
	"MyRedis/config"
	"MyRedis/core"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func readCommand(connection io.ReadWriter) (*core.RedisCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := connection.Read(buf[:]) // fires read system call - Blocking call
	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeArrayStrings(buf[:n])

	if err != nil {
		return nil, err
	}

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, connection io.ReadWriter) {
	connection.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(command *core.RedisCmd, connection io.ReadWriter) {
	err := core.EvalAndRespond(command, connection)
	if err != nil {
		respondError(err, connection)
	}
}

func RunSyncTCPServer() {
	log.Println("starting a synchronous TCP server on", config.Host, config.Port)

	var connected_clients = 0 // concurrenct connected clients

	listener, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port)) // Starting a tcp server
	if err != nil {
		panic(err)
	}

	for {
		connection, err := listener.Accept() // Blocking call

		if err != nil {
			panic(err)
		}

		connected_clients += 1
		log.Println("client connected with address: ", connection.RemoteAddr(), "concurrent clients", connected_clients)

		for {
			commd, err := readCommand(connection)

			if err != nil {
				connection.Close()
				connected_clients -= 1

				log.Println("Client disconnected", connection.RemoteAddr(), "concurrent clients", connected_clients)

				if err == io.EOF {
					break
				}
				log.Println("Error: ", err)
			}
			respond(commd, connection)
		}
	}
}
