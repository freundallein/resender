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
