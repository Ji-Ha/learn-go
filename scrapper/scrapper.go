package scrapper

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

type extractedJobs struct {
	id       string
	location string
	title    string
	salary   string
	summary  string
}

//Scrape Indeed by a term
func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []extractedJobs
	c := make(chan []extractedJobs)
	totalPages := getPages(baseURL) / 2
	for i := 0; i < totalPages; i++ {
		go getPage(baseURL, i, c)

	}

	for i := 0; i < totalPages; i++ {
		exJob := <-c
		jobs = append(jobs, exJob...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func getPages(baseURL string) int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = (s.Find("a").Length())
	})
	return pages
}

func getPage(baseURL string, page int, mainC chan<- []extractedJobs) {
	var jobs []extractedJobs
	c := make(chan extractedJobs)
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting:", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractedJob(card, c)

	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extractedJob(card *goquery.Selection, c chan<- extractedJobs) {
	id, _ := card.Attr("data-jk")
	title := CleanString(card.Find(".title>a").Text())
	location := CleanString(card.Find(".sjcl").Text())
	salary := CleanString(card.Find(".salaryText").Text())
	summary := CleanString(card.Find(".summary").Text())
	c <- extractedJobs{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func writeJobs(jobs []extractedJobs) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	c := make(chan []string)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	//이부분도 goroutine을 이용해서 하면 더 빠를것.
	for _, job := range jobs {
		go writeJob(job, c)
		jobSlice := <-c
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func writeJob(job extractedJobs, c chan<- []string) {
	jobSlice := []string{"https://kr.indeed.com/viewjobs?jk=" + job.id, job.title, job.location, job.salary, job.summary}
	c <- jobSlice
}

//CleanString cleans string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
