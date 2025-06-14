package tfidf

import (
	"bufio"
	"fmt"
	"io"
	"math"
)

func ProcessReaderToTerms(data io.Reader) ([]string, error) {
	text := []string{}
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("tfidf/tfidfSerivce.ProcessReaderToTerms: [%w]", err)
	}
	return text, nil
}

func GetTermFrequency(document []string) map[string]uint64 {
	bagOfWords := map[string]uint64{}
	for _, word := range document {
		bagOfWords[word] = bagOfWords[word] + 1
	}
	return bagOfWords
}

func CalculateIdf(termsAmount uint64, termFrequency uint64) float64 {
	return math.Log10(float64(termsAmount) / float64(termFrequency))
}
