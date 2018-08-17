package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/covveco/micro-learning/handler/auth"
	"github.com/covveco/micro-learning/model"
	"golang.org/x/crypto/bcrypt"
)

// Login handles the login request
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	staffNo := r.FormValue("staff_no")
	password := r.FormValue("password")

	u, err := h.DB.GetUserByStaffNo(staffNo)
	if err != nil {
		redirect(w, r, "/login", "Login failed. Username or password invalid")
		return
	}

	if !u.Registered {
		redirect(w, r, "/login", "User is not registered. Check your email or with your administrator")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
	if err != nil {
		redirect(w, r, "/login", "Login failed. Username or password invalid")
		return
	}

	token, err := auth.GenerateToken(u.ID, h.AuthSecret)
	if err != nil {
		redirect(w, r, "/login", "Token generation failed")
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: token})

	redirect(w, r, "/", "")
}

// LoginForm renders the login form
func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	f, _ := GetFlash(w, r)
	renderView(w, "login", struct{ Flash string }{f})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: ""})
	redirect(w, r, "/", "")
}

// Registration handles the registration request
func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	password := r.FormValue("password")
	cPassword := r.FormValue("confirm_password")

	if password != cPassword {
		redirect(w, r, "/register", "Registration failed. Passwords do not match.")
		return
	}

	u, err := h.DB.GetUserByToken(token)
	if err != nil {
		redirect(w, r, "/register", "Registration failed. Token is invalid.")
		return
	}

	if u.Registered {
		redirect(w, r, "/login", "You are already registered")
		return
	}
	pHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		redirect(w, r, "/register", "Registration failed. Password invalid.")
		return
	}
	if err = h.DB.SetUserPassword(u, string(pHash)); err != nil {
		redirect(w, r, "/register", "Registration Failed. Password invalid.")
		return
	}
	redirect(w, r, "/login", "")

}

// RegistrationForm renders the registration form
func (h *Handler) RegistrationForm(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	_, err := h.DB.GetUserByToken(token)
	if err != nil {
		h.NotFoundError(w, r)
		return
	}

	params := struct {
		Flash string
		Token string
	}{
		Token: token,
	}

	f, err := GetFlash(w, r)
	if err == nil {
		params.Flash = f
	}
	renderView(w, "register", params)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	u := &model.User{}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := h.DB.CreateUser(u); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	uR, err := h.DB.GetUserByStaffNo(u.StaffNo)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(uR.ID)))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	u := &model.User{}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := h.DB.DeleteUserByID(u.ID); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(u.ID)))
}
