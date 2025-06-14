package services

type IdGenerator interface {
	GenerateId() string
}

type IdGeneratorFunc func() string

func (f IdGeneratorFunc) GenerateId() string {
	return f()
}
