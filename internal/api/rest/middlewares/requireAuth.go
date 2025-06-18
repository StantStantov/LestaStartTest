package middlewares

import (
	"Stant/LestaGamesInternship/internal/app/services/sesserv"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"log"
	"net/http"
)

func RequireAuth(userStore stores.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, ok := sesserv.GetSession(r.Context())
			if !ok {
				http.Error(w, `{"error": "Couldn't get Session info"}`, http.StatusUnauthorized)
				return
			}

			userId := session.UserId()
			if userId == nil {
				http.Error(w, `{"error": "Not authenticated"}`, http.StatusUnauthorized)
				return
			}
			exist, err := userStore.IsIdRegistered(r.Context(), *userId)
			if err != nil {
				log.Printf("middlewares/requireAuth.RequireAuth: [%v]", err)
				http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
				return
			}
			if !exist {
				http.Error(w, `{"error": "Not authenticated"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
