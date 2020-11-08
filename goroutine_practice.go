package main

import (
	"fmt"
	"time"
)

//GorountinePractice 연습용
func GorountinePractice() {
	c := make(chan string)
	people := [5]string{"nico", "flynn", "dal", "jihwan", "larry"}
	for _, person := range people {
		go isSexy(person, c)
	}

	for i := 0; i < len(people); i++ {
		fmt.Println(<-c)
	}
}

func isSexy(person string, c chan string) {

	time.Sleep(time.Second * 5)
	c <- person + "is sexy"
}

func sexyCount(person string) {
	for i := 0; i < 10; i++ {
		fmt.Println(person, "is sexy", i)
		time.Sleep(time.Second)
	}
}
