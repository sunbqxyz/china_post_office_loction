package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"strconv"
	"strings"
	"time"
)

func main() {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	pageTotal, err := getPageTotal(browser)
	if err != nil {
		panic(err)
	}
	posts := []*ChinaPostInfo{}
	for i := 0; i < pageTotal; i += 10 {
		posts = append(posts, getPost(i, browser)...)
	}
	fmt.Println(len(posts))

}

//getPageTotal 获取总页数
func getPageTotal(browser *rod.Browser) (int, error) {
	page := browser.MustPage("http://www.chinapost.com.cn/html1/folder/181312/9531-1.htm").MustWaitLoad()
	defer page.MustClose()
	page = page.MustElement(`iframe[src="http://iframe.chinapost.com.cn/jsp/type/institutionalsite/SiteSearchJT.jsp?community=ChinaPostJT&"]`).MustFrame()
	eles := page.MustElements(`#ali > a`)
	ele := eles.Last()
	href := ele.MustAttribute(`href`)
	index := strings.LastIndex(*href, "=")
	s := *href
	return strconv.Atoi(s[index+1:])

}

//获取post
func getPost(currentPage int, browser *rod.Browser) (result []*ChinaPostInfo) {
	err := rod.Try(func() {
		result = retry(
			currentPage,
			browser,
		)
	})
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		fmt.Printf("抓取失败，睡眠30秒")
		return getPost(currentPage, browser)
	} else {
		return result
	}
}
func retry(currentPage int, browser *rod.Browser) []*ChinaPostInfo {
	page := browser.MustPage(fmt.Sprintf("http://iframe.chinapost.com.cn/jsp/type/institutionalsite/SiteSearchJT.jsp?community=ChinaPostJT&pos=%d", currentPage)).MustWaitLoad()
	defer page.MustClose()
	eles := page.MustElements(`.wangd2 > tbody >tr`)
	list := []*ChinaPostInfo{}
	for i, ele := range eles {
		if i == 0 {
			continue
		}
		if ele.MustAttribute(`style`) != nil {
			continue
		}
		children := ele.MustElements(`td`)
		chinaPost := NewChinaPostInfo(
			children[0].MustText(),
			children[1].MustText(),
			children[2].MustText(),
			children[3].MustText(),
			children[4].MustText(),
			children[5].MustText(),
		)
		fmt.Printf("%s %s %s %s\n", chinaPost.Province, chinaPost.City, chinaPost.County, chinaPost.Post)

		list = append(list)
	}
	return list
}

type ChinaPostInfo struct {
	Province string
	City     string
	County   string
	Info     string
	Post     string
	Addr     string
}

func NewChinaPostInfo(province, city, county, info, post, addr string) *ChinaPostInfo {
	return &ChinaPostInfo{
		province,
		city,
		county,
		info,
		post,
		addr,
	}
}
