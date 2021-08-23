package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-rod/rod"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	MaxThread = 10
)

func main() {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	pageTotal, err := getPageTotal(browser)
	if err != nil {
		panic(err)
	}
	//var posts []*ChinaPostInfo
	//for i := 0; i < pageTotal; i += 10 {
	//	posts = append(posts, getPost(i, browser)...)
	//}
	posts := thread(pageTotal, browser)
	WriteJson(posts)
}

func thread(pageTotal int, browser *rod.Browser) []*ChinaPostInfo {
	var wg sync.WaitGroup
	//每个线程需要抓取的页数
	var size = pageTotal / MaxThread
	//剩余页数
	var surplusSize = pageTotal % MaxThread
	var posts []*ChinaPostInfo
	for i := 0; i < MaxThread; i++ {
		wg.Add(1)
		go threadFunc(i, size, posts, browser, &wg)
	}
	if surplusSize != 0 {
		wg.Add(1)
		go threadFunc(MaxThread, surplusSize, posts, browser, &wg)
	}
	wg.Wait()
	return posts
}
func threadFunc(currentThreadIndex int, threadSize int, posts []*ChinaPostInfo, browser *rod.Browser, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for i := currentThreadIndex * threadSize; i <= (currentThreadIndex+1)*threadSize; i++ {
		posts = append(posts, getPost(i, browser)...)
	}
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
	var list []*ChinaPostInfo
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

//WriteJson 写入json file
func WriteJson(area []*ChinaPostInfo) {
	areaBytes, err := json.Marshal(area)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fileName := "dist/china-post-%d.json"
	currentTime := time.Now().UnixNano() / 1e6
	fileName = fmt.Sprintf(fileName, currentTime)
	err = ioutil.WriteFile(fileName, areaBytes, 0666)
	if err != nil {
		fmt.Printf("create file error: %s", err.Error())
		return
	}
}

//ChinaPostInfo 中国邮政详情
type ChinaPostInfo struct {
	Province string `json:"province"`
	City     string `json:"city"`
	County   string `json:"county"`
	Info     string `json:"info"`
	Post     string `json:"post"`
	Addr     string `json:"addr"`
}

//NewChinaPostInfo 新建 ChinaPostInfo
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
