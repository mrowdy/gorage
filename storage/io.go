package storage

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Slemgrim/gorage"
)

type Io struct {
	BasePath   string
	DirLength  int
	BufferSize int
}

func (s Io) Write(file *gorage.File) error {
	dir := s.getDirPath(file.Hash)
	err := createDir(dir)

	if err != nil {
		return err
	}

	path := s.getFilePath(file.Hash)
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fo.Close()

	r := bufio.NewReader(bytes.NewBuffer(file.Content))

	writer := bufio.NewWriter(fo)

	buf := make([]byte, s.BufferSize)
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

func (s Io) Read(hash string) (gorage.FileContent, error) {
	fi, err := os.Open(s.getFilePath(hash))

	var empty gorage.FileContent

	if err != nil {
		return empty, err
	}

	defer fi.Close()

	bytes, err := ioutil.ReadAll(fi)

	if err != nil {
		return empty, err
	}

	return bytes, nil
}

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		error := os.MkdirAll(dir, 0755)
		if error != nil {
			return error
		}
	}

	return nil
}

func (s *Io) getDirPath(hash string) string {
	dir := make([]string, 0)
	dir = append(dir, s.BasePath)
	dir = append(dir, hash[:s.DirLength])
	return strings.Join(dir, "/")
}

func (s *Io) getFilePath(hash string) string {
	path := make([]string, 0)
	path = append(path, s.BasePath)
	path = append(path, hash[:s.DirLength])
	path = append(path, hash)
	return strings.Join(path, "/")
}
