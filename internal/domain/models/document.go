package models

import (
	"io"
)

type Document struct {
	id     string
	userId string
	name   string
	file   io.Reader
}

func NewDocument(id, userId, name string, file io.Reader) Document {
	return Document{
		id:     id,
		userId: userId,
		name:   name,
		file:   file,
	}
}

func (d *Document) Rename(newName string) {
	d.name = newName
}

func (d *Document) Replace(file io.Reader) {
	d.file = file
}

func (d *Document) Id() string {
	return d.id
}

func (d *Document) UserId() string {
	return d.userId
}

func (d *Document) Name() string {
	return d.name
}

func (d *Document) File() io.Reader {
	return d.file
}
