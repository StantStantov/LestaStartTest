package models

type Document struct {
	name string
	id   uint64
}

func NewDocument(id uint64, filename string) Document {
	return Document{
		name: filename,
		id:   id,
	}
}

func (d *Document) Name() string {
	return d.name
}

func (d *Document) Id() uint64 {
	return d.id
}
