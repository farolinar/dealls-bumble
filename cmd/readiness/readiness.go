package readiness

import (
	"net/http"

	"github.com/farolinar/dealls-bumble/config/postgres"
)

func DBReadinessHandler(rw http.ResponseWriter, r *http.Request) {
	db := postgres.GetDBConnection()
	err := db.Ping()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(200)
}
