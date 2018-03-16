package grmon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime/pprof"
)

func StartAgent() {
	http.HandleFunc("/debug/grmon", grmonHandler)
	go http.ListenAndServe(":1234", nil)
}

func grmonHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	pprof.Lookup("goroutine").WriteTo(&buf, 2)
	routines := ReadRoutines(buf)

	data, err := json.Marshal(routines)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
