package middlewares

import (
	"Stant/LestaGamesInternship/internal/app/services/sesserv"
	"Stant/LestaGamesInternship/internal/domain/models"
	"net/http"
)

func RequireCsrf(sessionService *sesserv.SessionService) func(http.Handler) http.Handler {
	csrfHeaderName := sessionService.CsrfHeaderName()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, ok := sesserv.GetSession(r.Context())
			if !ok {
				http.Error(w, `{"error": "Couldn't get Session info"}`, http.StatusUnauthorized)
				return
			}

			if r.Method == http.MethodPost || r.Method == http.MethodPut ||
				r.Method == http.MethodDelete || r.Method == http.MethodPatch {
				if !isValidCsrfToken(r, csrfHeaderName, session) {
					http.Error(w, `{"error": "Couldn't verify CSRF token"}`, http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isValidCsrfToken(r *http.Request, csrfHeaderName string, session models.Session) bool {
	sessionToken := session.CsrfToken()

	requestToken := r.Header.Get(csrfHeaderName)
	if requestToken == "" {
		return false
	}

	return sessionToken == requestToken
}
