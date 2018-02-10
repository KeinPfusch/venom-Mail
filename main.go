package main

import (
	"fmt"
	"os"

	lib "venom-Mail/lib"

	upnp "venom-Mail/upnp"

	"os/signal"
)

func init() {

	if (os.Getuid() == 0) && (os.Getgid() != 0) {
		fmt.Println("[MAIN] ROOT SHOULD NOT TO RUN THIS PROGRAM")
		os.Exit(1)
	}

	lib.LogEngineStart()
	upnp.UPNP_Worker_Start()

}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	os.Exit(0)

}
