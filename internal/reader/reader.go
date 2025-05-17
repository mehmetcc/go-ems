package reader

import (
	"bufio"
	"errors"
	"os"
	"strconv"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrParseNumber  = errors.New("failed to parse number")
	ErrScanFile     = errors.New("error while scanning file")
)

type FileReader struct {
	path    string
	file    *os.File
	scanner *bufio.Scanner
}

func New(path string) (*FileReader, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return &FileReader{
		path:    path,
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (r *FileReader) Read(batchSize int) ([]int, error) {
	batch := make([]int, 0, batchSize)

	for i := 0; i < batchSize && r.scanner.Scan(); i++ {
		num, err := strconv.Atoi(r.scanner.Text())
		if err != nil {
			return nil, ErrParseNumber
		}
		batch = append(batch, num)
	}

	if err := r.scanner.Err(); err != nil {
		return nil, ErrScanFile
	}

	if len(batch) == 0 {
		return nil, nil
	}

	return batch, nil
}

func (r *FileReader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}
