package core

import (
	"errors"
	"io"
	"strconv"
	"time"
)

var RESP_NIL []byte = []byte("$-1\r\n")

func evalPing(args []string, connection io.ReadWriter) error {
	var b []byte

	if len(args) >= 2 {
		return errors.New("ERR Wrong number of arguments for the 'ping' command")
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	}
	if len(args) == 1 {
		b = Encode(args[0], false)
	}

	_, err := connection.Write(b)

	return err
}

func evalSet(args []string, connection io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(error) ERR wrong number of arguments for 'set' command")
	}

	var key, value string
	var exDurationMs int64 = -1

	key, value = args[0], args[1]
	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			{
				i += 1
				if i == len(args) {
					return errors.New("(error) ERR syntax error")
				}

				exDurationSec, err := strconv.ParseInt(args[3], 10, 64)
				if err != nil {
					return errors.New("(error) ERR value if not an integer or out of range")
				}

				exDurationMs = exDurationSec * 1000
			}
		default:
			return errors.New("(error) ERR syntax error")
		}
	}

	Put(key, NewObj(value, exDurationMs))
	connection.Write([]byte("+OK\r\n"))
	return nil
}

func evalGet(args []string, connection io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'get' command")
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		connection.Write(RESP_NIL)
		return nil
	}

	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		connection.Write(RESP_NIL)
		return nil
	}

	connection.Write(Encode(obj.Value, false))
	return nil
}

func evalTtl(args []string, connection io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong number of arguments for 'ttl' command")
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		connection.Write([]byte(":-2\r\n"))
		return nil
	}

	if obj.ExpiresAt == -1 {
		connection.Write([]byte(":-1\r\n"))
		return nil
	}

	durationMs := obj.ExpiresAt - time.Now().UnixMilli()

	if durationMs < 0 {
		connection.Write([]byte(":-2\r\n"))
		return nil
	}

	connection.Write(Encode(int64(durationMs/1000), false))
	return nil
}

func evalDel(args []string, connection io.ReadWriter) error {
	var countDeleted int64 = 0
	for _, key := range args {

		if ok := Del(key); ok {
			countDeleted++
		}
	}

	connection.Write(Encode(countDeleted, false))
	return nil
}

func evalExpire(args []string, connection io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(error) ERR wrong number of arguments for 'expire' command")
	}

	var key string = args[0]

	exDurationSec, err := strconv.ParseInt(args[1], 10, 64)

	if err != nil {
		return errors.New("(error) ERR value if not an integer or out of range")
	}

	obj := Get(key)

	if obj == nil {
		connection.Write([]byte(":0\r\n"))
		return nil
	}

	obj.ExpiresAt = time.Now().UnixMilli() + exDurationSec*1000
	connection.Write([]byte(":1\r\n"))

	return nil
}

func EvalAndRespond(command *RedisCmd, connection io.ReadWriter) error {
	switch command.Cmd {
	case "PING":
		return evalPing(command.Args, connection)
	case "SET":
		return evalSet(command.Args, connection)
	case "GET":
		return evalGet(command.Args, connection)
	case "TTL":
		return evalTtl(command.Args, connection)
	case "DEL":
		return evalDel(command.Args, connection)
	case "EXPIRE":
		return evalExpire(command.Args, connection)
	default:
		return evalPing(command.Args, connection)
	}
}
