package search

import (
	"go-search/internal/model"
	"sync"
	"unicode"
)

type Input struct {
	text []rune
	line uint32
}

type Shard struct {
	data   map[string][]model.Position
	input  chan *Input
	buffer *sync.Pool
}

func (sh *Shard) matchAt(text []rune, word []rune, start int) bool {
	for j := range word {
		if unicode.ToLower(text[start+j]) != unicode.ToLower(word[j]) {
			return false
		}
	}
	return true
}

func (sh *Shard) worker(wg *sync.WaitGroup, words [][]rune) {
	defer wg.Done()
	// Iterate until channel is closed
	for input := range sh.input {
		for _, word := range words {
			// Iterate over text and find matches
			textLen := len(input.text)
			wordLen := len(word)
			if wordLen == 0 || wordLen > textLen {
				continue
			}
			for i := 0; i <= textLen-wordLen; {
				if !sh.matchAt(input.text, word, i) {
					i++
					continue
				}
				position := model.Position{input.line, uint32(i + 1)}
				key := string(word)
				if _, ok := sh.data[key]; !ok {
					// Micro optimization
					sh.data[key] = make([]model.Position, 0, 4)
				}
				sh.data[key] = append(sh.data[key], position)
				i += wordLen
			}
		}

		// Release memory to the buffer pool
		sh.buffer.Put(input)
	}
}
