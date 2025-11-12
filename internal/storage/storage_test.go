package storage_test

import (
	"io/fs"
	"testing"

	"github.com/boomhut/go-epub/internal/storage"
	"github.com/boomhut/go-epub/internal/storage/memory"
)

func TestReadFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  []byte
		wantErr  bool
	}{
		{
			name:     "read existing file",
			filename: "test.txt",
			content:  []byte("hello world"),
			wantErr:  false,
		},
		{
			name:     "read non-existent file",
			filename: "nonexistent.txt",
			content:  nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := memory.NewMemory()

			// Setup: write file if content is provided
			if tt.content != nil {
				err := fs.WriteFile(tt.filename, tt.content, 0644)
				if err != nil {
					t.Fatalf("Failed to setup test file: %v", err)
				}
			}

			// Test ReadFile
			got, err := storage.ReadFile(fs, tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(got) != string(tt.content) {
				t.Errorf("ReadFile() = %v, want %v", string(got), string(tt.content))
			}
		})
	}
}

func TestMkdirAll(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		perm    fs.FileMode
		wantErr bool
	}{
		{
			name:    "create single directory",
			dir:     "test/",
			perm:    0755,
			wantErr: false,
		},
		{
			name:    "create nested directories",
			dir:     "test/nested/deep/",
			perm:    0755,
			wantErr: false,
		},
		{
			name:    "create already existing directory",
			dir:     "existing/",
			perm:    0755,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := memory.NewMemory()

			// Setup: create existing directory for the last test
			if tt.name == "create already existing directory" {
				err := fs.Mkdir("existing", tt.perm)
				if err != nil {
					t.Fatalf("Failed to setup existing directory: %v", err)
				}
			}

			// Test MkdirAll - it creates the directory and all its parents
			err := storage.MkdirAll(fs, tt.dir, tt.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("MkdirAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMkdirAll_NestedPath(t *testing.T) {
	fs := memory.NewMemory()

	// Test creating deeply nested directories
	// MkdirAll should create the path itself and all parent directories
	err := storage.MkdirAll(fs, "a/b/c/d/e/", 0755)
	if err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}

	// Verify all directories were created
	dirs := []string{"a", "a/b", "a/b/c", "a/b/c/d", "a/b/c/d/e"}
	for _, dir := range dirs {
		_, err = fs.Stat(dir)
		if err != nil {
			t.Errorf("MkdirAll() did not create directory %q: %v", dir, err)
		}
	}
}
