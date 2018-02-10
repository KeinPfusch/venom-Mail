package upnp

import (
	// "fmt"
	"errors"
	"log"
	"sync"
	"tribes/tools"
)

type MappingPortStruct struct {
	lock         *sync.Mutex
	mappingPorts map[string][][]int
}

func (this *MappingPortStruct) addMapping(localPort, remotePort int, protocol string) {

	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		portMapping := [][]int{one, two}
		this.mappingPorts = map[string][][]int{protocol: portMapping}
		return
	}
	portMapping := this.mappingPorts[protocol]
	if portMapping == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		this.mappingPorts[protocol] = [][]int{one, two}
		return
	}
	one := portMapping[0]
	two := portMapping[1]
	one = append(one, localPort)
	two = append(two, remotePort)
	this.mappingPorts[protocol] = [][]int{one, two}
}

func (this *MappingPortStruct) delMapping(remotePort int, protocol string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		return
	}
	tmp := MappingPortStruct{lock: new(sync.Mutex)}
	mappings := this.mappingPorts[protocol]
	for i := 0; i < len(mappings[0]); i++ {
		if mappings[1][i] == remotePort {

			break
		}
		tmp.addMapping(mappings[0][i], mappings[1][i], protocol)
	}
	this.mappingPorts = tmp.mappingPorts
}
func (this *MappingPortStruct) GetAllMapping() map[string][][]int {
	return this.mappingPorts
}

type Upnp struct {
	Active             bool
	LocalHost          string
	GatewayInsideIP    string
	GatewayOutsideIP   string
	OutsideMappingPort map[string]int
	InsideMappingPort  map[string]int
	Gateway            *Gateway
	CtrlUrl            string
	MappingPort        MappingPortStruct
}

func (this *Upnp) SearchGateway() (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			log.Println("[UPnP] Timeout", errTemp)
			err = errTemp.(error)
		}
	}(err)

	if this.LocalHost == "" {
		this.MappingPort = MappingPortStruct{
			lock: new(sync.Mutex),
		}
		this.LocalHost = tools.ReadIpFromHost()
	}
	searchGateway := SearchGateway{upnp: this}
	if searchGateway.Send() {
		return nil
	}
	return errors.New("No UPnP Gateway")
}

func (this *Upnp) deviceStatus() {

}

func (this *Upnp) deviceDesc() (err error) {
	if this.GatewayInsideIP == "" {
		if err := this.SearchGateway(); err != nil {
			return err
		}
	}
	device := DeviceDesc{upnp: this}
	device.Send()
	this.Active = true

	return
}

func (this *Upnp) ExternalIPAddr() (err error) {
	if this.CtrlUrl == "" {
		if err := this.deviceDesc(); err != nil {
			return err
		}
	}
	eia := ExternalIPAddress{upnp: this}
	eia.Send()
	return nil

}

func (this *Upnp) AddPortMapping(localPort, remotePort int, protocol string) (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			log.Println("[UPnP] Network timeout/problem caused Exception, recovering.")
			err = errTemp.(error)
		}
	}(err)
	if this.GatewayOutsideIP == "" {
		if err := this.ExternalIPAddr(); err != nil {
			log.Println("[UPnP] AddPortMapping: ", err.Error())
			return err
		}
	}
	addPort := AddPortMapping{upnp: this}
	if issuccess := addPort.Send(localPort, remotePort, protocol); issuccess {
		this.MappingPort.addMapping(localPort, remotePort, protocol)
		log.Println("[UPnP] AddPortMapping: ", issuccess)
		return nil
	} else {
		this.Active = false
		log.Println("[UPnP] AddPortMapping: this.Active = false")
		return errors.New("Adding port mapping failed")
	}
}

func (this *Upnp) DelPortMapping(remotePort int, protocol string) bool {
	delMapping := DelPortMapping{upnp: this}
	issuccess := delMapping.Send(remotePort, protocol)
	if issuccess {
		this.MappingPort.delMapping(remotePort, protocol)
		log.Println("[UPnP] Removed port:", remotePort)
	}
	return issuccess
}

func (this *Upnp) Reclaim() {
	mappings := this.MappingPort.GetAllMapping()
	tcpMapping, ok := mappings["TCP"]
	if ok {
		for i := 0; i < len(tcpMapping[0]); i++ {
			this.DelPortMapping(tcpMapping[1][i], "TCP")
		}
	}
	udpMapping, ok := mappings["UDP"]
	if ok {
		for i := 0; i < len(udpMapping[0]); i++ {
			this.DelPortMapping(udpMapping[0][i], "UDP")
		}
	}
}

func (this *Upnp) GetAllMapping() map[string][][]int {
	return this.MappingPort.GetAllMapping()
}
