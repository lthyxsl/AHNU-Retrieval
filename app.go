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
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var cookies []*http.Cookie
var client = &http.Client{}
var wg sync.WaitGroup
var info map[int]string

type Seat struct {
	Name string
	Href string
}

func init() {
	info = make(map[int]string, 10)
	info[1] = "http://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=1&fid=1"   // 二 南
	info[2] = "http://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=6&fid=3"   // 三 南 自然科学
	info[3] = "http://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=5&fid=4"   // 三 北 社科 一
	info[4] = "http://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=3&fid=5"   // 四 南 社科三
	info[5] = "http://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=4&fid=6"   // 四 北 社科二
	info[6] = "http://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=13&fid=9"  // 三 公共区域 东
	info[7] = "http://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=14&fid=9"  // 三 公共区域 西
	info[8] = "http://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=15&fid=10" // 四 公共区域 东
	info[9] = "http://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=16&fid=10" // 四 公共区域 西
}

func login(stuId, password string) bool {
	form := url.Values{}
	form.Set("__VIEWSTATE", "/wEPDwULLTE0MTcxNzMyMjZkZJoL/NVYL0T+r5y3cXpfEFEzXz+dxNVtb7TlDKf8jIxz")
	form.Set("__VIEWSTATEGENERATOR", "F2D227C8")
	form.Set("__EVENTVALIDATION", "/wEWBQKV1czoDALyj/OQAgKXtYSMCgKM54rGBgKj48j5D1AZa5C6Zak6btNjhoHWy1AzD9qoyayyu5qGeLnFyXKG")
	form.Set("tbUserName", stuId)
	form.Set("tbPassWord", stuId)
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

	cookies = resp.Cookies()
	return true
}

func printInfo() {
	fmt.Println("花津校区图书馆2楼南 			请输入 1")
	fmt.Println("花津校区图书馆3楼南自然科学		请输入 2")
	fmt.Println("花津校区图书馆3楼北社科一	 	请输入 3")
	fmt.Println("花津校区图书馆4楼北社科三		请输入 4")
	fmt.Println("花津校区图书馆4楼南社科二		请输入 5")
	fmt.Println("花津校区图书馆3楼公共区域东		请输入 6")
	fmt.Println("花津校区图书馆3楼公共区域西		请输入 7")
	fmt.Println("花津校区图书馆4楼公共区域东		请输入 8")
	fmt.Println("花津校区图书馆4楼公共区域西		请输入 9")

}

func main() {
	var stuId string
	var password string
	fmt.Printf("请输入学号： ")
	fmt.Scanln(&stuId)
	fmt.Printf("请输入密码： ")
	fmt.Scanln(&password)
	logState := login(stuId, password)
	if !logState {
		fmt.Println("登陆失败，程序退出")
		return
	}
	fmt.Println("恭喜你，登陆成功， 欢迎进入此系统")
	fmt.Println()
	printInfo()
	var input int
	fmt.Printf("请输入预约楼层选项: ")
	fmt.Scanln(&input)
	var date string
	fmt.Printf("请输入预约时间:（如格式 2021-03-08）")
	fmt.Scanln(&date)
	if date == "" {
		date = time.Now().Format("2006-01-02")
		fmt.Println("使用当天默认时间查询", date)
	}
	link, exists := info[input]
	if exists {
		res, err := client.Get(link)
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
			go check(s, date)
		}

		wg.Wait()
	} else {
		fmt.Println("输入选项不正确，程序退出")
		return
	}

	fmt.Println("********************")
}

func check(s Seat, date string) {
	data := make(map[string]interface{})
	data["atdate"] = date
	data["sid"] = s.Href[strings.LastIndex(s.Href, "=")+1:]
	b, _ := json.Marshal(data)
	buff := bytes.NewBuffer(b)
	req, _ := http.NewRequest("POST", "http://libzwxt.ahnu.edu.cn/SeatWx/ajaxpro/SeatManage.Seat,SeatManage.ashx", buff)
	//req.Header.Add("X-Powered-By", "ASP.NET")
	//req.Header.Add("Proxy-Connection", "keep-alive")
	//req.Header.Add("Accept", "*/*")
	//req.Header.Add("Host", "libzwxt.ahnu.edu.cn")
	//req.Header.Add("Origin", "http://libzwxt.ahnu.edu.cn")
	//req.Header.Add("Referer", s.Href)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Mobile Safari/537.36")
	req.Header.Add("X-AjaxPro-Method", "GetSetInfo")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	all, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("正在查找: %s 座位情况\n", s.Name)
	if !strings.Contains(string(all), `<i class='on'></i>`) {
		fmt.Printf("----------------------检索全天空闲座位成功 %s \n", s.Name)
	}
	wg.Done()

}
