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
	threadCount = 10
)

func main() {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	pageTotal, err := getPageTotal(browser)
	if err != nil {
		panic(err)
	}
	posts := thread(pageTotal, browser)
	WriteJson(posts)
}

//forEach 普通方式
func forEach(totalPage int, browser *rod.Browser) []*ChinaPostInfo {
	var posts []*ChinaPostInfo
	for i := 0; i <= totalPage; i += 10 {
		posts = append(posts, getPost(i, browser)...)
	}
	return posts
}

//thread 多协程抓取
func thread(pageTotal int, browser *rod.Browser) []*ChinaPostInfo {
	var (
		totalPage = pageTotal / 10
		//每个线程要处理的页数
		everPage = totalPage / threadCount
		//前几个线程要增加一页
		everPageAdd = totalPage % threadCount
		startPage   = 0
		endPage     = 0
		flag        = 0
		wg          sync.WaitGroup
		posts       []*ChinaPostInfo
	)
	var value = make(chan []*ChinaPostInfo)
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
		go threadFunc(startPage, endPage, &wg, value, browser)
	}
	go func() {
		wg.Wait()
		defer close(value)
	}()
	for data := range value {
		if data == nil {
			break
		}
		posts = append(posts, data...)
	}

	return posts
}

func threadFunc(startPage, endPage int, wg *sync.WaitGroup, chanValue chan []*ChinaPostInfo, browser *rod.Browser) {
	var posts []*ChinaPostInfo
	for page := startPage; page <= endPage; page++ {
		posts = append(posts, getPost(page, browser)...)
	}
	chanValue <- posts
	defer wg.Done()
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
	page := browser.MustPage(fmt.Sprintf("http://iframe.chinapost.com.cn/jsp/type/institutionalsite/SiteSearchJT.jsp?community=ChinaPostJT&pos=%d", currentPage*10)).MustWaitLoad()
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
		fmt.Printf("%s %s %s %s %s\n", chinaPost.Province, chinaPost.City, chinaPost.County, chinaPost.Addr, chinaPost.Post)

		list = append(list, chinaPost)
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
