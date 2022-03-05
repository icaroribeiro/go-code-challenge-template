package healthcheck

import "net/http"

// IHandler interface is the healthcheck's handler contract.
type IHandler interface {
	GetStatus(w http.ResponseWriter, r *http.Request)
}
