package grmon

import (
	"net/http"
	_ "net/http/pprof"
)

func Start() { go http.ListenAndServe(":1234", nil) }
