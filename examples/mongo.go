package main

import (
	"fmt"
	"io/ioutil"

	"path/filepath"

	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/gorage/meta"
	"github.com/Slemgrim/gorage/relation"
	"github.com/Slemgrim/gorage/storage"
	"gopkg.in/mgo.v2"
)

/*
* Example for writing, reading and deleting files with gorage
 */
func main() {
	// Create a new storage instance. The actual files will be stored inside storage. Gorage includes a fs storage but it
	// would be easily possible to write other storage drivers for memcache, redis, etc.
	s := storage.Io{
		BasePath:   "./examples/tmp",
		DirLength:  6,
		BufferSize: 1024,
	}

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"192.168.99.100:32771"},
		Database: "gorage",
		Username: "",
		Password: "",
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//The relation instance keeps a one to one relation between the stored file and its charakteristics (filesize, mime-type, etc)
	r := relation.Mongo{Collection: session.DB("gorage").C("relation")}

	//The meta instance handles the different uploaded files. When you upload 2 files with the same content
	//they will be both available through the meta instance. (filename, upload date, etc., context)
	m := meta.Mongo{Collection: session.DB("gorage").C("meta")}

	gorage := gorage.NewGorage(s, r, m)

	//Load an existing file into a byte buffer (gorage only works with byte buffers)
	filename := "./examples/test-file.html"
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	//Save the file content with its filename
	savedFile1, err := gorage.Save(filepath.Base(filename), file, nil)
	if err != nil {
		panic(err)
	}

	//Save the same file again. Gorage will link this one to the previous saved file without
	//creating a new file in the filesystem
	//You can add a context as third param
	savedFile2, err := gorage.Save(filepath.Base(filename), file, 0.1)
	if err != nil {
		panic(err)
	}

	//Load the file with the ID given from the save method
	loadedFile, err := gorage.Load(savedFile1.ID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", loadedFile)

	//Delete the file with its ID
	err = gorage.Delete(savedFile1.ID)
	if err != nil {
		panic(err)
	}

	//Load now returns an error because the file was deleted. The file stil exists in the filesystem
	//because savedFile2 still references it
	_, err = gorage.Load(savedFile1.ID)
	if err != nil {
		fmt.Println(err)
	}

	//Loading savedFile2 still works fine
	_, err = gorage.Load(savedFile2.ID)
	if err != nil {
		fmt.Println(err)
	}

	//After deleting savedFile2 the actualFile will be deleted from the filesystem
	err = gorage.Delete(savedFile2.ID)
	if err != nil {
		panic(err)
	}
}
