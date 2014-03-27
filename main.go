package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shxsun/go-sh"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// FormatTime: 20060102 15:04:05

/*
	My dear friend,

	When I wrote this, God and I knew what it meant.

	Now, God only knows.

	CC @Artwalk
*/

type UsedData struct {
	Date      string
	MainPages []MainPage
}

type MainPage struct {
	Id         int // story id
	Title      string
	ShareImage string // download img
}

type FinalData struct {
	Useddata []UsedData
	Pagemark []int
}

var IMG = "static/img/"

//------------------------------------Main------------------------------------------

func main() {

	pages := make(map[int]FinalData)
	autoUpdate(pages)

	m := martini.Classic()
	m.Use(martini.Static("static"))
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {

		r.HTML(200, "content", []interface{}{pages[1]})
	})

	m.Get("/page/:id", func(params martini.Params, r render.Render) {

		id := atoi(params["id"])
		r.HTML(200, "content", []interface{}{pages[id]})
	})

	m.Get("/id/:id", func(params martini.Params, r render.Render) {

		id := atoi(params["id"])
		r.HTML(200, "id", id)
	})

	m.Get("/date/**", func(r render.Render) {
		r.HTML(200, "content", []interface{}{pages[1]})
	})

	m.Get("/url/**", func(r render.Render) {
		r.HTML(200, "content", []interface{}{pages[1]})
	})

	http.ListenAndServe("0.0.0.0:8000", m)
	m.Run()
}

//------------------------------------Pages------------------------------------------

func zhihuDailyJson(str string) UsedData {

	sj, _ := simplejson.NewJson([]byte(str))

	news, _ := sj.Get("news").Array()
	tmp, _ := time.Parse("20060102", sj.Get("date").MustString())
	date := tmp.Format("2006.01.02 Monday")

	var mainpages []MainPage

	for _, a := range news {

		m := a.(map[string]interface{})

		shareimageurl := ""
		shareimage := ""
		title := ""

		url := m["url"].(string)
		id := atoi(url[strings.LastIndexAny(url, "/")+1:])

		if m["share_image"] != nil {

			shareimageurl = m["share_image"].(string)

		} else { // no share_imag
			title = m["title"].(string)
			shareimageurl = m["image"].(string)
			//fmt.Println(id, title, shareimage)
		}

		shareimage = shareImgUrlToFilename(shareimageurl)
		mainpages = append(mainpages, MainPage{id, title, shareimage})

	}

	return UsedData{Date: date, MainPages: mainpages}
}

func renderPages(days int, pages map[int]FinalData) {

	var pagemark []int
	date := time.Now()

	if date.Format("MST") == "UTC" {
		date = date.Add(time.Hour * 8)
	}

	memoreyCache := QueryData()

	for i := 1; i <= len(memoreyCache)/days; i += 1 {
		pagemark = append(pagemark, i)
	}

	var newMainPages []MainPage

	for i := 1; i <= len(memoreyCache)/days; i += 1 {

		var finaldata FinalData
		var useddata []UsedData

		if i == 1 && date.Format("15") > "07" {
			todaydata := zhihuDailyJson(todayData())
			useddata = append(useddata, todaydata)

			//downloadDayShareImg(todaydata.MainPages)
			newMainPages = append(newMainPages, todaydata.MainPages...)
		}

		for j := 0; j < days; j++ {
			key := date.Format("20060102")

			data, ok := memoreyCache[atoi(key)]
			if !ok {
				data = getBeforeData(key)
			}

			beforeday := zhihuDailyJson(data)

			if i == 1 && j == 0 { // comment this line if you are first `go run main.go`
				//downloadDayShareImg(beforeday.MainPages)
				newMainPages = append(newMainPages, beforeday.MainPages...)
			} // comment this line if you are first `go run main.go`

			useddata = append(useddata, beforeday)
			date = date.AddDate(0, 0, -1)
		}
		finaldata.Useddata = useddata
		finaldata.Pagemark = pagemark
		pages[i] = finaldata
	}

	downloadDayShareImg(newMainPages)
}

func autoUpdate(pages map[int]FinalData) {

	// init
	days := 3
	renderPages(days, pages)

	ticker := time.NewTicker(time.Hour) // update every per hour
	go func() {
		for t := range ticker.C {
			fmt.Println("renderPages at ", t)
			renderPages(days, pages)
		}
	}()

}

