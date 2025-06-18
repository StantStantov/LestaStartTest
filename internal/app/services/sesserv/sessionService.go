package sesserv

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

const bytesEntropy = 32

type sessionContextKey struct{}

type SessionService struct {
	sessionStore stores.SessionStore

	sessionLifetime   time.Duration
	domain            string
	sessionCookieName string
	csrfCookieName    string
	csrfHeaderName    string
}

func NewSessionService(sessionStore stores.SessionStore,
	domain string,
	sessionCookieName string,
	sessionLifetime time.Duration,
	csrfCookieName string,
	csrfHeaderName string,
) *SessionService {
	return &SessionService{
		sessionStore:      sessionStore,
		domain:            domain,
		sessionLifetime:   sessionLifetime,
		sessionCookieName: sessionCookieName,
		csrfCookieName:    csrfCookieName,
		csrfHeaderName:    csrfHeaderName,
	}
}

func (s *SessionService) Start(r *http.Request) (models.Session, error) {
	var session models.Session

	cookie, err := r.Cookie(s.sessionCookieName)
	if err == nil {
		token := cookie.Value
		session, err = s.sessionStore.Find(r.Context(), token)
		if err != nil {
			return models.Session{}, fmt.Errorf("[%w]", err)
		}
	}

	if session == (models.Session{}) || s.IsExpired(session) {
		session = s.newSession(nil)
	}

	return session, nil
}

func (s *SessionService) Save(ctx context.Context, session models.Session) error {
	if err := s.sessionStore.Create(ctx, session); err != nil {
		return fmt.Errorf("sessions/sessionService.Save: [%w]", err)
	}

	return nil
}

func (s *SessionService) IsExpired(session models.Session) bool {
	return time.Now().After(session.ExpireOn())
}

func (s *SessionService) Migrate(ctx context.Context, newUserId *string, session *models.Session) error {
	if err := s.sessionStore.Delete(ctx, session.Token()); err != nil {
		return fmt.Errorf("[%w]", err)
	}

	session.Migrate(
		newUserId,
		generateToken(bytesEntropy),
		generateToken(bytesEntropy))

	return nil
}

func (s *SessionService) Stop(ctx context.Context, session *models.Session) error {
	err := s.sessionStore.Delete(ctx, session.Token())
	if err != nil {
		return err
	}

	session.Migrate(
		nil,
		generateToken(bytesEntropy),
		generateToken(bytesEntropy))

	return nil
}

func (s *SessionService) SessionLifetime() time.Duration {
	return s.sessionLifetime
}

func (s *SessionService) SessionCookieName() string {
	return s.sessionCookieName
}

func (s *SessionService) CsrfCookieName() string {
	return s.csrfCookieName
}

func (s *SessionService) CsrfHeaderName() string {
	return s.csrfHeaderName
}

func (s *SessionService) newSession(userId *string) models.Session {
	return models.NewSession(
		generateToken(bytesEntropy),
		userId,
		generateToken(bytesEntropy),
		time.Now().Add(s.sessionLifetime),
	)
}

func (s *SessionService) SetSessionCookie(w http.ResponseWriter, session models.Session) {
	sessionCookie := newSecureCookie(s.sessionCookieName, session.Token(), s.domain)
	sessionCookie.Expires = time.Now().Add(s.sessionLifetime)
	sessionCookie.MaxAge = int(s.sessionLifetime / time.Second)

	http.SetCookie(w, sessionCookie)
}

func (s *SessionService) ExpireSessionCookie(w http.ResponseWriter, session models.Session) {
	sessionCookie := newExpiredSecureCookie(s.sessionCookieName, session.Token(), s.domain)

	http.SetCookie(w, sessionCookie)
}

func (s *SessionService) SetCsrfCookie(w http.ResponseWriter, session models.Session) {
	csrfCookie := newSecureCookie(s.csrfCookieName, session.CsrfToken(), s.domain)

	http.SetCookie(w, csrfCookie)
}

func (s *SessionService) ExpireCsrfCookie(w http.ResponseWriter, session models.Session) {
	csrfCookie := newExpiredSecureCookie(s.csrfCookieName, session.CsrfToken(), s.domain)

	http.SetCookie(w, csrfCookie)
}

func SetSession(ctx context.Context, session models.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey{}, session)
}

func GetSession(ctx context.Context) (models.Session, bool) {
	session, ok := ctx.Value(sessionContextKey{}).(models.Session)
	return session, ok
}

func newSecureCookie(name, value, domain string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   domain,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func newExpiredSecureCookie(name, value, domain string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   domain,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now(),
		MaxAge:   -1,
	}
}

func generateToken(bytesEntropy uint) string {
	token := make([]byte, bytesEntropy)
	rand.Read(token)

	return base64.URLEncoding.EncodeToString(token)
}
