package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

func getStautss(doc *goquery.Document) []string {
	stautsMessage := make([]string, 0, 10)
	stauts := make([]string, 0, 10)
	statuss := make([]string, 0, 10)
	count := 0
	doc.Find("td.c1").Each(func(j int, t *goquery.Selection) {
		doc.Find("td.left").Each(func(i int, s *goquery.Selection) {
			if count == 0 && t.Text() == "HP" {
				stautsMessage = append(stautsMessage, t.Text())
				count++
			} else if i == j-5 && t.Text() != "タイプ" && t.Text() != "英語名" {
				stautsMessage = append(stautsMessage, t.Text())
			}
			if i < 7 && i == j-5 {
				stauts = append(stauts, s.Text())
			}
		})

	})
	for index1, value1 := range stautsMessage {
		for index2, value2 := range stauts {
			if index1 == index2 {
				statuss = append(statuss, value1+value2)
			}
		}
	}
	return statuss
}
func main() {
	engine := gin.Default()

	//var n string

	engine.LoadHTMLGlob("template/*")
	engine.Static("/assets", "./assets")
	engine.POST("/", func(c *gin.Context) {
		n := c.PostForm("id")
		url := "https://yakkun.com/sm/zukan/n" + n
		// Getリクエスト
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		// 読み取り
		buf, _ := ioutil.ReadAll(res.Body)

		// 文字コード判定
		det := chardet.NewTextDetector()
		detRslt, _ := det.DetectBest(buf)
		// => EUC-JP

		// 文字コード変換
		bReader := bytes.NewReader(buf)
		reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)

		// HTMLパース
		doc, _ := goquery.NewDocumentFromReader(reader)
		// rslt := doc.Find("td.left").Text()

		name := doc.Find("tr.head").Text()
		statuss := getStautss(doc)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"name":   name,
			"stauts": statuss,
		})
	})
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	engine.Run(":8080")
}
