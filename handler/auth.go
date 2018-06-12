package handler

import (
	"context"
	"net/http"

	"gitlab.com/covveco/special-needs/handler/auth"
)

func (h *Handler) AuthenticateUser(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			redirect(w, r, "/login", "You must login first")
			return
		}

		id, err := auth.ValidateToken(cookie.Value, h.AuthSecret)
		if err != nil {
			redirect(w, r, "/login", "You must login first")
			return
		}

		u, err := h.DB.GetUserByID(id)
		if err != nil {
			redirect(w, r, "/login", "You must login first")
			return
		}

		ctx := context.WithValue(r.Context(), "user", u)
		n.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) AuthenticateAdmin(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pw := r.Header.Get("X-Admin-Password")
		if h.AdminPassword != pw {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		n.ServeHTTP(w, r)
	})
}
