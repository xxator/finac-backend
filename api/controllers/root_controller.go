package controllers

import (
	"net/http"

	"github.com/xxator/finac-backend/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To Finac  API")

}