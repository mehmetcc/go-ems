package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("valid directory", func(t *testing.T) {
		dir := t.TempDir()
		w, err := New(dir)
		if err != nil {
			t.Fatalf("New() error = %v, want nil", err)
		}
		if w.outputDirectory != dir {
			t.Errorf("outputDirectory = %v, want %v", w.outputDirectory, dir)
		}
	})

	t.Run("empty directory path", func(t *testing.T) {
		w, err := New("")
		if err != ErrInvalidDir {
			t.Errorf("New() error = %v, want %v", err, ErrInvalidDir)
		}
		if w != nil {
			t.Errorf("New() writer = %v, want nil", w)
		}
	})
}

func TestWriter_Write(t *testing.T) {
	tests := []struct {
		name      string
		batch     []int
		iteration int
		wantErr   error
	}{
		{
			name:      "successful write",
			batch:     []int{1, 2, 3},
			iteration: 0,
			wantErr:   nil,
		},
		{
			name:      "nil batch",
			batch:     nil,
			iteration: 0,
			wantErr:   ErrNilBatch,
		},
		{
			name:      "empty batch",
			batch:     []int{},
			iteration: 0,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			w, err := New(dir)
			if err != nil {
				t.Fatalf("Failed to create writer: %v", err)
			}

			err = w.Write(tt.iteration, tt.batch)
			if err != tt.wantErr {
				t.Errorf("Write() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && tt.batch != nil && len(tt.batch) > 0 {
				// Verify file contents
				filename := filepath.Join(dir, fmt.Sprintf("tmp_batch_%d.txt", tt.iteration))
				content, err := os.ReadFile(filename)
				if err != nil {
					t.Fatalf("Failed to read output file: %v", err)
				}

				// Convert content back to numbers and compare
				lines := strings.Split(strings.TrimSpace(string(content)), "\n")
				if len(lines) != len(tt.batch) {
					t.Errorf("Written batch has %d numbers, want %d", len(lines), len(tt.batch))
				}

				for i, line := range lines {
					num, err := strconv.Atoi(line)
					if err != nil {
						t.Errorf("Invalid number at line %d: %v", i, err)
					}
					if num != tt.batch[i] {
						t.Errorf("Number at position %d = %d, want %d", i, num, tt.batch[i])
					}
				}
			}
		})
	}
}

func TestWriter_WriteMultiple(t *testing.T) {
	dir := t.TempDir()
	w, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}

	batches := []struct {
		iteration int
		numbers   []int
	}{
		{0, []int{1, 2, 3}},
		{1, []int{4, 5, 6}},
		{2, []int{7, 8, 9}},
	}

	for _, batch := range batches {
		if err := w.Write(batch.iteration, batch.numbers); err != nil {
			t.Errorf("Write() error = %v", err)
		}
	}

	// Verify all files exist with correct content
	for _, batch := range batches {
		filename := filepath.Join(dir, fmt.Sprintf("tmp_batch_%d.txt", batch.iteration))
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Errorf("Failed to read file for iteration %d: %v", batch.iteration, err)
			continue
		}

		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		if len(lines) != len(batch.numbers) {
			t.Errorf("Iteration %d: got %d numbers, want %d", batch.iteration, len(lines), len(batch.numbers))
		}
	}
}
