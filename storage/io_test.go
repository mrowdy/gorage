package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/Slemgrim/gorage"
)

var tempDir = ".//tmp"
var testFile = "./test-files/test-file.html"

func TestMain(m *testing.M) {
	createTempDir(tempDir)

	retCode := m.Run()
	removeTempDir(tempDir)
	os.Exit(retCode)
}

func createTempDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}
}

func removeTempDir(dir string) {
	os.RemoveAll(dir)
}

func TestWriteFileToGivenDirectory(t *testing.T) {
	f := gorage.NewFile("test-file.html", []byte("fooo"))
	w := Io{
		BasePath:   tempDir,
		DirLength:  6,
		BufferSize: 1024,
	}
	err := w.Write(f)

	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("%s/%s/%s", tempDir, f.Hash()[:6], f.Hash())

	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Fatal("file not written")
	}
}

func TestWriteDirectoryBasedOnWritterSetting(t *testing.T) {
	f := gorage.NewFile("test-file.html", []byte("fooo"))
	w := Io{
		BasePath:   tempDir,
		DirLength:  4,
		BufferSize: 1024,
	}

	err := w.Write(f)

	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("%s/%s/%s", tempDir, f.Hash()[:4], f.Hash())

	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Fatal("file not written")
	}
}

func TestWriteRealFile(t *testing.T) {
	fi, err := os.Open(testFile)
	if err != nil {
		panic(err)
	}

	defer fi.Close()

	bytes, err := ioutil.ReadAll(fi)
	f := gorage.NewFile("file.png", bytes)
	w := Io{
		BasePath:   tempDir,
		DirLength:  6,
		BufferSize: 1024,
	}

	err = w.Write(f)

	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("%s/%s/%s", tempDir, f.Hash()[:6], f.Hash())

	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Fatal("file not written")
	}
}

func TestItReturnsErrorWhenPathIsNotWriteable(t *testing.T) {
	f := gorage.NewFile("test-file.html", []byte("invalid"))
	w := Io{
		BasePath:   "invalid-path",
		DirLength:  6,
		BufferSize: 1024,
	}

	err := w.Write(f)

	if err == nil {
		log.Fatal("No error reportet")
	}

	path := fmt.Sprintf("%s/%s", "./invalid-path", f.Hash()[:6])

	_, err = os.Stat(path)
	os.IsNotExist(err)

	if err == nil {
		log.Fatal("No error reportet")
	}
}

func TestItDoesntOverwriteExistingFiles(t *testing.T) {
	f1 := gorage.NewFile("test-file.html", []byte("fooooobar"))
	f2 := gorage.NewFile("test-file2.html", []byte("fooooobar"))
	w := Io{
		BasePath:   tempDir,
		DirLength:  6,
		BufferSize: 1024,
	}

	w.Write(f1)
	path1 := fmt.Sprintf("%s/%s/%s", tempDir, f1.Hash()[:6], f1.Hash())
	file1, _ := os.Stat(path1)
	m1 := file1.ModTime()

	w.Write(f2)
	path2 := fmt.Sprintf("%s/%s/%s", tempDir, f2.Hash()[:6], f2.Hash())
	file2, _ := os.Stat(path2)
	m2 := file2.ModTime()

	if m1 != m2 {
		log.Fatal("file was overwritten")
	}
}

func TestErrorWhenFileDoesntExist(t *testing.T) {
	w := Io{
		BasePath:   tempDir,
		DirLength:  6,
		BufferSize: 1024,
	}
	_, err := w.Read("asdfasdfasdf")

	if err == nil {
		log.Fatal("No error for missing file")
	}
}

func TestBasicReadWrite(t *testing.T) {
	f1 := gorage.NewFile("Test", []byte("TestBasicReadWrite"))
	w := Io{
		BasePath:   tempDir,
		DirLength:  6,
		BufferSize: 1024,
	}

	w.Write(f1)
	f2, _ := w.Read(f1.Hash())

	if bytes.Compare(f1.Content, f2.Content) != 0 {
		log.Fatal("write and read doesn't work")
	}
}
