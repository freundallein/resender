package producers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/freundallein/resender/data"
)

func TestNew(t *testing.T) {
	expectedUrl := "test/url"
	observed := NewHttp(expectedUrl)
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf(&HttpProducer{})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
	if observed.Url != expectedUrl {
		t.Error("Expected", expectedUrl, "got", observed.Url)
	}
}

func TestGetName(t *testing.T) {
	expectedName := "[http]"
	observed := NewHttp("test/url")
	if observed.GetName() != expectedName {
		t.Error("Expected", expectedName, "got", observed.GetName())
	}
}
func TestValidate(t *testing.T) {
	observed := NewHttp("test/url")
	pkg := data.Package{
		ApID: "123",
		ProbeRequests: []*data.Probe{
			&data.Probe{
				MAC:       "1",
				Timestamp: "2",
				BSSID:     "",
				SSID:      "",
			},
		},
	}
	err := observed.Validate(pkg)
	if err != nil {
		t.Error(err)
	}
	probe := pkg.ProbeRequests[0]
	expectedBSSID := "FF-FF-FF-FF-FF-FF"
	if probe.BSSID != expectedBSSID {
		t.Error("Expected", expectedBSSID, "got", probe.BSSID)
	}
	expectedSSID := "UNKNOWN"
	if probe.SSID != expectedSSID {
		t.Error("Expected", expectedSSID, "got", probe.SSID)
	}
}
func TestValidateEmpty(t *testing.T) {
	observed := NewHttp("test/url")
	pkg := data.Package{
		ApID:          "123",
		ProbeRequests: []*data.Probe{},
	}
	err := observed.Validate(pkg)
	if err == nil {
		t.Error("Should be data.ErrNoProbeRequests")
	}
	if err != data.ErrNoProbeRequests {
		t.Error("Expected", data.ErrNoProbeRequests, "got", err)
	}
}
func TestValidateInvalidApID(t *testing.T) {
	observed := NewHttp("test/url")
	pkg := data.Package{
		ApID: "",
		ProbeRequests: []*data.Probe{
			&data.Probe{
				MAC:       "1",
				Timestamp: "2",
				BSSID:     "",
				SSID:      "",
			},
		},
	}
	err := observed.Validate(pkg)
	if err == nil {
		t.Error("Should be data.ErrInvalidApID")
	}
	if err != data.ErrInvalidApID {
		t.Error("Expected", data.ErrInvalidApID, "got", err)
	}
}

func TestProduce(t *testing.T) {
	var received data.Package
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&received)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()
	observed := NewHttp(ts.URL)
	pkg := data.Package{
		ApID: "123",
		ProbeRequests: []*data.Probe{
			&data.Probe{
				MAC:       "1",
				Timestamp: "2",
				BSSID:     "3",
				SSID:      "4",
			},
		},
	}
	data, err := pkg.Bytes()
	if err != nil {
		t.Error(err)
	}
	observed.Produce("test", data)
	if received.ApID != pkg.ApID {
		t.Error("Expected", pkg.ApID, "got", received.ApID)
	}
	receivedProbe := received.ProbeRequests[0]
	expectedProbe := pkg.ProbeRequests[0]
	if receivedProbe.MAC != expectedProbe.MAC {
		t.Error("Expected", expectedProbe.MAC, "got", receivedProbe.MAC)
	}
	if receivedProbe.Timestamp != expectedProbe.Timestamp {
		t.Error("Expected", expectedProbe.Timestamp, "got", receivedProbe.Timestamp)
	}
	if receivedProbe.BSSID != expectedProbe.BSSID {
		t.Error("Expected", expectedProbe.BSSID, "got", receivedProbe.BSSID)
	}
	if receivedProbe.SSID != expectedProbe.SSID {
		t.Error("Expected", expectedProbe.SSID, "got", receivedProbe.SSID)
	}

}
