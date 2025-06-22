package handlers

import (
	"Stant/LestaGamesInternship/internal/app/services/sesserv"
	"Stant/LestaGamesInternship/internal/app/services/usrserv"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"encoding/json"
	"log"
	"net/http"
)

// @Summary Регистрация пользователя
// @Description Зарегестрировать нового пользователя.
// @Tags Аунтефикация
// @Produce json
// @Success 200 {object} dto.SuccessMessage
// @Router /api/register [post]
func HandlePostRegister(userService *usrserv.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Println(err)
			return
		}

		username := body.Username
		password := body.Password
		if !userService.IsValidUsername(username) ||
			!userService.IsValidPassword(password) {
			return
		}
		if err := userService.Register(r.Context(), username, password); err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// @Summary Вход в аккаунт
// @Description Начинает новую сессию и выдает куки.
// @Tags Аунтефикация
// @Produce json
// @Success 200 {object} dto.SuccessMessage
// @Router /api/login [post]
func HandlePostLogin(
	userService *usrserv.UserService,
	userStore stores.UserStore,
	sessionService *sesserv.SessionService,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		username := body.Username
		password := body.Password
		if !userService.IsValidUsername(username) ||
			!userService.IsValidPassword(password) {
			return
		}
		registered, err := userService.IsRegistered(r.Context(), username, password)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if !registered {
			log.Println("incorrect username/password")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// TODO simplify
		user, err := userStore.FindByName(r.Context(), username)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		userId := user.Id()

		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if err := sessionService.Migrate(r.Context(), &userId, &session); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if err := sessionService.Save(r.Context(), session); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		sessionService.SetSessionCookie(w, session)
		sessionService.SetCsrfCookie(w, session)
		w.WriteHeader(http.StatusOK)
	})
}

// @Summary Выход из аккаунта
// @Description Окончивает сессию пользователя и удаляет куки.
// @Tags Аунтефикация
// @Produce json
// @Success 200 {object} dto.SuccessMessage
// @Router /api/logout [get]
func HandlePostLogout(sessionService *sesserv.SessionService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, "Failed to get Session info", http.StatusInternalServerError)
			return
		}

		if err := sessionService.Stop(r.Context(), &session); err != nil {
			log.Println(err)
			http.Error(w, "Failed to stop Session", http.StatusInternalServerError)
			return
		}

		sessionService.ExpireSessionCookie(w, session)
		sessionService.ExpireCsrfCookie(w, session)
		w.WriteHeader(http.StatusOK)
	})
}

// @Summary Смена пароля
// @Description Меняет пароль пользователя.
// @Tags Аунтефикация
// @Produce json
// @Success 200 {object} dto.SuccessMessage
// @Router /api/user/{user_id} [patch]
func HandlePatchUser(pathValue string, userService *usrserv.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUserId := r.PathValue(pathValue)
		if requestUserId == "" {
			http.Error(w, `{"error" : "Empty path value"}`, http.StatusBadRequest)
			return
		}

		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, `{"error" : "Session does not exist"}`, http.StatusBadRequest)
			return
		}
		sessionUserId := session.UserId()
		if sessionUserId == nil || *sessionUserId != requestUserId {
			http.Error(w, `{"error": "User ID mismatch"}`, http.StatusUnauthorized)
			return
		}

		requestBody := struct {
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			log.Printf("handlers/user.HandlePatchUser: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		if err := userService.ChangePassword(r.Context(),
			requestUserId,
			requestBody.OldPassword,
			requestBody.NewPassword); err != nil {
			log.Printf("handlers/user.HandlePatchUser: [%v]", err)
			http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// @Summary Регистрация пользователя
// @Description Зарегестрировать нового пользователя.
// @Tags Аунтефикация
// @Produce json
// @Success 200 {object} dto.SuccessMessage
// @Router /api/user/{user_id} [delete]
func HandleDeleteUser(
	pathValue string,
	userService *usrserv.UserService,
	sessionService *sesserv.SessionService,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUserId := r.PathValue(pathValue)
		if requestUserId == "" {
			http.Error(w, `{"error" : "Empty path value"}`, http.StatusBadRequest)
			return
		}

		session, ok := sesserv.GetSession(r.Context())
		if !ok {
			http.Error(w, `{"error" : "Something went wrong"}`, http.StatusInternalServerError)
			return
		}
		sessionUserId := session.UserId()
		if sessionUserId == nil || *sessionUserId != requestUserId {
			http.Error(w, `{"error" : "User ID mismatch"}`, http.StatusUnauthorized)
			return
		}

		if err := userService.Deregister(r.Context(), requestUserId); err != nil {
			log.Printf("handlers/user.HandlePatchUser: [%v]", err)
			http.Error(w, `{"error" : "Something went wrong"}`, http.StatusInternalServerError)
			return
		}

		sessionService.ExpireSessionCookie(w, session)
		sessionService.ExpireCsrfCookie(w, session)
		w.WriteHeader(http.StatusOK)
	})
}
