package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Storage interface {
	Write(file *File) error
	Read(hash string) (File, error)
}

type IoStorage struct {
	basePath   string
	dirLength  int
	bufferSize int
}

func NewStorage(basePath string, dirLength int, bufferSize int) *IoStorage {
	s := new(IoStorage)
	s.basePath = basePath
	s.dirLength = dirLength
	s.bufferSize = bufferSize
	return s
}

func (s *IoStorage) Write(file *File) error {

	fileHash := file.Hash()

	dir := s.getDirPath(fileHash)
	err := createDir(dir)

	if err != nil {
		return err
	}

	path := s.getFilePath(fileHash)
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fo.Close()

	r := bufio.NewReader(bytes.NewBuffer(file.Content))

	writer := bufio.NewWriter(fo)

	buf := make([]byte, s.bufferSize)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := writer.Write(buf[:n]); err != nil {
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (s *IoStorage) Read(hash string) (File, error) {
	fi, err := os.Open(s.getFilePath(hash))
	if err != nil {
		return File{}, err
	}

	defer fi.Close()

	bytes, err := ioutil.ReadAll(fi)

	if err != nil {
		return File{}, err
	}

	file := new(File)
	file.Content = bytes

	return *file, nil
}

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		error := os.Mkdir(dir, 0777)
		if error != nil {
			return error
		}
	}

	return nil
}

func (s *IoStorage) getDirPath(hash string) string {
	dir := make([]string, 0)
	dir = append(dir, s.basePath)
	dir = append(dir, hash[:s.dirLength])
	return strings.Join(dir, "/")
}

func (s *IoStorage) getFilePath(hash string) string {
	path := make([]string, 0)
	path = append(path, s.basePath)
	path = append(path, hash[:s.dirLength])
	path = append(path, hash)
	return strings.Join(path, "/")
}
