package upnp

import (
	"log"
	"net"
	"strings"
	"time"
	// "net/http"
)

type Gateway struct {
	GatewayName   string
	Host          string
	DeviceDescUrl string
	Cache         string
	ST            string
	USN           string
	deviceType    string
	ControlURL    string
	ServiceType   string
}





type SearchGateway struct {
	searchMessage string
	upnp          *Upnp
}

func (this *SearchGateway) Send() bool {
	this.buildRequest()
	c := make(chan string)
	go this.send(c)
	result := <-c
	if result == "" {

		this.upnp.Active = false
		return false
	}
	this.resolve(result)

	this.upnp.Gateway.ServiceType = "urn:schemas-upnp-org:service:WANIPConnection:1"
	this.upnp.Active = true
	return true
}
func (this *SearchGateway) send(c chan string) {

	var conn *net.UDPConn
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	go func(conn *net.UDPConn) {
		defer func() {
			if r := recover(); r != nil {

			}
		}()

		time.Sleep(time.Second * 3)
		c <- ""
		conn.Close()
	}(conn)
	remotAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	if err != nil {
		log.Println("[UPnP] Multicast address is ok")
	}
	locaAddr, err := net.ResolveUDPAddr("udp", this.upnp.LocalHost+":50700")

	if err != nil {
		log.Println("[UpnP] IP Format is Correct")
	}
	conn, err = net.ListenUDP("udp", locaAddr)
	defer conn.Close()
	if err != nil {
		log.Println("[UPnP] Error in listening to UDP")
	}

	_, err = conn.WriteToUDP([]byte(this.searchMessage), remotAddr)
	if err != nil {
		log.Println("[UPnP] Wrong Multicast Address")
	}
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Println("[UPnP] Can't read UDP packets")
	}

	result := string(buf[:n])
	c <- result
}
func (this *SearchGateway) buildRequest() {
	this.searchMessage = "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"ST: urn:schemas-upnp-org:service:WANIPConnection:1\r\n" +
		"MAN: \"ssdp:discover\"\r\n" + "MX: 3\r\n\r\n"
}

func (this *SearchGateway) resolve(result string) {
	this.upnp.Gateway = &Gateway{}

	lines := strings.Split(result, "\r\n")
	for _, line := range lines {

		nameValues := strings.SplitAfterN(line, ":", 2)
		if len(nameValues) < 2 {
			continue
		}
		switch strings.ToUpper(strings.Trim(strings.Split(nameValues[0], ":")[0], " ")) {
		case "ST":
			this.upnp.Gateway.ST = nameValues[1]
		case "CACHE-CONTROL":
			this.upnp.Gateway.Cache = nameValues[1]
		case "LOCATION":
			urls := strings.Split(strings.Split(nameValues[1], "//")[1], "/")
			this.upnp.Gateway.Host = urls[0]
			this.upnp.Gateway.DeviceDescUrl = "/" + urls[1]
		case "SERVER":
			this.upnp.Gateway.GatewayName = nameValues[1]
		default:
		}
	}
}
