package models

import (
	"os"
)

type Document struct {
	id     string
	userId string
	name   string
	file   *os.File
}

func NewDocument(id, userId, name string, file *os.File) Document {
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

func (d *Document) Replace(file *os.File) {
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

func (d *Document) File() *os.File {
	return d.file
}
