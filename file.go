package gorage

import (
	"crypto/sha1"
	"fmt"
)

type FileContent []byte

type File struct {
	ID       string
	Name     string
	Hash     string
	MimeType string
	Content  FileContent
	Size     int
}

func (f *File) CalculateHash() string {
	h := sha1.New()
	h.Write(f.Content)
	return fmt.Sprintf("%x", h.Sum(nil))
}
