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

func toArrayStrings(ai []interface{}) ([]string, error) {
	as := make([]string, len(ai))
	for i := range ai {
		as[i] = ai[i].(string)
	}
	return  as, nil
}

func readCommands(connection io.ReadWriter) (core.RedisCmds, error) {
	var buf []byte = make([]byte, 512)
	n, err := connection.Read(buf[:]) // fires read system call - Blocking call
	if err != nil {
		return nil, err
	}

	values, err := core.Decode(buf[:n])

	if err != nil {
		return nil, err
	}

	var cmds []*core.RedisCmd = make([]*core.RedisCmd, 0)

	for _, value := range values {
		tokens, err := toArrayStrings(value.([]interface{}))

		if err != nil {
			return  nil, err
		}

		cmds = append(cmds, &core.RedisCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}

	return cmds, nil
}

func respondError(err error, connection io.ReadWriter) {
	connection.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(commands core.RedisCmds, connection io.ReadWriter) {
	core.EvalAndRespond(commands, connection)
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
			commds, err := readCommands(connection)

			if err != nil {
				connection.Close()
				connected_clients -= 1

				log.Println("Client disconnected", connection.RemoteAddr(), "concurrent clients", connected_clients)

				if err == io.EOF {
					break
				}
				log.Println("Error: ", err)
			}
			respond(commds, connection)
		}
	}
}
