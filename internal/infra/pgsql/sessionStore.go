package pgsql

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type SessionStore struct {
	dbConn DBConn
}

func NewSessionStore(dbConn DBConn) *SessionStore {
	return &SessionStore{dbConn: dbConn}
}

const insertSession = `
  INSERT INTO lesta_start.sessions
  (token, user_id, csrf_token, expire_on)
  VALUES
  ($1, $2, $3, $4)
	ON CONFLICT (token) DO UPDATE
	SET user_id = $2, csrf_token = $3, expire_on = $4
  ;
`

func (s *SessionStore) Create(ctx context.Context, session models.Session) error {
	if _, err := s.dbConn.Exec(ctx, insertSession, session.Token(), session.UserId(), session.CsrfToken(), session.ExpireOn()); err != nil {
		return fmt.Errorf("pgsql/sessionStore.Create: [%w]", err)
	}

	return nil
}

const selectSession = `
  SELECT token, user_id, csrf_token, expire_on
	FROM lesta_start.sessions
  WHERE token = $1
	AND current_timestamp < expire_on
  LIMIT 1
  ;
`

func (s *SessionStore) Find(ctx context.Context, token string) (models.Session, error) {
	row := s.dbConn.QueryRow(ctx, selectSession, token)

	session, err := s.scanSession(row)
	if err != nil {
		return models.Session{}, fmt.Errorf("pgsql/sessionStore.Find: [%w]", err)
	}

	return session, nil
}

const deleteSessionByToken = `
  DELETE FROM lesta_start.sessions
  WHERE token = $1 
  ;
`

func (s *SessionStore) Delete(ctx context.Context, token string) error {
	if _, err := s.dbConn.Exec(ctx, deleteSessionByToken, token); err != nil {
		return fmt.Errorf("pgsql/sessionStore.Delete: [%w]", err)
	}
	return nil
}

const deleteExpiredSessions = `
  DELETE FROM lesta_start.sessions
  WHERE expire_on < current_timestamp 
  ;
`

func (s *SessionStore) DeleteAllExpired() error {
	if _, err := s.dbConn.Exec(context.Background(), deleteExpiredSessions); err != nil {
		return fmt.Errorf("pgsql/sessionStore.DeleteAllExpired: [%w]", err)
	}
	return nil
}

func (s *SessionStore) scanSession(row pgx.Row) (models.Session, error) {
	var (
		token     string
		userId    *string
		csrfToken string
		expireOn  time.Time
	)
	if err := row.Scan(&token, &userId, &csrfToken, &expireOn); err != nil {
		return models.Session{}, fmt.Errorf("pgsql/sessionStore.scanSession: [%w]", err)
	}

	return models.NewSession(token, userId, csrfToken, expireOn), nil
}

func (s *SessionStore) StartCleanup(
	interval time.Duration,
) (chan<- struct{}, <-chan struct{}) {
	quit, done := make(chan struct{}), make(chan struct{})
	go s.clean(interval, quit, done)

	return quit, done
}

func (s *SessionStore) StopCleanup(quit chan<- struct{}, done <-chan struct{}) {
	quit <- struct{}{}
	<-done
}

func (s *SessionStore) clean(
	interval time.Duration,
	quit <-chan struct{},
	done chan<- struct{},
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			done <- struct{}{}
			return
		case <-ticker.C:
			if err := s.DeleteAllExpired(); err != nil {
				log.Printf("pgsql/sessionStore.clean: [%v]", err)
			}
		}
	}
}
