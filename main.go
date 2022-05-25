package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// var cookies []*http.Cookie
var client = &http.Client{}
var wg sync.WaitGroup
var config = Config{}

type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var finalSeats = StringSlice{}

type Seat struct {
	Name string
	Href string
}

type UrlsItem struct {
	Index int    `json:"index"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

type Config struct {
	TbUserName string     `json:"tbUserName"`
	TbPassWord string     `json:"tbPassWord"`
	Date       string     `json:"date"`
	Option     int        `json:"option"`
	Urls       []UrlsItem `json:"urls"`
}

func init() {
	bytes, err := os.ReadFile("./conf/config.json")
	if err != nil {
		fmt.Println("加载配置文件失败")
		os.Exit(-1)
	}
	config = Config{}
	json.Unmarshal(bytes, &config)
	if config.Date == "" {
		config.Date = time.Now().Format("2006-01-02")
		fmt.Println("使用当天默认时间查询", config.Date)
	}
}

func login() bool {
	form := url.Values{}
	form.Set("__VIEWSTATE", "/wEPDwULLTE0MTcxNzMyMjZkZAl5GTLNAO7jkaD1B+BbDzJTZe4WiME3RzNDU4obNxXE")
	form.Set("__VIEWSTATEGENERATOR", "F2D227C8")
	form.Set("__EVENTVALIDATION", "/wEWBQK1odvtBQLyj/OQAgKXtYSMCgKM54rGBgKj48j5D4sJr7QMZnQ4zS9tzQuQ1arifvSWo1qu0EsBRnWwz6pw")
	form.Set("tbUserName", config.TbUserName)
	form.Set("tbPassWord", config.TbPassWord)
	form.Set("Button1", "登 录")
	form.Set("hfurl", "")
	b := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequest("POST", "http://libzwxt.ahnu.edu.cn/SeatWx/login.aspx?url=http%3a%2f%2flibzwxt.ahnu.edu.cn%2fSeatWx%2findex.aspx", b)
	if err != nil {
		err := fmt.Errorf("登录发起post请求时，http.NewRequest错误:%v", err)
		fmt.Println(err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar

	resp, err := client.Do(req)

	if err != nil {
		err := fmt.Errorf("登录发起post请求时，client.Do错误:%v", err)
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err := fmt.Errorf("登录发起post请求时，response.StatusCode:%v", resp.StatusCode)
		fmt.Println(err)
		return false
	}

	// cookies = resp.Cookies()
	// fmt.Printf("cookies: %#v\n", cookies)
	return true
}

func main() {

	if !login() {
		fmt.Println("登陆失败，程序退出")
		return
	}
	request()
}

// 开始请求
func request() {
	urlItem := config.Urls[config.Option]
	res, err := client.Get(urlItem.Url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	children := doc.Find("#ulSeat li").Children()
	seats := make([]Seat, 0, children.Length())
	// Find the review items
	doc.Find("#ulSeat li").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		name := s.Find("a").Text()
		href, _ := s.Find("a").Attr("href")
		//fmt.Printf("Review %d: %s-%s \n", i, name, href)
		seat := Seat{
			Name: name,
			Href: "http://libzwxt.ahnu.edu.cn/SeatWx/" + href,
		}
		seats = append(seats, seat)
	})

	for _, s := range seats {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(2000)))
		wg.Add(1)
		go check(s, config.Date)
	}
	wg.Wait()

	fmt.Println("********************")
	if len(finalSeats) > 0 {
		fmt.Println("********************完整座位列表")
		sort.Sort(finalSeats)
		for _, seat := range finalSeats {
			fmt.Printf("%v ", seat)
		}
		fmt.Println()
	}
}

func check(s Seat, date string) {
	data := make(map[string]interface{})
	data["atdate"] = date
	data["sid"] = s.Href[strings.LastIndex(s.Href, "=")+1:]
	b, _ := json.Marshal(data)
	buff := bytes.NewBuffer(b)
	req, _ := http.NewRequest("POST", "http://libzwxt.ahnu.edu.cn/SeatWx/ajaxpro/SeatManage.Seat,SeatManage.ashx", buff)
	req.Header.Add("X-Powered-By", "ASP.NET")
	//req.Header.Add("Proxy-Connection", "keep-alive")
	//req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "libzwxt.ahnu.edu.cn")
	req.Header.Add("Origin", "http://libzwxt.ahnu.edu.cn")
	//req.Header.Add("Referer", s.Href)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Mobile Safari/537.36")
	req.Header.Add("X-AjaxPro-Method", "GetSetInfo")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	all, _ := ioutil.ReadAll(res.Body)
	// fmt.Printf("-------------正在查找: %s 座位情况\n", s.Name)
	if !strings.Contains(string(all), `<i class='on'></i>`) {
		fmt.Printf(">>>>>>>>>>>>>>>> 座位号 %s 全天空闲 \n", s.Name)
		finalSeats = append(finalSeats, s.Name)
	}
	wg.Done()

}
