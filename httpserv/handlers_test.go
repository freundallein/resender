package httpserv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/freundallein/resender/data"
	"github.com/freundallein/resender/producers"
	"github.com/freundallein/resender/uidgen"
)

func TestIndexClosure(t *testing.T) {
	handlerFunc := Index(nil)
	observedType := reflect.TypeOf(handlerFunc)
	expectedType := reflect.TypeOf(func(w http.ResponseWriter, r *http.Request) {})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}

}
func TestIndex405(t *testing.T) {
	handlerFunc := Index(nil)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestIndex404(t *testing.T) {
	handlerFunc := Index(nil)
	req, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Error(err)
	}
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
func TestIndex400(t *testing.T) {
	handlerFunc := Index(nil)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte("i am bad")))
	if err != nil {
		t.Error(err)
	}
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
func TestIndex200(t *testing.T) {
	opts := &Options{
		Gen: uidgen.New(1),
	}
	handlerFunc := Index(opts)
	body := `{"ap_id":"123","probe_requests":[{"mac":"1","timestamp":"2","bssid":"3","ssid":"4"}]}`
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Error(err)
	}
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestIndexWithHttpProducer(t *testing.T) {
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
	opts := &Options{
		Gen:       uidgen.New(1),
		Producers: []producers.Producer{producers.NewHttp(ts.URL)},
	}
	handlerFunc := Index(opts)
	body := `{"ap_id":"A8-F9-4B-B6-87-FF","probe_requests":[{"mac":"88-1D-FC-DF-6F-C1","timestamp":"1579782767"}]}`
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Error(err)
	}
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	time.Sleep(100 * time.Millisecond)
	if received.ApID != "A8-F9-4B-B6-87-FF" {
		t.Error("Expected A8-F9-4B-B6-87-FF got", received.ApID)
	}
	receivedProbe := received.ProbeRequests[0]
	if receivedProbe.MAC != "88-1D-FC-DF-6F-C1" {
		t.Error("Expected 88-1D-FC-DF-6F-C1 got", receivedProbe.MAC)
	}
	if receivedProbe.Timestamp != "1579782767" {
		t.Error("Expected 1579782767 got", receivedProbe.Timestamp)
	}
	if receivedProbe.BSSID != "FF-FF-FF-FF-FF-FF" {
		t.Error("Expected FF-FF-FF-FF-FF-FF got", receivedProbe.BSSID)
	}
	if receivedProbe.SSID != "UNKNOWN" {
		t.Error("Expected UNKNOWN got", receivedProbe.SSID)
	}
}

// // Index - main http handler
// func Index(options *Options) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		for _, producer := range options.Producers {
// 			go func(prd producers.Producer) {
// 				if err := prd.Validate(pkg); err != nil {
// 					log.Println(prd.GetName(), uid, err)
// 					return
// 				}
// 				data, err := pkg.Bytes()
// 				if err != nil {
// 					log.Println("[data]", err)
// 				}
// 				prd.Produce(uid, data)
// 			}(producer)
// 		}
// 	}
// }
