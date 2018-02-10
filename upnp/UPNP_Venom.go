package upnp

import (
	"log"
	"strings"
	"time"
)

func init() {

	go UPNP_Worker_Start()

}

func UPNP_Worker_Start() {

	upnp_renew := time.NewTicker(5 * time.Minute)
	mapping_DHT := new(Upnp)

	ClusterPort := 35000

	err := mapping_DHT.SearchGateway()
	if err != nil {
		log.Printf("[UPnP] Problem getting the gateway:  %s...", err.Error())
	} else {
		log.Printf("[UPnP] Local ip address: %s", mapping_DHT.LocalHost)
		torn := strings.Split(mapping_DHT.Gateway.Host, ":")
		log.Printf("[UPnP] Gateway ip address is %s on port: %s", torn[0], torn[1])
	}

	err = mapping_DHT.ExternalIPAddr()
	if err != nil {
		log.Printf("[UPnP] Problem getting my external IP:  %s...", err.Error())

	} else {
		log.Printf("[UPnP] WAN ip address: %s", mapping_DHT.GatewayOutsideIP)

	}

	log.Printf("[UPnP] UPnP on TCP %d...", ClusterPort)

	for {

		if err = mapping_DHT.AddPortMapping(ClusterPort, ClusterPort, "TCP"); err == nil {
			log.Printf("[UPnP] UPnP redirect TCP %d : no errors from %s", ClusterPort, mapping_DHT.Gateway.Host)
		} else {
			log.Printf("[UPnP] No UPnP on TCP %d:  %s", ClusterPort, err.Error())
		}

		<-upnp_renew.C
		log.Printf("[UPnP] Renew the UPnP lease")
	}

}

func UPNP_Engine_Start() {

	log.Printf("[UPnP] Engine start")

}
