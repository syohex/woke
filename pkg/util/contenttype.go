package util

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	ErrFileEmpty   = errors.New("file is empty")
	ErrIsDir       = errors.New("file is a directory")
	ErrFileNotText = errors.New("file is not a text file")
)

func detectContentType(file io.Reader) string {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)

	return http.DetectContentType(buffer[:n])
}

func isTextFile(file *os.File) bool {
	contentType := detectContentType(file)

	return strings.HasPrefix(contentType, "text/plain")
}

func IsTextFileFromFilename(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return IsTextFile(f)
}

// IsTextFile returns an error if the file is not of content-type 'text/plain'
func IsTextFile(file *os.File) error {
	e, err := file.Stat()
	if err != nil {
		return err
	}
	if e.IsDir() {
		return ErrIsDir
	}

	if e.Size() == 0 {
		return ErrFileEmpty
	}

	if !isTextFile(file) {
		return ErrFileNotText
	}

	return nil
}
