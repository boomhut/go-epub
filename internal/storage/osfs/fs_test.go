package osfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewOSFS(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)
	if osfs == nil {
		t.Fatal("NewOSFS() returned nil")
	}
	if osfs.rootDir != tmpDir {
		t.Errorf("NewOSFS() rootDir = %v, want %v", osfs.rootDir, tmpDir)
	}
}

func TestOSFS_WriteFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)
	testContent := []byte("test content")
	testFile := "test.txt"

	err = osfs.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Verify file was created
	fullPath := filepath.Join(tmpDir, testFile)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("File content = %v, want %v", string(content), string(testContent))
	}
}

func TestOSFS_Mkdir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)
	testDir := "testdir"

	err = osfs.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	// Verify directory was created
	fullPath := filepath.Join(tmpDir, testDir)
	info, err := os.Stat(fullPath)
	if err != nil {
		t.Fatalf("Failed to stat created directory: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("Created path is not a directory")
	}
}

func TestOSFS_RemoveAll(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)
	testDir := "testdir"

	// Create a directory first
	err = osfs.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	// Remove it
	err = osfs.RemoveAll(testDir)
	if err != nil {
		t.Fatalf("RemoveAll() error = %v", err)
	}

	// Verify it was removed
	fullPath := filepath.Join(tmpDir, testDir)
	_, err = os.Stat(fullPath)
	if !os.IsNotExist(err) {
		t.Errorf("RemoveAll() did not remove directory")
	}
}

func TestOSFS_Create(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)
	testFile := "test.txt"

	file, err := osfs.Create(testFile)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer file.Close()

	// Write to the file
	testContent := []byte("test content")
	_, err = file.Write(testContent)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Close and verify
	file.Close()
	fullPath := filepath.Join(tmpDir, testFile)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("File content = %v, want %v", string(content), string(testContent))
	}
}

func TestOSFS_Stat(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)
	testFile := "test.txt"
	testContent := []byte("test content")

	// Create a file first
	err = osfs.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Stat it
	info, err := osfs.Stat(testFile)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	if info.Name() != testFile {
		t.Errorf("Stat() name = %v, want %v", info.Name(), testFile)
	}

	if info.Size() != int64(len(testContent)) {
		t.Errorf("Stat() size = %v, want %v", info.Size(), len(testContent))
	}
}

func TestOSFS_Open(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)
	testFile := "test.txt"
	testContent := []byte("test content")

	// Create a file first
	err = osfs.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Open and read it
	file, err := osfs.Open(testFile)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer file.Close()

	content := make([]byte, len(testContent))
	_, err = file.Read(content)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("File content = %v, want %v", string(content), string(testContent))
	}
}

func TestOSFS_Integration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)

	// Create nested directories
	err = osfs.Mkdir("a", 0755)
	if err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	err = osfs.Mkdir("a/b", 0755)
	if err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	// Write a file in the nested directory
	testFile := "a/b/test.txt"
	testContent := []byte("nested content")
	err = osfs.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Read it back
	file, err := osfs.Open(testFile)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	content := make([]byte, len(testContent))
	_, err = file.Read(content)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	file.Close() // Close before RemoveAll

	if string(content) != string(testContent) {
		t.Errorf("File content = %v, want %v", string(content), string(testContent))
	}

	// Remove the nested structure
	err = osfs.RemoveAll("a")
	if err != nil {
		t.Fatalf("RemoveAll() error = %v", err)
	}

	// Verify it's gone
	fullPath := filepath.Join(tmpDir, "a")
	_, err = os.Stat(fullPath)
	if !os.IsNotExist(err) {
		t.Errorf("RemoveAll() did not remove nested directory structure")
	}
}

func TestOSFS_OpenNonExistent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)

	_, err = osfs.Open("nonexistent.txt")
	if err == nil {
		t.Error("Open() should return error for non-existent file")
	}
}

func TestOSFS_StatNonExistent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "osfs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	osfs := NewOSFS(tmpDir)

	_, err = osfs.Stat("nonexistent.txt")
	if err == nil {
		t.Error("Stat() should return error for non-existent file")
	}
}
