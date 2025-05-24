package models

type Term struct {
	word      string
	frequency uint64
	idf       float64
}

func NewTerm(word string, frequency uint64, idf float64) Term {
	return Term{
		word:      word,
		frequency: frequency,
		idf:       idf,
	}
}

func (t Term) Word() string {
	return t.word
}

func (t Term) Frequency() uint64 {
	return t.frequency
}

func (t Term) Idf() float64 {
	return t.idf
}
