package transport

import (
	"finance-tracker/internal/middleware"
	"finance-tracker/internal/service"
	"finance-tracker/internal/session"
	"net/http"
	"net/url"
)

type Handler struct {
	authSvc    *service.AuthService
	expenseSvc *service.ExpenseService
	sessMgr    *session.Manager
	middleware *middleware.AuthMiddleware
}

func NewHandler(auth *service.AuthService, exp *service.ExpenseService, sm *session.Manager, mw *middleware.AuthMiddleware) *Handler {
	return &Handler{
		authSvc:    auth,
		expenseSvc: exp,
		sessMgr:    sm,
		middleware: mw,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	// Корневой маршрут теперь просто редиректит на логин
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/logout", h.Logout)

	// Protected routes
	mux.HandleFunc("/dashboard", h.middleware.RequireAuth(h.Dashboard))
	mux.HandleFunc("/expenses/add", h.middleware.RequireAuth(h.AddExpense))
	mux.HandleFunc("/password/update", h.middleware.RequireAuth(h.UpdatePassword))

	return mux
}

// Вспомогательные функции (ОСТАВЛЯЕМ ТОЛЬКО ЗДЕСЬ)
func setFlash(w http.ResponseWriter, name, value string) {
	encodedValue := url.QueryEscape(value)
	http.SetCookie(w, &http.Cookie{Name: name, Value: encodedValue, Path: "/", MaxAge: 5, HttpOnly: true})
}

func getFlash(w http.ResponseWriter, r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	decodedValue, err := url.QueryUnescape(c.Value)
	if err != nil {
		return ""
	}
	return decodedValue
}
