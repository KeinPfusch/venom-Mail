package upnp

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

type DeviceDesc struct {
	upnp *Upnp
}

func (this *DeviceDesc) Send() bool {
	request := this.BuildRequest()
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return false
	}

	resultBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		this.resolve(string(resultBody))
		return true
	}
	return false
}
func (this *DeviceDesc) BuildRequest() *http.Request {

	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("User-Agent", "preston")
	header.Set("Host", this.upnp.Gateway.Host)
	header.Set("Connection", "keep-alive")

	request, err := http.NewRequest("GET", "http://"+this.upnp.Gateway.Host+this.upnp.Gateway.DeviceDescUrl, nil)

	if err != nil {
		return nil
	}

	request.Header = header

	return request
}

func (this *DeviceDesc) resolve(resultStr string) {
	inputReader := strings.NewReader(resultStr)

	lastLabel := ""

	ISUpnpServer := false

	IScontrolURL := false
	var controlURL string //`controlURL`

	decoder := xml.NewDecoder(inputReader)
	for t, err := decoder.Token(); err == nil && !IScontrolURL; t, err = decoder.Token() {
		switch token := t.(type) {

		case xml.StartElement:
			if ISUpnpServer {
				name := token.Name.Local
				lastLabel = name
			}

		case xml.EndElement:

		case xml.CharData:

			content := string([]byte(token))

			if content == this.upnp.Gateway.ServiceType {
				ISUpnpServer = true
				continue
			}

			if ISUpnpServer {
				switch lastLabel {
				case "controlURL":

					controlURL = content
					IScontrolURL = true
				case "eventSubURL":
					// eventSubURL = content
				case "SCPDURL":
					// SCPDURL = content
				}
			}
		default:
			// ...
		}
	}
	this.upnp.CtrlUrl = controlURL
}
