package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

var RESP_NIL []byte = []byte("$-1\r\n")
var RESP_OK []byte = []byte("+OK\r\n")

func evalPing(args []string) []byte {
	var b []byte

	if len(args) >= 2 {
		return Encode(errors.New("ERR Wrong number of arguments for the 'ping' command"), false)
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	}
	if len(args) == 1 {
		b = Encode(args[0], false)
	}

	return b
}

func evalSet(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'set' command"), false)
	}

	var key, value string
	var exDurationMs int64 = -1

	key, value = args[0], args[1]
	oType, oEncoding := deduceTypeEncoding(value)

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			{
				i += 1
				if i == len(args) {
					return Encode(errors.New("ERR syntax error"), false)
				}

				exDurationSec, err := strconv.ParseInt(args[3], 10, 64)
				if err != nil {
					return Encode(errors.New("ERR value if not an integer or out of range"), false)
				}

				exDurationMs = exDurationSec * 1000
			}
		default:
			return Encode(errors.New("ERR syntax error"), false)
		}
	}

	Put(key, NewObj(value, exDurationMs, oType, oEncoding))
	return RESP_OK
}

func evalGet(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'get' command"), false)
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		return RESP_NIL
	}

	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		return RESP_NIL
	}

	return Encode(obj.Value, false)
}

func evalTtl(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ttl' command"), false)
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		return []byte(":-2\r\n")
	}

	if obj.ExpiresAt == -1 {
		return []byte(":-1\r\n")
	}

	durationMs := obj.ExpiresAt - time.Now().UnixMilli()

	if durationMs < 0 {
		return []byte(":-2\r\n")
	}

	return Encode(int64(durationMs/1000), false)
}

func evalDel(args []string) []byte {
	var countDeleted int64 = 0
	for _, key := range args {

		if ok := Del(key); ok {
			countDeleted++
		}
	}

	return Encode(countDeleted, false)
}

func evalExpire(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'expire' command"), false)
	}

	var key string = args[0]

	exDurationSec, err := strconv.ParseInt(args[1], 10, 64)

	if err != nil {
		return Encode(errors.New("ERR value if not an integer or out of range"), false)
	}

	obj := Get(key)

	if obj == nil {
		return []byte(":0\r\n")
	}

	obj.ExpiresAt = time.Now().UnixMilli() + exDurationSec*1000
	return []byte(":1\r\n")
}

func evalBGREWRITEAOF(args []string) []byte {
	DumpAllAOF()
	return RESP_OK
}

func evalIncr(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'incr' command"), false)
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		obj = NewObj("0", -1, OBJ_TYPE_STRING, OBJ_ENCODING_INT)
		Put(key, obj)
	}

	if err := assertType(obj.TypeEncoding, OBJ_TYPE_STRING); err != nil {
		return Encode(err, false)
	}

	if err := assertEncoding(obj.TypeEncoding, OBJ_ENCODING_INT); err != nil {
		return Encode(err, false)
	}

	i, _ := strconv.ParseInt(obj.Value.(string), 10, 64)
	i++

	obj.Value = strconv.FormatInt(i, 10)
	return Encode(i, false)
}

func evalInfo(args []string) []byte {
	var info[] byte
	buf := bytes.NewBuffer(info)
	buf.WriteString("# Keyspace\r\n")

	for i := range KeyspaceStat {
		buf.WriteString(fmt.Sprintf("db%d:keys=%d, expires=0, avg_ttl=0\r\n", i, KeyspaceStat[i]["Keys"]))
	}
	return  Encode(buf.String(), false)
}

func evalClinet(args []string) []byte {
	return RESP_OK
}

func evalLatency(args []string) []byte {
	return Encode([]string{}, false)
}


func EvalAndRespond(commands RedisCmds, connection io.ReadWriter) {
	var response []byte
	buf := bytes.NewBuffer(response)

	for _, command := range commands {
		switch command.Cmd {
		case "PING":
			buf.Write(evalPing(command.Args))
		case "SET":
			buf.Write(evalSet(command.Args))
		case "GET":
			buf.Write(evalGet(command.Args))
		case "TTL":
			buf.Write(evalTtl(command.Args))
		case "DEL":
			buf.Write(evalDel(command.Args))
		case "EXPIRE":
			buf.Write(evalExpire(command.Args))
		case "BGREWRITEAOF":
			buf.Write(evalBGREWRITEAOF(command.Args))
		case "INCR":
			buf.Write(evalIncr(command.Args))
		case "INFO": 
			buf.Write(evalInfo(command.Args))
		case "CLIENT": 
			buf.Write(evalClinet(command.Args))
		case "LATENCY": 
			buf.Write(evalLatency(command.Args))
		default:
			buf.Write(evalPing(command.Args))
		}
	}

	connection.Write(buf.Bytes())
}
