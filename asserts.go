package main

import (
	"log"
	"os"
	"runtime"
)

func assertLinux() {
	if runtime.GOOS != "linux" {
		log.Printf("Error found: This service can only run on linux. You are running %s.", runtime.GOOS)
		os.Exit(4)
	}
}

func assert(err error) {
	if err != nil {
		log.Printf("Error found: %s", err.Error())
	}
}

func assertExit(msg string, err error, exitCode int) {
	if err != nil {
		log.Printf("Error found: %s", msg)
		log.Printf(err.Error())
		os.Exit(exitCode)
	}
}

func logMsg(msg string) {
	log.Println(msg)
}
