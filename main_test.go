package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"sync"
	"testing"
)

func Test_getPost(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	list := getPost(52060, browser)
	fmt.Println(len(list))
}

func Test_thread(t *testing.T) {
	totalPage := 40
	threadCount := 10
	threadPage := totalPage / threadCount
	for i := 0; i < threadCount; i++ {
		start := i*threadPage + 1
		end := (i + 1) * threadPage
		fmt.Printf("thread %d start %d end %d \r\n", i, start, end)
	}
}
func Test_thread2(t *testing.T) {
	totalPage := 48420 / 10
	threadCount := 10
	everPage := totalPage / threadCount
	everPageAdd := totalPage % threadCount
	starPage := 0
	endPage := 0
	var flag = 0
	for page := 0; page < threadCount; page++ {
		if page == 0 {
			starPage = 0
		} else {
			starPage = endPage + 1
		}
		if page < everPageAdd {
			endPage = (page+1)*everPage + flag
			flag++
		} else {
			endPage = (page+1)*everPage + flag
		}
		fmt.Printf("thread %d start %d end %d \r\n", page, starPage, endPage)
	}
}

func Test_thread1(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	posts := thread(48420, browser)
	fmt.Println(len(posts))
}

func Test_goroutine(t *testing.T) {
	totalPage := 48420 / 10
	threadCount := 10
	everPage := totalPage / threadCount
	everPageAdd := totalPage % threadCount
	startPage := 0
	endPage := 0
	var flag = 0
	var wg sync.WaitGroup
	for page := 0; page < threadCount; page++ {
		if page == 0 {
			startPage = 0
		} else {
			startPage = endPage + 1
		}
		if page < everPageAdd {
			endPage = (page+1)*everPage + flag
			flag++
		} else {
			endPage = (page+1)*everPage + flag
		}
		wg.Add(1)
		go threadTest(startPage, endPage, &wg)
	}
	defer wg.Wait()
}

func threadTest(startPage, endPage int, wg *sync.WaitGroup) {
	for i := startPage; i <= endPage; i++ {
		fmt.Printf("page %d \r\n", i)
	}
	defer wg.Done()
}
