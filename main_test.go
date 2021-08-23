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
