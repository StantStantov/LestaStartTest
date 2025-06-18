package models

import (
	"time"
)

type Session struct {
	expireOn  time.Time
	token     string
	userId    *string
	csrfToken string
}

func NewSession(token string, userId *string, csrfToken string, expireOn time.Time) Session {
	return Session{
		token:     token,
		userId:    userId,
		csrfToken: csrfToken,
		expireOn:  expireOn,
	}
}

func (s *Session) Migrate(userId *string, newSessionToken, newCsrfToken string) {
	s.userId = userId
	s.token = newSessionToken
	s.csrfToken = newCsrfToken
}

func (s Session) Token() string {
	return s.token
}

func (s Session) UserId() *string {
	return s.userId
}

func (s Session) CsrfToken() string {
	return s.csrfToken
}

func (s Session) ExpireOn() time.Time {
	return s.expireOn
}
