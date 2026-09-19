package core

import (
	"errors"
	"io"
)



func evalPing(args []string, connection io.ReadWriter) error {
	var b []byte;

	if len(args) >=2 {
		return  errors.New("ERR Wrong number of arguments for the 'ping' command")
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} 
	if(len(args) == 1) {
		b = Encode(args[0], false)
	}

	_, err := connection.Write(b)

	return  err
}


func EvalAndRespond(command *RedisCmd, connection io.ReadWriter) error {
	switch command.Cmd {
	case "PING": return evalPing(command.Args, connection)
	default: return evalPing(command.Args, connection)	
	}
}