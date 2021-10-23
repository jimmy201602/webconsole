package website

import (
	"net/http"
	"runtime/pprof"
)

func init() {
	Add_HandleFunc("get", "/pprof", Pprof_handler)
}

func Pprof_handler(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	user, pwd, ok := r.BasicAuth()
	if ok {
		if "admin" == user && "password" == pwd {
			ctx.SetContentType("text/plain")
			// p := pprof.Lookup("goroutine")
			p := pprof.Lookup("heap")
			p.WriteTo(w, 1)
			return
		} else {
			ctx.BasicAuth("Please enter the password authentication information")
		}
	} else {
		ctx.BasicAuth("Please enter the password authentication information")
	}
}
