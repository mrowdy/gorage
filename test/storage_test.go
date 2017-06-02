package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/slemgrim/treasury/treasury"
)

func TestMain(m *testing.M) {
	tempDir := "./../tmp"

	createTempDir(tempDir)

	retCode := m.Run()
	//removeTempDir(tempDir)
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
	f := treasury.NewFile("derp.html", []byte("fooo"))
	w := treasury.NewStorage("./../tmp", 6, 1024)
	err := w.Write(f)

	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("%s/%s/%s", "./../tmp", f.Hash()[:6], f.Hash())

	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Fatal("file not written")
	}
}

func TestWriteDirectoryBasedOnWritterSetting(t *testing.T) {
	f := treasury.NewFile("derp.html", []byte("fooo"))
	w := treasury.NewStorage("./../tmp", 4, 1024)
	err := w.Write(f)

	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("%s/%s/%s", "./../tmp", f.Hash()[:4], f.Hash())

	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Fatal("file not written")
	}
}

func TestWriteRealFile(t *testing.T) {
	fi, err := os.Open("./../files/derp.html")
	if err != nil {
		panic(err)
	}

	defer fi.Close()

	bytes, err := ioutil.ReadAll(fi)
	f := treasury.NewFile("file.png", bytes)
	w := treasury.NewStorage("./../tmp", 6, 1024)
	err = w.Write(f)

	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("%s/%s/%s", "./../tmp", f.Hash()[:6], f.Hash())

	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Fatal("file not written")
	}
}

func TestItReturnsErrorWhenPathIsNotWriteable(t *testing.T) {
	f := treasury.NewFile("derp.html", []byte("invalid"))
	w := treasury.NewStorage("./../invalid-path", 6, 1024)
	err := w.Write(f)

	if err == nil {
		log.Fatal("No error reportet")
	}

	path := fmt.Sprintf("%s/%s", "./../invalid-path", f.Hash()[:6])

	_, err = os.Stat(path)
	os.IsNotExist(err)

	if err == nil {
		log.Fatal("No error reportet")
	}
}

func TestItDoesntOverwriteExistingFiles(t *testing.T) {
	f1 := treasury.NewFile("derp.html", []byte("fooooobar"))
	f2 := treasury.NewFile("foo.html", []byte("fooooobar"))
	w := treasury.NewStorage("./../tmp", 6, 1024)

	w.Write(f1)
	path1 := fmt.Sprintf("%s/%s/%s", "./../tmp", f1.Hash()[:6], f1.Hash())
	file1, _ := os.Stat(path1)
	m1 := file1.ModTime()

	w.Write(f2)
	path2 := fmt.Sprintf("%s/%s/%s", "./../tmp", f2.Hash()[:6], f2.Hash())
	file2, _ := os.Stat(path2)
	m2 := file2.ModTime()

	if m1 != m2 {
		log.Fatal("file was overwritten")
	}
}

func TestErrorWhenFileDoesntExist(t *testing.T) {
	w := treasury.NewStorage("./../tmp", 6, 1024)
	_, err := w.Read("asdfasdfasdf")

	if err == nil {
		log.Fatal("No error for missing file")
	}
}

func TestBasicReadWrite(t *testing.T) {
	f1 := treasury.NewFile("Test", []byte("TestBasicReadWrite"))
	w := treasury.NewStorage("./../tmp", 6, 1024)

	w.Write(f1)
	f2, _ := w.Read(f1.Hash())

	if bytes.Compare(f1.Content, f2.Content) != 0 {
		log.Fatal("write and read doesn't work")
	}
}
