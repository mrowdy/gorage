package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/slemgrim/treasury/treasury"
)

func main() {

	db, err := gorm.Open("mysql", "treasury:test1234@/treasury?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if !db.HasTable(&treasury.DbFile{}) {
		fmt.Println("Create File Table")
		db.CreateTable(&treasury.DbFile{})
	}

	s := treasury.NewStorage("./tmp/", 6, 1024)
	treasury := treasury.NewTreasury(s, db)
	file, err := treasury.Save("files/derp.html")
	if err != nil {
		panic(err)
	}

	fmt.Println(file)
}
