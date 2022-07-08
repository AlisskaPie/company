package middleware

import (
	"companies/ipapi"
	"companies/responses"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

var locationName = os.Getenv("LOCATION")

func LocationIP(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wantLocation := locationName
		// creation operation must be allowed only for requests received from users located in Cyprus
		loc, ok := ipapi.IsExpectedLocation(wantLocation)
		if !ok {
			err := errors.New("wrong location for this operation")
			err = errors.Wrapf(err, "got location: %v, want location: %v", loc, wantLocation)
			responses.ERROR(w, http.StatusPreconditionFailed, err)
			return
		}

		handler(w, r)
	}
}
