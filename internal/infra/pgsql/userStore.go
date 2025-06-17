package pgsql

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type UserStore struct {
	dbConn DBConn
}

func NewUserStore(conn DBConn) *UserStore {
	return &UserStore{dbConn: conn}
}

const insertUser = `
	INSERT INTO lesta_start.users
	(id, username, password)
	VALUES
	($1, $2, $3)
	;
`

func (s *UserStore) Register(ctx context.Context, user models.User) error {
	if _, err := s.dbConn.Exec(ctx, insertUser, user.Id(), user.Name(), user.HashedPassword()); err != nil {
		return fmt.Errorf("pgsql/userStore.Register: [%w]", err)
	}

	return nil
}

const checkUserById = `
	SELECT EXISTS
	(SELECT 1 
	FROM lesta_start.users
	WHERE id = $1
	LIMIT 1)
	;
`

func (s *UserStore) IsIdRegistered(ctx context.Context, id string) (bool, error) {
	isRegisterd := false

	row := s.dbConn.QueryRow(ctx, checkUserById, id)
	if err := row.Scan(&isRegisterd); err != nil {
		return false, fmt.Errorf("pgsql/userStore.IsIdRegistered: [%w]", err)
	}

	return isRegisterd, nil
}

const checkUserByName = `
	SELECT EXISTS
	(SELECT 1 
	FROM lesta_start.users
	WHERE username = $1
	LIMIT 1)
	;
`

func (s *UserStore) IsNameRegistered(ctx context.Context, name string) (bool, error) {
	isRegisterd := false

	row := s.dbConn.QueryRow(ctx, checkUserByName, name)
	if err := row.Scan(&isRegisterd); err != nil {
		return false, fmt.Errorf("pgsql/userStore.IsNameRegistered: [%w]", err)
	}

	return isRegisterd, nil
}

const selectUserById = `
	SELECT 
	id, username, password
	FROM lesta_start.users	
	WHERE id = $1
	LIMIT 1
	;
`

func (s *UserStore) FindById(ctx context.Context, id string) (models.User, error) {
	row := s.dbConn.QueryRow(ctx, selectUserById, id)

	user, err := s.scanUser(row)
	if err != nil {
		return models.User{}, fmt.Errorf("pgsql/userStore.FindById: [%w]", err)
	}

	return user, nil
}

const selectUserByName = `
	SELECT 
	id, username, password
	FROM lesta_start.users	
	WHERE username = $1
	LIMIT 1
	;
`

func (s *UserStore) FindByName(ctx context.Context, name string) (models.User, error) {
	row := s.dbConn.QueryRow(ctx, selectUserByName, name)

	user, err := s.scanUser(row)
	if err != nil {
		return models.User{}, fmt.Errorf("pgsql/userStore.FindByName: [%w]", err)
	}

	return user, nil
}

const updateUserPassword = `
	UPDATE
	lesta_start.users
	SET password = $2
	WHERE id = $1
	;
`

func (s *UserStore) Update(ctx context.Context, user models.User) error {
	tag, err := s.dbConn.Exec(ctx, updateUserPassword, user.Id(), user.HashedPassword())
	if err != nil {
		return fmt.Errorf("pgsql/userStore.Update: [%w]", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pgsql/userStore.Update: [User No%q doesn't exist]", user.Id())
	}

	return nil
}

const deleteUserById = `
	DELETE
	FROM lesta_start.users	
	WHERE id = $1
	;
`

func (s *UserStore) Deregister(ctx context.Context, id string) error {
	tag, err := s.dbConn.Exec(ctx, deleteUserById, id)
	if err != nil {
		return fmt.Errorf("pgsql/userStore.Deregister: [%w]", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pgsql/userStore.Deregister: [User No%q doesn't exist]", id)
	}

	return nil
}

func (s *UserStore) scanUser(row pgx.Row) (models.User, error) {
	var (
		id       string
		name     string
		password string
	)
	if err := row.Scan(&id, &name, &password); err != nil {
		return models.User{}, fmt.Errorf("pgsql/userStore.scanUser: [%w]", err)
	}

	return models.NewUser(id, name, password), nil
}
