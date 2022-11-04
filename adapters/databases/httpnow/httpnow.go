package httpnow

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type HttpNow struct {
	lock        *sync.Mutex
	lastRequest []byte
}

func NewHttpNow(listenPort int) *HttpNow {

	res := &HttpNow{
		lock: &sync.Mutex{},
	}

	http.HandleFunc("/last", res.handleLastRequest)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
		if err != nil {
			log.Fatalf("failed to start HTTP listener on port %d: %s", listenPort, err)
		}
	}()

	return res
}

func (h *HttpNow) handleLastRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if h.lastRequest == nil {
		w.Write([]byte("{}"))
		return
	}

	h.lock.Lock()
	w.Write(h.lastRequest)
	h.lock.Unlock()
}

func (h *HttpNow) InsertRecord(measurement map[string]interface{}) error {
	measurement["Timestamp"] = time.Now().UnixNano() / int64(time.Second)

	h.lock.Lock()
	h.lastRequest, _ = json.MarshalIndent(measurement, "", "  ")
	h.lock.Unlock()

	return nil
}
