package writer

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrNilBatch        = errors.New("cannot write nil batch")
	ErrCreateFile      = errors.New("failed to create file")
	ErrWriteToFile     = errors.New("failed to write to file")
	ErrInvalidDir      = errors.New("invalid output directory")
	ErrCreateDirectory = errors.New("failed to create directory")
)

type Writer struct {
	outputDirectory string
}

func New(outputDirectory string) (*Writer, error) {
	if outputDirectory == "" {
		return nil, ErrInvalidDir
	}

	if err := os.MkdirAll(outputDirectory, 0755); err != nil {
		return nil, ErrCreateDirectory
	}

	return &Writer{
		outputDirectory: outputDirectory,
	}, nil
}

func (w *Writer) Write(iteration int, batch []int) error {
	if batch == nil {
		return ErrNilBatch
	}

	if len(batch) == 0 {
		return nil
	}

	filename := filepath.Join(w.outputDirectory, fmt.Sprintf("tmp_batch_%d.txt", iteration))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCreateFile, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, num := range batch {
		if _, err := fmt.Fprintln(writer, num); err != nil {
			return fmt.Errorf("%w: %s", ErrWriteToFile, err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("%w: %s", ErrWriteToFile, err)
	}

	return nil
}
