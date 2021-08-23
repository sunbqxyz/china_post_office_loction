package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"testing"
)

func Test_getPost(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	list := getPost(52060, browser)
	fmt.Println(len(list))
}

func Test_thread(t *testing.T) {
	pageTotal := 52131
	thread := 10
	var threadSize int
	//5213
	threadSize = pageTotal / thread
	//剩余page
	//thread := pageTotal % thread

	for i := 0; i < thread; i++ {
		for j := i * threadSize; j <= (i+1)*threadSize; j++ {
			fmt.Printf("thread %d page %d \n", i, j)
		}
	}
}

func Test_thread1(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	posts := thread(10, browser)
	fmt.Println(len(posts))
}
