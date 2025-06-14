package pgsql

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type DocumentStore struct {
	dbConn    DBConn
	fileStore stores.FileStore
}

func NewDocumentStore(dbConn DBConn, fileStore stores.FileStore) *DocumentStore {
	return &DocumentStore{
		dbConn:    dbConn,
		fileStore: fileStore,
	}
}

const createDocument = `
	INSERT INTO lesta_start.documents	
	(id, user_id, document)
	VALUES
	($1, $2, $3)
	;
`

func (s *DocumentStore) Save(ctx context.Context, document models.Document) error {
	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pgsql/documentStore.Save: [%w]", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, createDocument, document.Id(), document.UserId(), document.Name()); err != nil {
		return fmt.Errorf("pgsql/documentStore.Save: [%w]", err)
	}
	filename := s.formatDocumentName(document.Id(), document.UserId(), document.Name())
	if err := s.fileStore.Save(filename, document.File()); err != nil {
		return fmt.Errorf("pgsql/documentStore.Save: [%w]", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("pgsql/documentStore.Save: [%w]", err)
	}

	return nil
}

const selectDocumentByName = `
	SELECT id, user_id, document
	FROM lesta_start.documents
	WHERE user_id = $1 AND document = $2
	LIMIT 1
	;
`

func (s *DocumentStore) IsNameExist(ctx context.Context, userId, name string) (bool, error) {
	row := s.dbConn.QueryRow(ctx, selectDocumentByName, userId, name)

	var id string
	if err := row.Scan(&id, &userId, &name); err != nil {
		return false, fmt.Errorf("pgsql/documentStore.IsNameExist: [%w]", err)
	}
	isExist, err := s.fileStore.IsExist(s.formatDocumentName(id, userId, name))
	if err != nil {
		return false, fmt.Errorf("pgsql/documentStore.IsNameExist: [%w]", err)
	}

	return isExist, nil
}

const selectDocumentById = `
	SELECT id, user_id, document
	FROM lesta_start.documents
	WHERE id = $1
	LIMIT 1
	;
`

func (s *DocumentStore) IsIdExist(ctx context.Context, id string) (bool, error) {
	row := s.dbConn.QueryRow(ctx, selectDocumentById, id)

	var userId string
	var name string
	if err := row.Scan(&id, &userId, &name); err != nil {
		return false, fmt.Errorf("pgsql/documentStore.IsNameExist: [%w]", err)
	}
	isExist, err := s.fileStore.IsExist(s.formatDocumentName(id, userId, name))
	if err != nil {
		return false, fmt.Errorf("pgsql/documentStore.IsNameExist: [%w]", err)
	}

	return isExist, nil
}

func (s *DocumentStore) Open(ctx context.Context, id string) (models.Document, error) {
	row := s.dbConn.QueryRow(ctx, selectDocumentById, id)

	document, err := s.scanDocument(row)
	if err != nil {
		return document, fmt.Errorf("pgsql/documentStore.Open: [%w]", err)
	}

	return document, nil
}

const updateDocumentById = `
	UPDATE
	lesta_start.documents new
	SET document = $2
	FROM (SELECT 
		id, user_id, document 
		FROM lesta_start.documents 
		WHERE id = $1 
		FOR UPDATE) old
	WHERE new.id = old.id
	RETURNING old.user_id, old.document
	;
`

func (s *DocumentStore) Rename(ctx context.Context, id, newName string) error {
	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pgsql/documentStore.Rename: [%w]", err)
	}
	defer tx.Rollback(ctx)

	var (
		userId  string
		oldName string
	)
	row := tx.QueryRow(ctx, updateDocumentById, id, newName)
	if err := row.Scan(&userId, &oldName); err != nil {
		return fmt.Errorf("pgsql/documentStore.Rename: [%w]", err)
	}
	oldFilename := s.formatDocumentName(id, userId, oldName)
	newFilename := s.formatDocumentName(id, userId, newName)
	if err := s.fileStore.Rename(oldFilename, newFilename); err != nil {
		return fmt.Errorf("pgsql/documentStore.Rename: [%w]", err)
	}

	return nil
}

const deleteDocumentById = `
	DELETE 
	FROM lesta_start.documents
	WHERE id = $1
	RETURNING id, user_id, document
	;
`

func (s *DocumentStore) Delete(ctx context.Context, id string) error {
	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pgsql/documentStore.Delete: [%w]", err)
	}
	defer tx.Rollback(ctx)

	var (
		userId string
		name   string
	)
	row := tx.QueryRow(ctx, deleteDocumentById, id)
	if err := row.Scan(&id, &userId, &name); err != nil {
		return fmt.Errorf("pgsql/documentStore.Delete: [%w]", err)
	}
	filename := s.formatDocumentName(id, userId, name)
	if err := s.fileStore.Delete(filename); err != nil {
		return fmt.Errorf("pgsql/documentStore.Delete: [%w]", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("pgsql/documentStore.Delete: [%w]", err)
	}

	return nil
}

func (s *DocumentStore) scanDocument(row pgx.Row) (models.Document, error) {
	var (
		id     string
		userId string
		name   string
		file   *os.File
	)
	if err := row.Scan(&id, &userId, &name); err != nil {
		return models.Document{}, fmt.Errorf("pgsql/documentStore.scanDocument: [%w]", err)
	}
	filename := s.formatDocumentName(id, userId, name)
	file, err := s.fileStore.Open(filename)
	if err != nil {
		return models.Document{}, fmt.Errorf("pgsql/documentStore.scanDocument: [%w]", err)
	}

	return models.NewDocument(id, userId, name, file), nil
}

func (s *DocumentStore) formatDocumentName(id, userId, documentName string) string {
	return fmt.Sprintf("%s_%s_%s", id, userId, documentName)
}
