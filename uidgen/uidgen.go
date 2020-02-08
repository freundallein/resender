package uidgen

import (
	"sync/atomic"
	"time"
)

//Generator - generate unique uint64 ids
type Generator struct {
	machineId uint8
	counter   uint64
}

// New - constructor
func New(id uint8) *Generator {
	return &Generator{machineId: id, counter: 0}
}

//GetId - return unique uint64 id
func (gen *Generator) GetId() uint64 {
	var id uint64
	epoch := uint64(time.Now().UTC().UnixNano()) / 1000000
	id = epoch << (64 - 42)                      // 42 bits for epoch
	id |= uint64(gen.machineId) << (64 - 42 - 8) // 8 bit for machine id
	id |= gen.counter % (1<<14 - 1)              // 14 bit for counter
	atomic.AddUint64(&gen.counter, 1)
	return id
}
