package gorage

import (
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
	Context    interface{}
	UploadedAt time.Time
	DeletedAt  *time.Time
}
