package httpserv

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/freundallein/resender/data"
	"github.com/freundallein/resender/producers"
)

// Index - main http handler
func Index(options *Options) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println("[server] POST from", r.RemoteAddr)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var pkg data.Package
		err := decoder.Decode(&pkg)
		if err != nil {
			log.Println("[parse]", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		uid := strconv.FormatUint(options.Gen.GetId(), 10)
		log.Println("[server] Received", uid, pkg)
		for _, producer := range options.Producers {
			go func(prd producers.Producer) {
				if err := prd.Validate(pkg); err != nil {
					log.Println(prd.GetName(), uid, err)
					return
				}
				data, err := pkg.Bytes()
				if err != nil {
					log.Println("[data]", err)
				}
				prd.Produce(uid, data)
			}(producer)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(uid))
	}
}
