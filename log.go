package main

import (
	"fmt"
	"log"
	"time"
)

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	if flagSet.Lookup("debug").Value.String() != "disable" {
		return fmt.Print("["+time.Now().UTC().Format("2006-01-02 15:04:05")+"] ", string(bytes))
	}
	return fmt.Print(string(bytes))
}

func setLogDebug() {
	if flagSet.Lookup("debug").Value.String() != "disable" {
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetFlags(0)
	}
}
