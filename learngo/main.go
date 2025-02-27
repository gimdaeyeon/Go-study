package main

import (
	"fmt"
	"learngo/mydict"
)

func main() {
	dict := mydict.Dictionary{"first": "First word"}
	baseWord := "hello"
	dict.Add(baseWord, "First")
	dict.Search(baseWord)
	dict.Delete(baseWord)

	word, err := dict.Search(baseWord)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(word)

}
