package epub

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestEpubWriteTo(t *testing.T) {
	e := NewEpub(testEpubTitle)
	var b bytes.Buffer
	n, err := e.WriteTo(&b)
	if err != nil {
		t.Fatal(err)
	}
	if int64(len(b.Bytes())) != n {
		t.Fatalf("Expected size %v, got %v", len(b.Bytes()), n)
	}
}

func TestWriteToErrors(t *testing.T) {
	t.Run("CSS", func(t *testing.T) {
		e := NewEpub(testEpubTitle)
		testWriteToErrors(t, e, e.AddCSS, "cover.css")
	})
	t.Run("Font", func(t *testing.T) {
		e := NewEpub(testEpubTitle)
		testWriteToErrors(t, e, e.AddFont, "redacted-script-regular.ttf")
	})
	t.Run("Image", func(t *testing.T) {
		e := NewEpub(testEpubTitle)
		testWriteToErrors(t, e, e.AddImage, "gophercolor16x16.png")
	})
	t.Run("Video", func(t *testing.T) {
		e := NewEpub(testEpubTitle)
		testWriteToErrors(t, e, e.AddVideo, "sample_640x360.mp4")
	})
	t.Run("Audio", func(t *testing.T) {
		e := NewEpub(testEpubTitle)
		testWriteToErrors(t, e, e.AddAudio, "sample_audio.wav")
	})
}

func testWriteToErrors(t *testing.T, e *Epub, adder func(string, string) (string, error), name string) {
	// Copy testdata to temp file
	data, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("cannot open testdata: %v", err)
	}
	defer data.Close()
	temp, err := ioutil.TempFile("", "temp")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	io.Copy(temp, data)
	temp.Close()
	// Add temp file to epub
	if _, err := adder(temp.Name(), ""); err != nil {
		t.Fatalf("unable to add temp file: %v", err)
	}
	// Delete temp file
	if err := os.Remove(temp.Name()); err != nil {
		t.Fatalf("unable to delete temp file: %v", err)
	}
	// Write epub to buffer
	var b bytes.Buffer
	if _, err := e.WriteTo(&b); err == nil {
		t.Fatal("Expected error")
	}
}

// Test that UnableToCreateEpubError.Error() is correctly formatted
func TestUnableToCreateEpubErrorMessage(t *testing.T) {
	err := &UnableToCreateEpubError{
		Path: "/invalid/path/test.epub",
		Err:  io.ErrClosedPipe,
	}
	errStr := err.Error()
	if errStr == "" {
		t.Error("UnableToCreateEpubError.Error() returned empty string")
	}
	// Just verify it contains key information
	expectedSubstrings := []string{"/invalid/path/test.epub", "Error creating EPUB"}
	for _, substr := range expectedSubstrings {
		if !contains(errStr, substr) {
			t.Errorf("UnableToCreateEpubError.Error() = %q, should contain %q", errStr, substr)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
