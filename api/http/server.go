package http

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"users-api/models"
)

type Server struct {
	Storage models.UserStorage
	Port    string
}

type ErrorResponse struct {
	Reason string `json:"reason"`
}

func (s *Server) Start() error {
	router := chi.NewRouter()

	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Route("/users", s.users)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", s.Port), router))
	return nil
}

func JSONError(w http.ResponseWriter, r *http.Request, code int, reason string) {
	render.Status(r, code)

	if reason != "" && code != http.StatusInternalServerError {
		render.JSON(w, r, ErrorResponse{Reason: reason})
		return
	}
}