// ----------------------------Download----------------------------------------------

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func download(url string) {

	filename := shareImgUrlToFilename(url)
	index := strings.LastIndexAny(filename, "_")

	if index > -1 {

		if !Exist(IMG + "croped/" + filename) {

			resp, err := http.Get(url)
			checkErr(err)

			defer resp.Body.Close()

			file, err := os.Create(IMG + filename)
			checkErr(err)

			io.Copy(file, resp.Body)

			fmt.Println("download: "+url+" -> ", filename)

			cropImage(filename)
		}
	}
}

func muiltDownload(urls []string, threads int) {
	if threads == 1 {
		go func() {
			for _, url := range urls {
				download(url)
			}
		}()
	} else {
		threads /= 2
		mid := len(urls) / 2
		muiltDownload(urls[:mid], threads)
		muiltDownload(urls[mid:], threads)
	}
}

func cropImage(filename string) {
	session := sh.NewSession()
	session.Command("convert", filename, "-resize", "440>", "-crop", "x275+0+0", "croped/"+filename, sh.Dir(IMG)).Run()
	session.Command("rm", IMG+filename).Run()
}

// --------------------------------DataBase------------------------------------------
func getData(url string) string {
	resp, err := http.Get(url)
	checkErr(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	return string(body)
}

func getBeforeData(date string) string {
	url := "http://news.at.zhihu.com/api/1.2/news/before/" + date
	data := getData(url)

	writeToDB(atoi(date), data)

	return data
}

func todayData() string {
	url := "http://news.at.zhihu.com/api/1.2/news/latest"

	return getData(url)
}

func QueryData() map[int]string {

	memoryCache := make(map[int]string)

	db, err := sql.Open("sqlite3", "./main.db")
	checkErr(err)

	rows, err := db.Query("SELECT * FROM datainfo")
	checkErr(err)

	db.Close()

	for rows.Next() {
		var date int
		var data string
		err = rows.Scan(&date, &data)
		memoryCache[date] = data
	}

	return memoryCache
}

func QueryID(id int, data string) string {

	db, err := sql.Open("sqlite3", "./main.db")
	checkErr(err)

	rows, err := db.Query("SELECT * FROM id")
	checkErr(err)

	db.Close()

	for rows.Next() {
		var index int
		var data string
		err = rows.Scan(&index, &data)
		if id == index {
			return data
		}
	}

	return ""
}

func writeToDB(date int, data string) {

	db, err := sql.Open("sqlite3", "./main.db")
	checkErr(err)
	//插入数据
	stmt, err := db.Prepare("INSERT INTO datainfo(date, data) values(?,?)")
	checkErr(err)

	res, err := stmt.Exec(date, data)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	db.Close()
}

// -----------------------------------Tools------------------------------------------
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(s string) int {
	dateInt, _ := strconv.Atoi(s)
	return dateInt
}

func idToUrl(id int) string {
	return "http://daily.zhihu.com/api/1.2/news/" + strconv.Itoa(id)
}

func filenameToShareImgUrl(filename string) string {
	url := ""
	if !strings.Contains(filename, "-") {
		url = "http://d0.zhimg.com/" + strings.Replace(filename, "_", "/", 1)
	} else {
		url = strings.Replace(filename, "_", "/", -1)
		url = strings.Replace(url, "-", ":", -1)
	}

	return url
}

func shareImgUrlToFilename(shareImgUrl string) string {

	filename := ""

	if strings.Contains(shareImgUrl, "http://d0.zhimg.com/") {
		str := strings.Replace(shareImgUrl, "http://d0.zhimg.com/", "", 1)
		filename = strings.Replace(str, "/", "_", 1)
	} else {
		filename = strings.Replace(shareImgUrl, "/", "_", -1)
		filename = strings.Replace(filename, ":", "-", -1)
	}

	return filename
}

// notice: do not call at once
func downloadDayShareImg(mainpages []MainPage) {

	var urls []string

	for _, mainpage := range mainpages {
		urls = append(urls, filenameToShareImgUrl(mainpage.ShareImage))
	}

	// 8 thread
	muiltDownload(urls, 8)
}
