package stores

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"fmt"
	"slices"
)

type InMemoryTermStore struct {
	values []models.Term
	length int
}

func NewInMemoryTermStore() *InMemoryTermStore {
	return &InMemoryTermStore{values: []models.Term{}, length: 0}
}

func (ts *InMemoryTermStore) Create(term models.Term) error {
	ts.values = append(ts.values, term)
	ts.length++
	return nil
}

func (ts *InMemoryTermStore) Read(id int) (models.Term, error) {
	if 0 > id || id >= ts.length {
		return models.Term{}, fmt.Errorf("InMemoryTermStore.Read: Term does not exist")
	}
	return ts.values[id], nil
}

func (ts *InMemoryTermStore) ReadAll() ([]models.Term, error) {
	return slices.Clone(ts.values), nil
}

func (ts *InMemoryTermStore) CountAll() (int, error) {
	return ts.length, nil
}

func (ts *InMemoryTermStore) Update(id int, term models.Term) error {
	if 0 > id || id >= ts.length {
		return fmt.Errorf("InMemoryTermStore.Update: Term does not exist")
	}
	ts.values[id] = term
	return nil
}

func (ts *InMemoryTermStore) Delete(id int) error {
	if 0 > id || id >= ts.length {
		return fmt.Errorf("InMemoryTermStore.Delete: Term does not exist")
	}
	ts.values = slices.Delete(ts.values, id, id+1)
	ts.length--
	return nil
}
