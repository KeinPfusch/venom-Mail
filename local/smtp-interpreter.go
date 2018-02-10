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
	ehloRegexp   = "(?i)^EHLO[ ]+(.*)"
	mailToRegexp = "(?i)^MAIL[ ]+FROM:(.*)<(.*@.*)>"
	rcptToRegexp = "(?i)^RCPT[ ]+TO:(.*)<(.*)@venom>.*"
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
			continue
		}

		if matches, _ := regexp.MatchString(mailToRegexp, message); matches == true {
			log.Printf("[INFO] SMTP %s from %s ", message, remote_client)
			re, _ := regexp.Compile(mailToRegexp)
			match := re.FindStringSubmatch(message)
			log.Printf("[INFO] SMTP from %s ,email: %s ", match[1], match[2])
			conn.Write([]byte("250 2.0.0 " + match[1] + " ...Sender Ok"))
			continue
		}

		if matches, _ := regexp.MatchString(rcptToRegexp, message); matches == true {
			log.Printf("[INFO] SMTP %s from %s ", message, remote_client)
			re, _ := regexp.Compile(rcptToRegexp)
			match := re.FindStringSubmatch(message)
			log.Printf("[INFO] SMTP RCPT TO to %s ,ToxID: %s ", match[1], match[2])
			conn.Write([]byte("250 " + match[2] + "... Recipient ok"))
			continue
		}

		// TO BE COMPLETED WITH THE SMTP VOCABULARY

		if message == "" {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("[INFO] SMTP BULLSHIT >%s< from %s ", message, remote_client)
		time.Sleep(1 * time.Second)
		conn.Write([]byte("500 Command not understood\r\n"))
		break

	}
	conn.Close()
}
