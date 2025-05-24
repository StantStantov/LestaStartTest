package stores

import "Stant/LestaGamesInternship/internal/models"

type TermStore interface {
	Create(term models.Term) error
	Read(id int) (models.Term, error)
	ReadAll() ([]models.Term, error)
	CountAll() (int, error)
	Update(id int, term models.Term) error
	Delete(id int) error
}
