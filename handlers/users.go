package handlers

import (
	"golang-boiler-plate/services"
	"golang-boiler-plate/utils/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UsersHandler struct {
	Service services.UserService
}

func NewUserHandler() *UsersHandler {
	return &UsersHandler{
		Service: services.NewUserService(),
	}
}

func (h *UsersHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/{username}", h.GetUserByUsername)

	return r
}

func (h *UsersHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	res := response.New(r.Context(), w)

	username := chi.URLParam(r, "username")
	user, err := h.Service.GetByUsername(r.Context(), username)
	if err != nil {
		res.SetCode(http.StatusInternalServerError).SetMessage("Error occured").Write()
		return
	}

	res.SetCode(http.StatusOK).SetData(user).Write()
}
