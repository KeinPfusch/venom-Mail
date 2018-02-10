package local

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"time"
	lib "venom-Mail/lib"
)

var SmtpMessage string
var dataPhase bool

const (
	ehloRegexp        = "(?i)^EHLO[ ]+(.*)"
	mailFromRegexp    = "(?i)^MAIL[ ]+FROM:[ ](.*)(<.*@.*>)"
	rcptToRegexp      = "(?i)^RCPT[ ]+TO:[ ](.*)(<.*@venom>)"
	endOfDataRegexp   = "^.$"
	startOfDataRegexp = "(?i)^DATA$"
	quitRegexp        = "(?i)^QUIT$"
)

func SMTP_Interpret(conn net.Conn) {

	dataPhase = false

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
			dataPhase = false
			continue
		}

		if matches, _ := regexp.MatchString(mailFromRegexp, message); matches == true {
			log.Printf("[INFO] SMTP %s from %s ", message, remote_client)
			re, _ := regexp.Compile(mailFromRegexp)
			match := re.FindStringSubmatch(message)
			log.Printf("[INFO] SMTP from %s ,email: %s ", match[1], match[2])
			conn.Write([]byte("250 2.0.0 " + match[1] + " ...Sender Ok"))
			SmtpMessage = SmtpMessage + "From: " + match[1] + " " + match[2] + "\r\n"
			dataPhase = false
			continue
		}

		if matches, _ := regexp.MatchString(rcptToRegexp, message); matches == true {
			log.Printf("[INFO] SMTP %s from %s ", message, remote_client)
			re, _ := regexp.Compile(rcptToRegexp)
			match := re.FindStringSubmatch(message)
			log.Printf("[INFO] SMTP RCPT TO to %s ,ToxID: %s ", match[1], match[2])
			conn.Write([]byte("250 " + match[2] + "... Recipient ok"))
			SmtpMessage = SmtpMessage + "To: " + match[1] + " " + match[2] + "\r\n"
			dataPhase = false
			continue
		}

		if matches, _ := regexp.MatchString(startOfDataRegexp, message); matches == true {
			log.Printf("[INFO] SMTP %s from %s ", message, remote_client)
			log.Printf("[INFO] SMTP DATA PHASE STARTED")
			conn.Write([]byte("354 Enter mail, end with \".\" on a line by itself"))
			dataPhase = true
			continue
		}

		if matches, _ := regexp.MatchString(endOfDataRegexp, message); matches == true {
			log.Printf("[INFO] SMTP %s from %s ", message, remote_client)
			log.Printf("[INFO] SMTP DATA PHASE STOPPED")
			seq := lib.RandSeq(32)
			conn.Write([]byte("250 " + seq + " Message accepted for delivery"))
			SmtpMessage = SmtpMessage + "Message-ID: " + seq + "\r\n"
			dataPhase = false
			continue
		}

		if dataPhase == true {

			SmtpMessage += message

		}

		if matches, _ := regexp.MatchString(quitRegexp, message); matches == true {
			log.Printf("[INFO] SMTP %s from %s ", message, remote_client)
			log.Printf("[INFO] SMTP CONN QUIT")
			conn.Write([]byte("221 venom closing connection"))
			dataPhase = false
			lib.SpoolWrite(lib.VenomQueue, SmtpMessage)
			break
		}

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
