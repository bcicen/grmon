package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/bcicen/grmon"
)

type sortFn func(grmon.Routine, grmon.Routine) bool

var (
	client = &http.Client{Timeout: 10 * time.Second}

	sortKey = "num"
	sorters = map[string]sortFn{
		"num": func(r1, r2 grmon.Routine) bool {
			return r1.Num < r2.Num
		},
		"state": func(r1, r2 grmon.Routine) bool {
			return r1.State < r2.State
		},
	}
)

type Routines []grmon.Routine

func (r Routines) Len() int           { return len(r) }
func (r Routines) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Routines) Less(i, j int) bool { return sorters[sortKey](r[i], r[j]) }

func poll() (routines Routines) {
	url := fmt.Sprintf("http://%s%s", *hostFlag, *endpointFlag)
	r, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&routines)
	if err != nil {
		panic(err)
	}

	sort.Sort(routines)
	return
}
