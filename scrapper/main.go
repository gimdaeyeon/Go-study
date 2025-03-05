package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string // 구인id
	title    string // 제목
	location string // 지역
	summary  string // 설명
	company  string // 회사명
}

var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=golang"
var viewURL string = "https://www.saramin.co.kr/zf_user/jobs/relay/view?isMypage=no&rec_idx=%s&view_type=search&searchword=golang&searchType=search&gz=1&t_ref_content=generic&t_ref=search"

func main() {
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages()

	for i := 1; i <= totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageUrl := baseURL + "&recruitPage=" + strconv.Itoa(page)

	fmt.Println("Requesting: ", pageUrl)

	res, err := http.Get(pageUrl)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".item_recruit")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("value")
	title := cleanString(card.Find(".job_tit>a").Text())
	location := cleanString(card.Find(".job_condition>span").First().Text())
	summary := cleanString(card.Find(".job_sector").Text())
	company := cleanString(card.Find(".corp_name>a").Text())

	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		summary:  summary,
		company:  company,
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
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

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Id", "Title", "Location", "Summary", "Company"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		url := fmt.Sprintf(viewURL, job.id)
		jobSlice := []string{url, job.title, job.location, job.summary, job.company}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
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
