package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	summary  string
	company  string
}

var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=golang"

func main() {
	totalPages := getPages()

	for i := 1; i <= totalPages; i++ {
		getPage(i)
	}
}

func getPage(page int) {
	pageUrl := baseURL + "&recruitPage=" + strconv.Itoa(page)

	res, err := http.Get(pageUrl)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".item_recruit")

	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("value")
		title := card.Find(".job_tit>a").Text()
		location := card.Find(".job_condition>span").First().Text()
		fmt.Println(id, title, location)
	})

}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: ", res.StatusCode)
	}
}
