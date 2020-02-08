package data

import (
	"encoding/json"
	"errors"
)

var (
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

// Bytes - json serializer
func (p *Package) Bytes() ([]byte, error) {
	return json.Marshal(p)
}
