package middleware

// The original work was derived from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/pressly/lg middleware pkgs.

func (m *Middleware) Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {

				if m.log != nil {
					errMsg := fmt.Sprintf("%v", rvr)
					stackTrace := getStackTrace()
					m.log.Error(errMsg, slog.String("stack", stackTrace))
				} else {
					fmt.Println(rvr)
				}

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func getStackTrace() string {
	var buf [4096]byte
	stackSize := runtime.Stack(buf[:], false)
	stack := string(buf[:stackSize])
	return stack
}
