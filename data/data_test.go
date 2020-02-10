package data

import (
	"testing"
)

func TestBytes(t *testing.T) {
	p := &Package{
		ApID: "123",
		ProbeRequests: []*Probe{
			&Probe{
				MAC:       "1",
				Timestamp: "2",
				BSSID:     "3",
				SSID:      "4",
			},
		},
	}
	observed, err := p.Bytes()
	if err != nil {
		t.Error(err)
	}
	expected := []byte(`{"ap_id":"123","probe_requests":[{"mac":"1","timestamp":"2","bssid":"3","ssid":"4"}]}`)
	for i := range observed {
		if observed[i] != expected[i] {
			t.Error("Expected", string(expected), "got", string(observed))
		}
	}
}

func TestPatchFieldsEmpty(t *testing.T) {
	observed := []*Probe{
		&Probe{
			MAC:       "1",
			Timestamp: "2",
			BSSID:     "",
			SSID:      "",
		},
		&Probe{
			MAC:       "1",
			Timestamp: "2",
			BSSID:     "",
			SSID:      "",
		},
	}
	p := &Package{
		ApID: "123",
		ProbeRequests: observed,
	}
	p.PatchFields()
	expectedBSSID := "FF-FF-FF-FF-FF-FF"
	expectedSSID := "UNKNOWN"
	for i := range observed {
		if observed[i].BSSID != expectedBSSID {
			t.Error("Expected", expectedBSSID, "got", observed[i].BSSID)
		}
		if observed[i].SSID != expectedSSID {
			t.Error("Expected", expectedSSID, "got", observed[i].SSID)
		}
	}
}

func TestPatchFieldsFilled(t *testing.T) {
	observed := []*Probe{
		&Probe{
			MAC:       "1",
			Timestamp: "2",
			BSSID:     "3",
			SSID:      "4",
		},
		&Probe{
			MAC:       "1",
			Timestamp: "2",
			BSSID:     "3",
			SSID:      "4",
		},
	}
	p := &Package{
		ApID: "123",
		ProbeRequests: observed,
	}
	p.PatchFields()
	expectedBSSID := "3"
	expectedSSID := "4"
	for i := range observed {
		if observed[i].BSSID != expectedBSSID {
			t.Error("Expected", expectedBSSID, "got", observed[i].BSSID)
		}
		if observed[i].SSID != expectedSSID {
			t.Error("Expected", expectedSSID, "got", observed[i].SSID)
		}
	}
}