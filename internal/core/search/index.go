package search

import (
	"bufio"
	"encoding/json"
	"go-search/internal/model"
	"runtime"
	"sync"
)

const maxLineSize = 1024 * 1024     // 1MB buffer capacity
const initialBufferSize = 64 * 1024 // 64KB initial buffer
const channelBufferSize = 256

var CpuNum = runtime.NumCPU()

type Index struct {
	shards []*Shard
	words  [][]rune

	buffer  *sync.Pool
	scanner *bufio.Scanner
	res     map[string][]model.Position
}

func NewIndex(scanner *bufio.Scanner) *Index {
	return &Index{
		scanner: scanner,
		// Init buffer pool to optimize GC
		buffer: &sync.Pool{
			New: func() any {
				return &Input{}
			},
		},
	}
}

func (idx *Index) initOptimalStructures(words []string) {
	idx.words = make([][]rune, len(words))
	for i, word := range words {
		idx.words[i] = []rune(word)
	}
	idx.shards = make([]*Shard, CpuNum)

	for i := range idx.shards {
		idx.shards[i] = &Shard{
			data:   make(map[string][]model.Position, len(idx.words)),
			input:  make(chan *Input, channelBufferSize),
			buffer: idx.buffer,
		}
	}
}

func (idx *Index) Search(words []string) {
	// Use a smaller initial buffer and let it grow if needed
	buf := make([]byte, initialBufferSize)
	idx.scanner.Buffer(buf, maxLineSize)

	// Better memory allocation
	idx.initOptimalStructures(words)

	var wg sync.WaitGroup
	wg.Add(CpuNum)
	for i := 0; i < CpuNum; i++ {
		go idx.shards[i].worker(&wg, idx.words)
	}

	lineNum := 0
	currentWorker := 0

	for idx.scanner.Scan() {
		lineNum++

		in := idx.buffer.Get().(*Input)
		in.text = []rune(idx.scanner.Text())
		in.line = lineNum

		idx.shards[currentWorker].input <- in
		currentWorker = (currentWorker + 1) % CpuNum
	}

	// Close all channels
	for _, sh := range idx.shards {
		close(sh.input)
	}

	// Wait until all workers finish
	wg.Wait()

	// Composite all togethers
	idx.res = idx.mergeShards()

	// Clear the buffer to help GC
	buf = nil
}

func (idx *Index) Close() {
	idx.shards = nil
	idx.words = nil
	idx.res = nil
	idx.scanner = nil

	// Force GC to return the memory back to OS
	//runtime.GC()
	//debug.FreeOSMemory()
}

func (idx *Index) mergeShards() map[string][]model.Position {
	res := make(map[string][]model.Position, len(idx.words))
	// Memory pre-allocation
	for _, word := range idx.words {
		counter := 0
		for _, shard := range idx.shards {
			key := string(word)
			if _, ok := shard.data[key]; ok {
				counter += len(shard.data[key])
			}
		}
		res[string(word)] = make([]model.Position, 0, counter)
	}

	// Native merge
	for _, shard := range idx.shards {
		for word, positions := range shard.data {
			if len(positions) > 0 {
				res[word] = append(res[word], positions...)
			}
		}
	}

	return res
}

func (idx *Index) GetJson() ([]byte, error) {
	jsonBytes, err := json.Marshal(idx.res)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

func (idx *Index) GetResult() map[string][]model.Position {
	return idx.res
}
