package main

import (
	"fmt"
	"github.com/Xuanwo/go-locale"
	"golang.org/x/text/language"
	"log"
)

func main() {
	tag, err := locale.Detect()
	if err != nil {
		log.Fatal(err)
	}
	r, c := tag.Base()
	fmt.Println("tag:", r, c)

	fmt.Println(language.Make("zh-HK").Base())

	// Have fun with language.Tag!

	tags, err := locale.DetectAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tags:", tags)
	// Get all available tags
}
