package search

import (
	"bufio"
	"go-search/internal/core/jobs"
)

type IndexFactory struct{}

func NewIndexFactory() *IndexFactory {
	return &IndexFactory{}
}

func (f *IndexFactory) CreateNewIndex(scanner *bufio.Scanner) jobs.Indexer {
	idx := NewIndex(scanner)

	return idx
}
