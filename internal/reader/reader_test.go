// STUB: generated
package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileReader(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("New", func(t *testing.T) {
		tests := []struct {
			name    string
			path    string
			wantErr error
		}{
			{
				name:    "file not found",
				path:    filepath.Join(tempDir, "nonexistent.txt"),
				wantErr: ErrFileNotFound,
			},
			{
				name:    "valid file",
				path:    createTestFile(t, tempDir, "valid.txt", []string{"123", "456"}),
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				reader, err := New(tt.path)
				if err != tt.wantErr {
					t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err == nil {
					defer reader.Close()
				}
			})
		}
	})

	t.Run("Read", func(t *testing.T) {
		tests := []struct {
			name      string
			content   []string
			batchSize int
			want      []int
			wantErr   error
		}{
			{
				name:      "read complete batch",
				content:   []string{"1", "2", "3"},
				batchSize: 3,
				want:      []int{1, 2, 3},
				wantErr:   nil,
			},
			{
				name:      "read partial batch",
				content:   []string{"1", "2"},
				batchSize: 3,
				want:      []int{1, 2},
				wantErr:   nil,
			},
			{
				name:      "invalid number",
				content:   []string{"1", "invalid", "3"},
				batchSize: 3,
				want:      nil,
				wantErr:   ErrParseNumber,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path := createTestFile(t, tempDir, "test.txt", tt.content)
				reader, err := New(path)
				if err != nil {
					t.Fatalf("Failed to create reader: %v", err)
				}
				defer reader.Close()

				got, err := reader.Read(tt.batchSize)
				if err != tt.wantErr {
					t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !sliceEqual(got, tt.want) {
					t.Errorf("Read() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Read Multiple Batches", func(t *testing.T) {
		content := []string{"1", "2", "3", "4", "5"}
		path := createTestFile(t, tempDir, "multi.txt", content)
		reader, err := New(path)
		if err != nil {
			t.Fatalf("Failed to create reader: %v", err)
		}
		defer reader.Close()

		batch1, err := reader.Read(2)
		if err != nil {
			t.Fatalf("Failed to read first batch: %v", err)
		}
		if !sliceEqual(batch1, []int{1, 2}) {
			t.Errorf("First batch = %v, want [1 2]", batch1)
		}

		batch2, err := reader.Read(2)
		if err != nil {
			t.Fatalf("Failed to read second batch: %v", err)
		}
		if !sliceEqual(batch2, []int{3, 4}) {
			t.Errorf("Second batch = %v, want [3 4]", batch2)
		}

		batch3, err := reader.Read(2)
		if err != nil {
			t.Fatalf("Failed to read third batch: %v", err)
		}
		if !sliceEqual(batch3, []int{5}) {
			t.Errorf("Third batch = %v, want [5]", batch3)
		}

		// Verify EOF
		batch4, err := reader.Read(2)
		if err != nil || batch4 != nil {
			t.Errorf("Expected EOF (nil, nil), got (%v, %v)", batch4, err)
		}
	})
}

func TestSequentialReads(t *testing.T) {
	tempDir := t.TempDir()
	content := []string{"1", "2", "3", "4", "5", "6"}
	path := createTestFile(t, tempDir, "sequential.txt", content)

	reader, err := New(path)
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}
	defer reader.Close()

	// First read should get first two numbers
	batch1, err := reader.Read(2)
	if err != nil {
		t.Fatalf("First read failed: %v", err)
	}
	if !sliceEqual(batch1, []int{1, 2}) {
		t.Errorf("First batch = %v, want [1 2]", batch1)
	}

	// Second read should get next two numbers
	batch2, err := reader.Read(2)
	if err != nil {
		t.Fatalf("Second read failed: %v", err)
	}
	if !sliceEqual(batch2, []int{3, 4}) {
		t.Errorf("Second batch = %v, want [3 4]", batch2)
	}

	// Third read should get last two numbers
	batch3, err := reader.Read(2)
	if err != nil {
		t.Fatalf("Third read failed: %v", err)
	}
	if !sliceEqual(batch3, []int{5, 6}) {
		t.Errorf("Third batch = %v, want [5 6]", batch3)
	}

	// Fourth read should return nil, nil (EOF)
	batch4, err := reader.Read(2)
	if err != nil || batch4 != nil {
		t.Errorf("Expected EOF (nil, nil), got (%v, %v)", batch4, err)
	}
}

func createTestFile(t *testing.T, dir, name string, content []string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	for _, line := range content {
		if _, err := file.WriteString(line + "\n"); err != nil {
			t.Fatalf("Failed to write to test file: %v", err)
		}
	}
	return path
}

func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
