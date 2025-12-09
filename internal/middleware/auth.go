package middleware

import (
	"context"
	"finance-tracker/internal/session"
	"net/http"
)

type AuthMiddleware struct {
	sm *session.Manager
}

func NewAuthMiddleware(sm *session.Manager) *AuthMiddleware {
	return &AuthMiddleware{sm: sm}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := m.sm.Check(r.Context(), cookie.Value)
		if err != nil || userID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Кладем userID в контекст
		ctx := context.WithValue(r.Context(), "userID", userID)
		next(w, r.WithContext(ctx))
	}
}
