package upnp

import (
	"io/ioutil"

	"net/http"
	"strconv"
	"strings"
)

type DelPortMapping struct {
	upnp *Upnp
}

func (this *DelPortMapping) Send(remotePort int, protocol string) bool {
	request := this.buildRequest(remotePort, protocol)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return false
	}

	resultBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		// log.Println(string(resultBody))
		this.resolve(string(resultBody))
		return true
	}
	return false
}
func (this *DelPortMapping) buildRequest(remotePort int, protocol string) *http.Request {

	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#DeletePortMapping"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")

	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:DeletePortMapping`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}
	childList1 := Node{Name: "NewExternalPort", Content: strconv.Itoa(remotePort)}
	childList2 := Node{Name: "NewProtocol", Content: protocol}
	childList3 := Node{Name: "NewRemoteHost"}
	childTwo.AddChild(childList1)
	childTwo.AddChild(childList2)
	childTwo.AddChild(childList3)
	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	bodyStr := body.BuildXML()

	request, err := http.NewRequest("POST", "http://"+this.upnp.Gateway.Host+this.upnp.CtrlUrl,
		strings.NewReader(bodyStr))

	if err != nil {
		return nil
	}
	request.Header = header
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(bodyStr))))
	return request
}

func (this *DelPortMapping) resolve(resultStr string) {
}
