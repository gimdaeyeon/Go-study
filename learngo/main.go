package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan string)
	people := [4]string{"nico", "flynn", "james", "john"}
	for _, person := range people {
		go isChill(person, c)
	}

	for i := 0; i < len(people); i++ {
		fmt.Println(<-c)
	}
}

func isChill(pesron string, c chan string) {
	time.Sleep(time.Second * 3)
	c <- pesron + " is so chill!"
}
