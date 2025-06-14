package pgsql

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"Stant/LestaGamesInternship/internal/domain/stores"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type CollectionStore struct {
	dbConn        DBConn
	documentStore stores.DocumentStore
}

func NewCollectionStore(dbConn DBConn, documentStore stores.DocumentStore) *CollectionStore {
	return &CollectionStore{dbConn: dbConn, documentStore: documentStore}
}

const insertCollection = `
	INSERT INTO lesta_start.collections
	(id, user_id, collection)
	VALUES
	($1, $2, $3)
	;
`

const pinDocumentToCollection = `
	INSERT INTO lesta_start.collection_documents
	(collection_id, document_id)
	VALUES
	($1, $2)
	;
`

func (s *CollectionStore) Save(ctx context.Context, collection models.Collection) error {
	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pgsql/collectionStore.Save: [%w]", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, insertCollection, collection.Id(), collection.UserId(), collection.Name()); err != nil {
		return fmt.Errorf("pgsql/collectionStore.Save: [%w]", err)
	}
	documents := collection.Documents()
	for _, document := range documents {
		if _, err := tx.Exec(ctx, pinDocumentToCollection, collection.Id(), document.Id()); err != nil {
			return fmt.Errorf("pgsql/collectionStore.Save: [%w]", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("pgsql/collectionStore.Save: [%w]", err)
	}

	return nil
}

func (s *CollectionStore) PinDocument(ctx context.Context, collectionId, documentId string) error {
	tag, err := s.dbConn.Exec(ctx, pinDocumentToCollection, collectionId, documentId)
	if err != nil {
		return fmt.Errorf("pgsql/collectionStore.PinDocument [%w]", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pgsql/collectionStore.PinDocument [%w]", err)
	}

	return nil
}

const checkCollectionById = `
	SELECT EXISTS
	(SELECT 1 FROM lesta_start.collections
	WHERE id = $1
	LIMIT 1)
	;
`

func (s *CollectionStore) IsExist(ctx context.Context, id string) (bool, error) {
	isExist := false

	row := s.dbConn.QueryRow(ctx, checkCollectionById, id)
	if err := row.Scan(&isExist); err != nil {
		return false, fmt.Errorf("pgsql/collectionStore.IsExist: [%w]", err)
	}

	return isExist, nil
}

const checkIsDocPinnedById = `
	SELECT EXISTS
	(SELECT 1 FROM lesta_start.collection_documents
	WHERE collection_id = $1 AND document_id = $2
	LIMIT 1)
	;
`

func (s *CollectionStore) IsPinned(ctx context.Context, collectionId, documentId string) (bool, error) {
	isExist := false

	row := s.dbConn.QueryRow(ctx, checkIsDocPinnedById, collectionId, documentId)
	if err := row.Scan(&isExist); err != nil {
		return false, fmt.Errorf("pgsql/collectionStore.IsPinned: [%w]", err)
	}

	return isExist, nil
}

const selectCollectionById = `
	SELECT c.id, c.user_id, c.collection, d.id
	FROM lesta_start.collection_documents AS cd
	INNER JOIN lesta_start.collections AS c ON c.id = cd.collection_id
	INNER JOIN lesta_start.documents AS d ON d.id = cd.document_id
	WHERE cd.collection_id = $1
	;
`

func (s *CollectionStore) FindById(ctx context.Context, id string) (*models.Collection, error) {
	rows, err := s.dbConn.Query(ctx, selectCollectionById, id)
	if err != nil {
		return nil, fmt.Errorf("pgsql/collectionStore.FindById: [%w]", err)
	}

	collection, err := s.scanCollection(ctx, rows)
	if err != nil {
		return nil, fmt.Errorf("pgsql/collectionStore.FindById: [%w]", err)
	}

	return collection, nil
}

const selectCollectionsByUserId = `
	SELECT c.id, c.user_id, c.collection, d.id
	FROM lesta_start.collection_documents AS cd
	INNER JOIN lesta_start.collections AS c ON c.id = cd.collection_id
	INNER JOIN lesta_start.documents AS d ON d.id = cd.document_id
	WHERE c.user_id = $1
	ORDER BY c.id
	;
`

func (s *CollectionStore) FindByUserId(ctx context.Context, userId string) ([]*models.Collection, error) {
	rows, err := s.dbConn.Query(ctx, selectCollectionsByUserId, userId)
	if err != nil {
		return nil, fmt.Errorf("pgsql/collectionStore.FindByUserId: [%w]", err)
	}

	collections, err := s.scanCollections(ctx, rows)
	if err != nil {
		return nil, fmt.Errorf("pgsql/collectionStore.FindByUserId: [%w]", err)
	}

	return collections, nil
}

const updateCollectionById = `
	UPDATE
	lesta_start.collections
	SET collection = $2 
	WHERE id = $1
	;
`

func (s *CollectionStore) Rename(ctx context.Context, id, newName string) error {
	tag, err := s.dbConn.Exec(ctx, updateCollectionById, id, newName)
	if err != nil {
		return fmt.Errorf("pgsql/collectionStore.Rename [%w]", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pgsql/collectionStore.Rename [%w]", err)
	}

	return nil
}

const unpinDocumentFromCollection = `
	DELETE 
	FROM lesta_start.collection_documents
	WHERE collection_id = $1 AND document_id = $2
	;
`

func (s *CollectionStore) UnpinDocument(ctx context.Context, collectionId, documentId string) error {
	tag, err := s.dbConn.Exec(ctx, unpinDocumentFromCollection, collectionId, documentId)
	if err != nil {
		return fmt.Errorf("pgsql/collectionStore.UnpinDocument [%w]", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pgsql/collectionStore.UnpinDocument [%w]", err)
	}

	return nil
}

const deleteCollectionById = `
	DELETE 
	FROM lesta_start.collections
	WHERE id = $1
	;
`

func (s *CollectionStore) Delete(ctx context.Context, id string) error {
	tag, err := s.dbConn.Exec(ctx, deleteCollectionById, id)
	if err != nil {
		return fmt.Errorf("pgsql/collectionStore.Delete [%w]", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pgsql/collectionStore.Delete [Collection No%q doesn't exist]", id)
	}

	return nil
}

func (s *CollectionStore) scanCollection(ctx context.Context, rows pgx.Rows) (*models.Collection, error) {
	var (
		id         string
		userId     string
		name       string
		documentId string
	)
	documents := []models.Document{}
	for rows.Next() {
		if err := rows.Scan(&id, &userId, &name, &documentId); err != nil {
			return nil, fmt.Errorf("pgsql/collectionStore.scanCollection: [%w]", err)
		}
		document, err := s.documentStore.Open(ctx, documentId)
		if err != nil {
			return nil, fmt.Errorf("pgsql/collectionStore.scanCollection: [%w]", err)
		}

		documents = append(documents, document)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("pgsql/collectionStore.scanCollection: [%w]", err)
	}

	collection := models.NewCollection(id, userId, name, documents)
	if rows.CommandTag().RowsAffected() == 0 {
		return collection, fmt.Errorf("pgsql/collectionStore.scanCollection: [%w]", pgx.ErrNoRows)
	}
	return collection, nil
}

func (s *CollectionStore) scanCollections(ctx context.Context, rows pgx.Rows) ([]*models.Collection, error) {
	var (
		id         string
		userId     string
		name       string
		documentId string
	)
	documents := []models.Document{}
	collections := []*models.Collection{}
	prevId := ""
	prevName := ""
	for rows.Next() {
		if err := rows.Scan(&id, &userId, &name, &documentId); err != nil {
			return nil, fmt.Errorf("pgsql/collectionStore.scanCollections: [%w]", err)
		}
		document, err := s.documentStore.Open(ctx, documentId)
		if err != nil {
			return nil, fmt.Errorf("pgsql/collectionStore.scanCollections: [%w]", err)
		}

		if prevId != "" && prevId != id {
			collections = append(collections, models.NewCollection(prevId, userId, prevName, documents))
			documents = []models.Document{}
		}

		documents = append(documents, document)
		prevId = id
		prevName = name
	}
	if prevId != "" {
		collections = append(collections, models.NewCollection(prevId, userId, prevName, documents))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("pgsql/collectionStore.scanCollections: [%w]", err)
	}

	if rows.CommandTag().RowsAffected() == 0 {
		return collections, fmt.Errorf("pgsql/collectionStore.scanCollections: [%w]", pgx.ErrNoRows)
	}
	return collections, nil
}
