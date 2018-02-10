package main

import (
	"fmt"
	"os"

	lib "venom-Mail/lib"
	loc "venom-Mail/local"
	upnp "venom-Mail/upnp"

	"os/signal"
)

func init() {

	if (os.Getuid() == 0) && (os.Getgid() != 0) {
		fmt.Println("[MAIN] ROOT SHALL NOT TO RUN THIS PROGRAM!!!!")
		os.Exit(1)
	}

	lib.LogEngineStart()
	lib.StartConfig()
	upnp.UPNP_Engine_Start()
	loc.SMTP_Engine_Start()

}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	os.Exit(0)

}
