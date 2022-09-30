package user

import "net/http"

// IHandler interface is a collection of function signatures that represents the user's handler contract.
type IHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
}
