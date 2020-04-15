package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// 电影信息
type Movie struct {
	name          string
	pictures      string
	director      string
	evaluationNum string
	eviews        string
}

// 电影集合
var Top250 []Movie

// 用于抓取信息的正则
var (
	MOVIENAME          = "<img width=\"100\" alt=\"(.*?)\""
	MOVIEPICTURES      = "src=\"(.*?)\" class=\"\""
	MOVIEDIRECTORY     = "导演: (.*?)&nbsp"
	MOVIEEVALUATiONNUM = "<span>(.*?)评价<\\/span>"
	MOVIEREVIEWS       = "<span class=\"inq\">(.*?)<"
)

var ch chan int = make(chan int, 10)

func main() {
	GetAllPagesMovieInfo()
	ApiSetup()
}

// 获取每页的 body 信息
func GetBody(url string) string {
	userAgent := `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36`
	c := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", userAgent)
	resp, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Failed to get the website information")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return string(body)
}

// 获取一页的信息
func GetOnePageMovieInfo(body string) {
	var aMovie Movie
	//fmt.Println(body)
	regName := regexp.MustCompile(MOVIENAME)
	names := regName.FindAllStringSubmatch(body, -1)

	regpictures := regexp.MustCompile(MOVIEPICTURES)
	pictures := regpictures.FindAllStringSubmatch(body, -1)

	regDirectory := regexp.MustCompile(MOVIEDIRECTORY)
	directorys := regDirectory.FindAllStringSubmatch(body, -1)

	regEvaluationNum := regexp.MustCompile(MOVIEEVALUATiONNUM)
	evaluationNums := regEvaluationNum.FindAllStringSubmatch(body, -1)

	regEviews := regexp.MustCompile(MOVIEREVIEWS)
	eviews := regEviews.FindAllStringSubmatch(body, -1)

	if len(names) != 25 || len(pictures) != 25 || len(directorys) != 25 || len(evaluationNums) != 25 || len(eviews) != 25 {
		fmt.Println("names:")
		fmt.Println(len(names))
		fmt.Println("pictures:")
		fmt.Println(len(pictures))
		fmt.Println("directorys:")
		fmt.Println(len(directorys))
		fmt.Println("evaluationNums:")
		fmt.Println(len(evaluationNums))
		fmt.Println("eviews:")
		fmt.Println(len(eviews))
	} else {
		for i := 0; i < 25; i++ {
			aMovie.name = names[i][1]
			aMovie.evaluationNum = evaluationNums[i][1]
			aMovie.eviews = eviews[i][1]
			aMovie.pictures = pictures[i][1]
			aMovie.director = directorys[i][1]
			Top250 = append(Top250, aMovie)
		}
	}
}

// 获取所有页的信息
func GetAllPagesMovieInfo() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			url := "https://movie.douban.com/top250?start=" + strconv.Itoa(i*25) + "&filter="
			body := GetBody(url)
			fmt.Println(url)
			GetOnePageMovieInfo(body)
			ch <- 0
		}(i)
		<-ch
	}
}

func ApiSetup() {
	r := gin.Default()
	r.GET("/top250", Handle)
	r.Run()
}

func Handle(c *gin.Context) {
	var data []gin.H

	for i := 0; i < 250; i++ {
		data = append(data, gin.H{
			"id":            i,
			"name":          Top250[i].name,
			"pictures":      Top250[i].pictures,
			"director":      Top250[i].director,
			"eviews":        Top250[i].eviews,
			"evaluationNum": Top250[i].evaluationNum,
		})
	}

	c.JSON(200, data)
}
