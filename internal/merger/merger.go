package merger

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var (
	ErrOpenChunk   = errors.New("failed to open chunk file")
	ErrReadChunk   = errors.New("failed to read from chunk")
	ErrWriteOutput = errors.New("failed to write to output")
)

type Merger struct {
	tempDir    string
	outputPath string
	chunkCount int
	bufferSize int
}

func New(tempDir, outputPath string, chunkCount int) *Merger {
	return &Merger{
		tempDir:    tempDir,
		outputPath: outputPath,
		chunkCount: chunkCount,
		bufferSize: 8192, // default buffer size
	}
}

func (m *Merger) Merge() error {
	outFile, err := os.Create(m.outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	h := &MinHeap{}
	heap.Init(h)

	readers := make([]*bufio.Scanner, m.chunkCount)
	for i := 0; i < m.chunkCount; i++ {
		chunkPath := filepath.Join(m.tempDir, fmt.Sprintf("tmp_batch_%d.txt", i))
		file, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrOpenChunk, err)
		}
		defer file.Close()

		readers[i] = bufio.NewScanner(file)

		if readers[i].Scan() {
			num, err := strconv.Atoi(readers[i].Text())
			if err != nil {
				return fmt.Errorf("%w: %s", ErrReadChunk, err)
			}
			heap.Push(h, Item{value: num, chunkID: i})
		}
	}

	for h.Len() > 0 {
		item := heap.Pop(h).(Item)

		if _, err := fmt.Fprintln(writer, item.value); err != nil {
			return fmt.Errorf("%w: %s", ErrWriteOutput, err)
		}

		if readers[item.chunkID].Scan() {
			num, err := strconv.Atoi(readers[item.chunkID].Text())
			if err != nil {
				return fmt.Errorf("%w: %s", ErrReadChunk, err)
			}
			heap.Push(h, Item{value: num, chunkID: item.chunkID})
		}
	}

	return nil
}

func (m *Merger) Cleanup() error {
	for i := 0; i < m.chunkCount; i++ {
		chunkPath := filepath.Join(m.tempDir, fmt.Sprintf("tmp_batch_%d.txt", i))
		if err := os.Remove(chunkPath); err != nil {
			return fmt.Errorf("cleaning up chunk %d: %w", i, err)
		}
	}
	return nil
}
