package core

import (
	"MyRedis/config"
	"fmt"
	"log"
	"os"
	"strings"
)

func dumpKey(fp *os.File, k string, obj *Obj) {
	cmd := fmt.Sprintf("SET %s %s", k, obj.Value)
	tokens := strings.Split(cmd, " ")
	fp.Write(Encode(tokens, false))
}

func DumpAllAOF() {
	fp, err := os.OpenFile(config.AOFFile, os.O_CREATE | os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Print("error", err)
		return 
	}

	log.Println("Rewriting AOF file at", config.AOFFile)
	for k, obj := range(store) {
		dumpKey(fp, k, obj)
	}
	log.Println("AOF file rewrite complete")
}