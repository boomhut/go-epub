package epub

import (
	"os"

	"github.com/boomhut/go-epub/internal/storage"
	"github.com/boomhut/go-epub/internal/storage/memory"
	"github.com/boomhut/go-epub/internal/storage/osfs"
)

// FSType identifies the storage implementation backing the EPUB writer.
type FSType int

// filesystem is the current filesytem used as the underlying layer to manage the files.
// See the storage.Use method to change it.
var filesystem storage.Storage = osfs.NewOSFS(os.TempDir())

const (
	// OsFS defines the local filesystem implementation.
	OsFS FSType = iota
	// MemoryFS defines an in-memory filesystem implementation.
	MemoryFS
)

// Use sets the default storage backend. It defaults to the local filesystem.
func Use(s FSType) {
	switch s {
	case OsFS:
		filesystem = osfs.NewOSFS(os.TempDir())
	case MemoryFS:
		//TODO
		filesystem = memory.NewMemory()
	default:
		panic("unexpected FSType")
	}
}
