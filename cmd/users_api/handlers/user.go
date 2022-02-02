package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandler struct {
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/v1/users/profile", h.Profile).Methods(http.MethodGet)
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {

}
