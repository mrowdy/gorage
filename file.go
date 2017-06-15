package gorage

import (
	"crypto/sha1"
	"fmt"
	"time"
)

type FileContent []byte

type File struct {
	ID         string
	Name       string
	Hash       string
	MimeType   string
	Content    FileContent
	Size       int
	UploadedAt time.Time
	DeletedAt  *time.Time
}

func (f *FileContent) CalculateHash() string {
	h := sha1.New()
	h.Write(*f)
	return fmt.Sprintf("%x", h.Sum(nil))
}
