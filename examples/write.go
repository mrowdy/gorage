package main

import (
	"fmt"

	"github.com/Slemgrim/gorage"
	"github.com/Slemgrim/gorage/meta"
	"github.com/Slemgrim/gorage/relation"
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

	s := storage.Io{
		BasePath:   "./examples/tmp",
		DirLength:  6,
		BufferSize: 1024,
	}

	r := relation.NewDb("relation", db)
	m := meta.NewDb("meta", db)

	gorage := gorage.NewGorage(s, r, m)

	file, err := gorage.Save("test-file", []byte("File Content"))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", file)
}
