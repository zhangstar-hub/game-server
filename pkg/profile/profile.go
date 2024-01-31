package profile

import (
	"net/http"
	_ "net/http/pprof"
)

func StartProfile() {
	http.ListenAndServe("0.0.0.0:6060", nil)
}
