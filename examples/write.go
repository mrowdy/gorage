package main

import (
	"fmt"

	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/gorage/persistance"
	"github.com/Slemgrim/gorage/storage"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", "gorage:gorage@/gorage?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	storage := storage.Io{
		BasePath:   "./examples/tmp",
		DirLength:  6,
		BufferSize: 1024,
	}

	persistance := persistance.Db{
		Db: db,
	}
	persistance.Init()

	gorage := gorage.NewGorage(storage, persistance)
	file, err := gorage.Save("./examples/example-files/test-file.html")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", file)
}
