package local

import (
	"log"
	"net"
	lib "venom-Mail/lib"
)

func init() {

	go SMTP_Frontend()

}

func SMTP_Frontend() {

	// setting up the tcp connection

	listenSmtp := "127.0.0.1:" + lib.VConfig["smtpport"]

	ln, err := net.Listen("tcp", listenSmtp)
	if err == nil {
		log.Printf("[NNTP] TCP listening at %s ", listenSmtp)

	} else {
		log.Printf("[WTF] TCP CANNOT listen at %s. SYSADMIIIIN!!", "127.0.0.1:11119")
		return
	}

	defer ln.Close()

	for {

		// start listening at it

		server, err := ln.Accept()
		tcp_client := server.RemoteAddr()

		if err == nil {

			log.Printf("[SMTP][SERVER][INFO] SMTP accepted connection from %s ", tcp_client)

		} else {
			log.Printf("[WTF] SMTP something went wrong at %s. SYSADMIIIIN!!", listenSmtp)
		}

		// start the NNTP interpreter in background.
		go SMTP_Interpret(server)

	}

}

func SMTP_Engine_Start() {

	log.Printf("[SMTP] SMTP Engine Starting")

}
