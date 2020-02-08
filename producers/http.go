package producers

import (
	"bytes"
	"log"
	"net/http"

	"github.com/freundallein/resender/data"
)

// HttpProducer - send pacakges to external url
type HttpProducer struct {
	Name string
	Url  string
}

//NewHttp - constructor
func NewHttp(url string) *HttpProducer {
	return &HttpProducer{
		Name: "[http]",
		Url:  url,
	}
}

// GetName - getter
func (p *HttpProducer) GetName() string {
	return p.Name
}

//Validate - findout if we need to send it, also fill empty fields
func (p *HttpProducer) Validate(pkg data.Package) error {
	if len(pkg.ProbeRequests) == 0 {
		return data.ErrNoProbeRequests
	}
	if pkg.ApID == "" {
		return data.ErrInvalidApID
	}
	for _, probe := range pkg.ProbeRequests {
		if probe.BSSID == "" {
			probe.BSSID = "FF-FF-FF-FF-FF-FF"
		}
		if probe.SSID == "" {
			probe.SSID = "UNKNOWN"
		}

	}
	return nil
}

// Produce - package sending
func (p *HttpProducer) Produce(uid string, data []byte) {
	log.Println(p.Name, uid, "sending to", p.Url)
	req, err := http.NewRequest("POST", p.Url, bytes.NewBuffer(data))
	if err != nil {
		log.Println(p.Name, uid, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println(p.Name, uid, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(p.Name, uid, "got", resp.StatusCode)
		return
	}
	log.Println(p.Name, uid, "sent")
}
