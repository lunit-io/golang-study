package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

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

func extractJob(card *goquery.Selection, c chan extractedJob) {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".slaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func getPage(URL string, jobs_c chan []extractedJob) {
	var jobs []extractedJob
	job_c := make(chan extractedJob)
	fmt.Println("Requesting to " + URL)
	res, err := http.Get(URL)

	checkErr(err)
	checkStatusCode(res.StatusCode)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	jobCards := doc.Find(".jobsearch-SerpJobCard")
	jobCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, job_c)

	})

	for i := 0; i < jobCards.Length(); i++ {
		jobs = append(jobs, <-job_c)
	}

	jobs_c <- jobs
}

func getPages(noPages int) []extractedJob {
	var jobs []extractedJob
	c := make(chan []extractedJob)
	for i := 0; i < noPages; i++ {
		go getPage(getURL(i), c)
	}

	for i := 0; i < noPages; i++ {
		jobs = append(jobs, <-c...)
	}

	return jobs
}

func writeCSVLine(w *csv.Writer, jobSlice []string) {
	jwErr := w.Write(jobSlice)
	checkErr(jwErr)
}

func writeCSVLineWithWait(w *csv.Writer, jobSlice []string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()   // later send signal 'Done'
	defer mu.Unlock() // first unlock

	mu.Lock()
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

	// https://stackoverflow.com/questions/18207772/how-to-wait-for-all-goroutines-to-finish-without-using-time-sleep
	var wg sync.WaitGroup
	// https://stackoverflow.com/questions/29981050/concurrent-writing-to-a-file
	var mu sync.Mutex

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}

		wg.Add(1)
		go writeCSVLineWithWait(w, jobSlice, &wg, &mu)
		// writeCSVLine(w, jobSlice)
	}
	wg.Wait()
}

func writeJobsAsJSON(jobs []extractedJob) {
	file, err := os.Create("jobs.json")
	checkErr(err)

	fmt.Println(jobs)
	b, err := json.Marshal(jobs)
	checkErr(err)
	// os.Stdout.Write(b)
	fmt.Println(string(b))
	n, err := file.Write(b)
	checkErr(err)
	fmt.Println(n)
}

func main() {
	fmt.Println(getURL(200))
	noPages := getNoPage(getURL(200))
	fmt.Println(noPages)
	jobs := getPages(noPages)

	fmt.Println(len(jobs))
	writeJobsAsCSV(jobs)
	// writeJobsAsJSON(jobs)
}
