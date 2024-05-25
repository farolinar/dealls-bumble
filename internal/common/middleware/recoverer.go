package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/farolinar/dealls-bumble/internal/common/response"
	servicebase "github.com/farolinar/dealls-bumble/services/base"
	"github.com/rs/zerolog/log"
)

func PanicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				if r != http.ErrAbortHandler {
					log.Error().Msg(fmt.Sprintf("Recovered from panic: %s", string(debug.Stack())))
				}
				response.JSON(w, http.StatusInternalServerError, servicebase.ResponseBody{
					Message: "Internal server error",
					Code:    servicebase.Code5XX,
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
