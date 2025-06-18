package middlewares

import (
	"Stant/LestaGamesInternship/internal/app/services/sesserv"
	"log"
	"net/http"
)

func CheckSession(sessionService *sesserv.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sessionService.Start(r)
			if err != nil {
				log.Printf("middlewares/checkSession.CheckSession: [%v]", err)
				http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
				return
			}
			ctx := sesserv.SetSession(r.Context(), session)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
