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
	id       string
	title    string
	location string
	salary   string
	summary  string
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatusCode(statusCode int) {
	if statusCode != 200 {
		log.Fatalln("Request failed with status:", statusCode)
	}
}

func getURL(page int) string {
	var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50&start=%d"
	return fmt.Sprintf(baseURL, page*50)
}

func getNoPage(URL string) int {
	lastNum := 0
	res, err := http.Get(URL)

	checkErr(err)
	checkStatusCode(res.StatusCode)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".pagination-list>li>b").Each(func(i int, s *goquery.Selection) {
		lastNumStr, err := s.Html()
		checkErr(err)
		lastNum, err = strconv.Atoi(strings.TrimSpace(lastNumStr))
		checkErr(err)
	})
	return lastNum
}

func cleanString(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".slaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	return extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func getPage(URL string) []extractedJob {
	var jobs []extractedJob
	fmt.Println("Requesting to " + URL)
	res, err := http.Get(URL)

	checkErr(err)
	checkStatusCode(res.StatusCode)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	jobCards := doc.Find(".jobsearch-SerpJobCard")
	jobCards.Each(func(i int, card *goquery.Selection) {
		jobs = append(jobs, extractJob(card))
	})

	return jobs
}

func getPages(noPages int) []extractedJob {
	var jobs []extractedJob
	for i := 0; i < noPages; i++ {
		jobs = append(jobs, getPage(getURL(i))...)
	}

	return jobs
}

func writeCSVLine(w *csv.Writer, jobSlice []string) {
	jwErr := w.Write(jobSlice)
	checkErr(jwErr)
}

func writeJobsAsCSV(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)

	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}

		writeCSVLine(w, jobSlice)
	}
}

func main() {
	fmt.Println(getURL(200))
	noPages := getNoPage(getURL(200))
	fmt.Println(noPages)
	jobs := getPages(noPages)

	fmt.Println(len(jobs))
	writeJobsAsCSV(jobs)
}
