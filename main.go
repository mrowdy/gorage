package main

import (
	"fmt"
)

func main() {
	s := NewStorage("./tmp/", 6, 1024)
	gorage := NewGorage(s)
	file, err := gorage.Save("files/derp.html")
	if err != nil {
		panic(err)
	}

	fmt.Println(file)
}
