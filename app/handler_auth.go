package app

import (
	"fmt"
	"net/http"

	"github.com/jritsema/gotoolbox/web"
)

// GET /auth/login
func (a *App) LoginForm(w http.ResponseWriter, r *http.Request) {
	res := web.HTML(http.StatusOK, a.Templ, "auth-login.html", nil, nil)
	res.Write(w)
}

// GET /auth/register
func (a *App) RegisterForm(w http.ResponseWriter, r *http.Request) {
	res := web.HTML(http.StatusOK, a.Templ, "auth-register.html", nil, nil)
	res.Write(w)
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println(email, password)
	w.Write([]byte("Login"))
}

func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register"))
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout"))
}
