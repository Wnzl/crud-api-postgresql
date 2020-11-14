package http

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"users-api/models"
)

func (s *Server) users(r chi.Router) {
	r.Get("/", s.getAllUsers)
	r.Post("/", s.addUser)
	r.Route("/{userID}", func(r chi.Router) {
		r.Use(s.UserCtx)
		r.Get("/", s.getUser)
		r.Put("/", s.updateUser)
		r.Delete("/", s.deleteUser)
	})
}

func (s *Server) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.User

		if userID := chi.URLParam(r, "userID"); userID != "" {
			ID, err := strconv.Atoi(userID)
			if err != nil {
				JSONError(w, r, http.StatusBadRequest, "user ID can't be parsed")
				return
			}

			user, err = s.Storage.Get(ID)
			if err != nil {
				JSONError(w, r, http.StatusBadRequest, "user not found")
				return
			}
		} else {
			JSONError(w, r, http.StatusBadRequest, "user ID is empty")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	render.JSON(w, r, user)
}

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.Storage.GetAll()
	if err != nil {
		JSONError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, r, users)
}

func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := render.Bind(r, user); err != nil {
		JSONError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	exist, err := s.Storage.UserExist(user)
	if err != nil {
		logrus.WithError(err).Fatal("checking if user exist")
		JSONError(w, r, http.StatusInternalServerError, "error checking if user exist")
		return
	}

	if exist {
		JSONError(w, r, http.StatusBadRequest, "user already exist")
		return
	}

	_, err = s.Storage.Store(user)
	if err != nil {
		JSONError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, r, user)
}

func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	if err := render.Bind(r, user); err != nil {
		JSONError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := s.Storage.Update(user.ID, user)
	if err != nil {
		JSONError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, r, user)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	err := s.Storage.Delete(user.ID)
	if err != nil {
		JSONError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(w, r, "successfully deleted")
}
