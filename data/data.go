package data

import (
	"encoding/json"
	"errors"
)

var (
	// Err... data errors
	ErrNoProbeRequests = errors.New("no probe requests")
	ErrInvalidApID     = errors.New("invalid ap_id")
	ErrInvalidMAC      = errors.New("invalid mac")
	ErrInvalidBSSID    = errors.New("invalid bssid")
)

// Probe - data from equipment
type Probe struct {
	MAC       string `json:"mac"`
	Timestamp string `json:"timestamp"`
	BSSID     string `json:"bssid"`
	SSID      string `json:"ssid"`
}

// Package - collection of probes
type Package struct {
	ApID          string   `json:"ap_id"`
	ProbeRequests []*Probe `json:"probe_requests"`
}

// PatchFields - fill empty fields
func (p *Package) PatchFields() {
	for _, probe := range p.ProbeRequests {
		if probe.BSSID == "" {
			probe.BSSID = "FF-FF-FF-FF-FF-FF"
		}
		if probe.SSID == "" {
			probe.SSID = "UNKNOWN"
		}
	}
}

// Bytes - json serializer
func (p *Package) Bytes() ([]byte, error) {
	return json.Marshal(p)
}
