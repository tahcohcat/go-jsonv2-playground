package main

import (
	"log"
	"math/big"
	"net/http"
	"time"

	json "encoding/json/v2"
)

type Payload struct {
	ID        int        `json:"id"`
	Timestamp time.Time  `json:"timestamp"`
	Value     *big.Float `json:"value"`
	Message   string     `json:"message"`
}

// Custom marshalers (wrapped as options)
var timeRFC3339 = json.MarshalFunc(func(t time.Time) ([]byte, error) {
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
})
var bigFloatAsString = json.MarshalFunc(func(f any) ([]byte, error) {
	if bf, ok := f.(*big.Float); ok {
		return []byte(`"` + bf.Text('g', -1) + `"`), nil
	}
	return nil, nil
})

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	p := Payload{1, time.Now(), big.NewFloat(12345.6789), "default encoding"}
	data, err := json.Marshal(p, json.StringifyNumbers(false))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func customHandler(w http.ResponseWriter, r *http.Request) {
	p := Payload{2, time.Now(), big.NewFloat(98765.4321), "custom encoding"}

	data, err := json.Marshal(p,
		json.WithMarshalers(timeRFC3339),
		json.WithMarshalers(bigFloatAsString),
		json.OmitZeroStructFields(true),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/default", defaultHandler)
	http.HandleFunc("/custom", customHandler)
	log.Println("Serving on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
