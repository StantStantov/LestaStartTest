package models

import (
	"fmt"
	"maps"
	"slices"
)

type Collection struct {
	documents map[string]Document
	id        string
	userId    string
	name      string
}

func NewCollection(id, userId, name string, documents []Document) *Collection {
	documentsMap := make(map[string]Document, len(documents))
	for _, document := range documents {
		documentsMap[document.name] = document
	}

	return &Collection{
		documents: documentsMap,
		id:        id,
		userId:    userId,
		name:      name,
	}
}

func NewEmptyCollection(id, userId, name string) *Collection {
	return &Collection{
		documents: make(map[string]Document, 0),
		id:        id,
		userId:    userId,
		name:      name,
	}
}

func (c *Collection) AddDocument(document Document) error {
	if _, present := c.documents[document.name]; present {
		return fmt.Errorf("models/collection.AddDocument: [Collection No%q already contains Document %q]", c.id, document.name)
	}
	c.documents[document.name] = document

	return nil
}

func (c *Collection) FindDocument(name string) (Document, error) {
	document, present := c.documents[name]
	if !present {
		return document, fmt.Errorf("models/collection.FindDocument: [Collection No%q doesn't contain Document %q]", c.id, name)
	}

	return document, nil
}

func (c *Collection) RemoveDocument(name string) error {
	if _, present := c.documents[name]; !present {
		return fmt.Errorf("models/collection.RemoveDocument: [Collection No%q doesn't contain Document %q]", c.id, name)
	}
	delete(c.documents, name)

	return nil
}

func (c *Collection) Documents() []Document {
	return slices.Collect(maps.Values(c.documents))
}

func (c *Collection) Id() string {
	return c.id
}

func (c *Collection) UserId() string {
	return c.userId
}

func (c *Collection) Name() string {
	return c.name
}
