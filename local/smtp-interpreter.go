package local

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"time"
)

type SmtpMessage string

const (
	ehloRegexp   = "(?i)^EHLO(.*)"
	mailToRegexp = "(?i)^MAIL[ ]+FROM:(.*)"
)

func SMTP_Interpret(conn net.Conn) {

	remote_client := conn.RemoteAddr()
	greetings := "220 Venom ESMTP ToxGateway READY"
	conn.Write([]byte(greetings + "\r\n"))
	for {
		linea := make([]byte, 1024)

		linea, _, _ = bufio.NewReader(conn).ReadLine()

		message := string(linea)

		// decides WTF to do with the string

		if matches, _ := regexp.MatchString(ehloRegexp, message); matches == true {

			re, _ := regexp.Compile(ehloRegexp)
			match := re.FindStringSubmatch(message)

			log.Printf("[INFO] SMTP %s from %s ", message, remote_client)

			conn.Write([]byte("250-Venom Hello" + match[1] + " ,pleased to meet you \r\n"))
			conn.Write([]byte("250-8BITMIME\r\n"))
			conn.Write([]byte("250-SIZE 36700160\r\n"))
			break
		}

		if matches, _ := regexp.MatchString(mailToRegexp, message); matches == true {
			log.Printf("[INFO] SMTP %s from %s ", message, remote_client)
			conn.Write([]byte("205 closing connection - goodbye!"))
			conn.Close()
			break
		}

		// TO BE COMPLETED WITH THE SMTP VOCABULARY

		if message == "" {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("[INFO] SMTP BULLSHIT >%s< from %s ", message, remote_client)
		time.Sleep(1 * time.Second)
		conn.Write([]byte("500 Command not understood\r\n"))

	}
	conn.Close()
}
