package uidgen

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	observed := New(uint8(1))
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf(&Generator{})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
	if observed.machineId != 1 {
		t.Error("Expected", 1, "got", observed.machineId)
	}
	expectedCounter := uint64(0)
	if reflect.TypeOf(observed.counter) != reflect.TypeOf(expectedCounter) {
		t.Error("Expected", expectedCounter, "got", reflect.TypeOf(observed.counter))
	}
}

func TestGetId(t *testing.T) {
	N := 100000
	seen := map[uint64]struct{}{}
	gen := New(uint8(1))
	for i := 0; i < N; i++ {
		val := gen.GetId()
		_, ok := seen[val]
		if ok {
			t.Error("id is not unique")
		}
		seen[val] = struct{}{}
	}

}
