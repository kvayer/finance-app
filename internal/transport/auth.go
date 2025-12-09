package transport

import (
	"finance-tracker/internal/models"
	"html/template"
	"net/http"
	"time"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Если уже залогинен - в дашборд
	if _, err := r.Cookie("session_token"); err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("ui/html/base.html", "ui/html/register.html"))
		// Передаем IsLoggedIn: false
		tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
			"IsLoggedIn": false,
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	req := models.RegisterRequest{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}

	if err := h.authSvc.Register(r.Context(), req); err != nil {
		// Показываем ошибку на странице (просто текстом для надежности)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Если уже залогинен - в дашборд
	if _, err := r.Cookie("session_token"); err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("ui/html/base.html", "ui/html/login.html"))
		tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
			"IsLoggedIn": false,
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	req := models.LoginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	token, err := h.sessMgr.Create(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
