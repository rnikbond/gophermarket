package handler

import (
	"log"
	"net/http"
)

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {

	log.Println("login endpoint")
}
