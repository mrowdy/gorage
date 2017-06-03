package main

import (
	"fmt"

	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/gorage/storage"
)

func main() {
	storage := storage.Io{
		BasePath:   "./examples/tmp",
		DirLength:  6,
		BufferSize: 1024,
	}

	gorage := gorage.NewGorage(storage)
	file, err := gorage.Save("./examples/example-files/test-file.html")
	if err != nil {
		panic(err)
	}

	fmt.Println(file)
}
