// Package osfs implements the Storage interface for os' filesystems

package osfs

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/boomhut/go-epub/internal/storage"
)

// OSFS provides a storage.Storage implementation backed by the OS filesystem.
type OSFS struct {
	rootDir string
	fs.FS
}

// NewOSFS creates an OS-backed storage rooted at rootDir.
func NewOSFS(rootDir string) *OSFS {
	return &OSFS{
		rootDir: rootDir,
		FS:      os.DirFS(rootDir),
	}
}

// WriteFile writes data to a file relative to the root directory.
func (o *OSFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filepath.Join(o.rootDir, name), data, perm)
}

// Mkdir creates a directory relative to the root directory.
func (o *OSFS) Mkdir(name string, perm fs.FileMode) error {
	return os.Mkdir(filepath.Join(o.rootDir, name), perm)
}

// RemoveAll removes the path and all descendants relative to the root directory.
func (o *OSFS) RemoveAll(name string) error {
	return os.RemoveAll(filepath.Join(o.rootDir, name))
}

// Create opens a writable file relative to the root directory.
func (o *OSFS) Create(name string) (storage.File, error) {
	return os.Create(filepath.Join(o.rootDir, name))
}

// Stat retrieves file information relative to the root directory.
func (o *OSFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(o.rootDir, name))
}

// Open opens a file relative to the root directory.
func (o *OSFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(o.rootDir, name))
}
