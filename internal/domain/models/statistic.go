package models

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
)

const MaxStatisticTermsAmount = 50

type Statistic struct {
	terms        map[string]Term
	collectionId uint64
}

func NewStatistic(collectionId uint64, terms [MaxStatisticTermsAmount]Term) *Statistic {
	termsMap := make(map[string]Term, len(terms))
	for _, term := range terms {
		termsMap[term.word] = term
	}

	return &Statistic{
		terms:        make(map[string]Term, MaxStatisticTermsAmount),
		collectionId: collectionId,
	}
}

func NewEmptyStatistic(collectionId uint64) *Statistic {
	return &Statistic{
		terms:        make(map[string]Term, MaxStatisticTermsAmount),
		collectionId: collectionId,
	}
}

func (s *Statistic) AddTerm(term Term) error {
	if len(s.terms) >= MaxStatisticTermsAmount {
		return fmt.Errorf("models/statistic.AddTerm: [Statistic No%d is at max capacity]", s.collectionId)
	}

	if _, present := s.terms[term.word]; present {
		return fmt.Errorf("models/statistic.AddTerm: [Statistic No%d already contains Term %q]", s.collectionId, term.word)
	}
	s.terms[term.word] = term

	return nil
}

func (s *Statistic) FindTerm(word string) (Term, error) {
	term, present := s.terms[word]
	if !present {
		return term, fmt.Errorf("models/statistic.FindTerm: [Statistic No%d doesn't contain Term %q]", s.collectionId, term.word)
	}

	return term, nil
}

func (s *Statistic) Contains(word string) bool {
	_, present := s.terms[word]

	return present
}

func (s *Statistic) RemoveTerm(word string) error {
	if _, present := s.terms[word]; !present {
		return fmt.Errorf("models/statistic.RemoveTerm: [Statistic No%d doesn't contain Term %q]", s.collectionId, word)
	}

	return nil
}

func (s *Statistic) Terms() []Term {
	terms := slices.Collect(maps.Values(s.terms))
	slices.SortFunc(terms, compareTermsByIdf)

	return terms
}

func (s *Statistic) CollectionId() uint64 {
	return s.collectionId
}

func compareTermsByIdf(a, b Term) int {
	return cmp.Compare(a.idf, b.idf)
}
