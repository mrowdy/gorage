package main

import (
	"crypto/sha1"
	"fmt"
)

type File struct {
	name    string
	Content []byte
}

func NewFile(name string, content []byte) *File {
	file := new(File)
	file.name = name
	file.Content = content
	return file
}

func (f *File) Hash() string {
	h := sha1.New()
	h.Write(f.Content)
	return fmt.Sprintf("%x", h.Sum(nil))
}
