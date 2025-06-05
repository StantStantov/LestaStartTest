package models

import (
	"fmt"
	"maps"
	"slices"
)

type Collection struct {
	documents map[string]Document
	name      string
	userUid   string
	id        uint64
}

func NewCollection(id uint64, name, userUid string, documents []Document) *Collection {
	documentsMap := make(map[string]Document, len(documents))
	for _, document := range documents {
		documentsMap[document.name] = document
	}

	return &Collection{
		documents: documentsMap,
		name:      name,
		userUid:   userUid,
		id:        id,
	}
}

func NewEmptyCollection(id uint64, name, userUid string) *Collection {
	return &Collection{
		documents: make(map[string]Document, 0),
		name:      name,
		userUid:   userUid,
		id:        id,
	}
}

func (c *Collection) AddDocument(document Document) error {
	if _, present := c.documents[document.name]; present {
		return fmt.Errorf("models/collection.AddDocument: [Collection No%d already contains Document %q]", c.id, document.name)
	}
	c.documents[document.name] = document

	return nil
}

func (c *Collection) FindDocument(name string) (Document, error) {
	document, present := c.documents[name]
	if !present {
		return document, fmt.Errorf("models/collection.AddDocument: [Collection No%d doesn't contain Document %q]", c.id, name)
	}

	return document, nil
}

func (c *Collection) RemoveDocument(name string) error {
	if _, present := c.documents[name]; !present {
		return fmt.Errorf("models/collection.RemoveDocument: [Collection No%d doesn't contain Document %q]", c.id, name)
	}
	delete(c.documents, name)

	return nil
}

func (c *Collection) Documents() []Document {
	return slices.Collect(maps.Values(c.documents))
}

func (c *Collection) Name() string {
	return c.name
}

func (c *Collection) UserUid() string {
	return c.userUid
}

func (c *Collection) Id() uint64 {
	return c.id
}
