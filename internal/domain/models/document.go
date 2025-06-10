package models

import (
	"os"
)

type Document struct {
	id   string
	name string
	file *os.File
}

func NewDocument(id, name string, file *os.File) Document {
	return Document{
		id:   id,
		name: name,
		file: file,
	}
}

func (d *Document) Id() string {
	return d.id
}

func (d *Document) Name() string {
	return d.name
}

func (d *Document) File() *os.File {
	return d.file
}
