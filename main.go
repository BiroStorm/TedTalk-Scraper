package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

type tedTalk struct {
	Title       string
	Description string
	//Thumb       string
	Views  string
	Author string
	Date   string
	Tags   []string
	Link   string
}

var tagsCollector *colly.Collector

func getTags(link string, c chan string) {
	defer close(c)
	fmt.Println("Arriva a getTags")

	tagsCollector.Visit(link)
}

type listTalk []tedTalk

func main() {

	allTalks := make(listTalk, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("www.ted.com"),
	)

	var visitedLink map[string]struct{} = make(map[string]struct{}, 0)

	detailCollector := c.Clone()
	tagsCollector = c.Clone()

	//talks := make([]tedTalk, 100)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		if e.Attr("class") == " ga-link" {
			if _, ok := visitedLink[e.Attr("href")]; !ok {
				visitedLink[e.Attr("href")] = struct{}{}
				//fmt.Printf("Link found: %s\n", e.Attr("href"))
				detailCollector.Visit("https://www.ted.com" + e.Attr("href"))
			}

		}

	})

	detailCollector.OnHTML(`html`, func(e *colly.HTMLElement) {

		/*
			fmt.Printf("Talk Title: %s\n", e.ChildAttr(`meta[itemprop="name"]`, "content"))
			fmt.Printf("Talk Descr: %s\n", e.ChildAttr(`meta[itemprop="description"]`, "content"))
			fmt.Printf("Talk Link: %s\n", e.ChildAttr(`meta[property="og:url"]`, "content"))
			fmt.Printf("Talk Author: %s\n", e.ChildAttr(`span meta[itemprop="name"]`, "content"))
			fmt.Printf("Talk Date: %s\n", strings.Split(e.ChildAttr(`meta[itemprop="uploadDate"]`, "content"), "T")[0])
		*/
		tags := make([]string, 0)
		detailCollector.Visit(e.ChildAttr(`link`, "href"))
		e.ForEach(`meta[property="og:video:tag"]`, func(_ int, el *colly.HTMLElement) {
			link := el.Attr("content")
			if link == "" {
				return
			}
			tags = append(tags, link)
		})

		// unlucky can't get the exact point where is printed the number of views.
		views := e.ChildTexts(`span`)

		newTalk := tedTalk{
			Title:       e.ChildAttr(`meta[itemprop="name"]`, "content"),
			Description: e.ChildAttr(`meta[itemprop="description"]`, "content"),
			Views:       strings.Split(views[3], "\n")[0],
			//Thumb:       e.ChildAttr(`link[itemprop="thumbnailUrl"`, "content"), // Not Working...
			Author: e.ChildAttr(`span meta[itemprop="name"]`, "content"),
			Date:   strings.Split(e.ChildAttr(`meta[itemprop="uploadDate"]`, "content"), "T")[0],
			Tags:   tags,
			Link:   e.ChildAttr(`meta[property="og:url"]`, "content"),
		}
		allTalks = append(allTalks, newTalk)

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.ted.com/talks")

	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal("Error during CSV Creation", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	headers := allTalks.getHeaders()
	if err := writer.Write(headers); err != nil {
		log.Fatal("Errore durante scrittura di Headers", err)
	}
	for _, value := range allTalks {
		if err := writer.Write(value.exportCSV()); err != nil {

			log.Fatal("Errore durante scrittura dei valori", err)
		}
	}
}

func (talks *listTalk) getHeaders() []string {
	//return []string{`Title;Description;Thumb;Author;Date;Tags;Link`}
	return []string{`Title;Description;Views;Author;Date;Tags;Link`}
}

func (talks *tedTalk) exportCSV() []string {
	str := make([]string, 0)
	str = append(str, talks.Title, talks.Description, talks.Views, talks.Author, talks.Date, "[\""+strings.Join(talks.Tags, "\",\"")+"\"]", talks.Link)
	return str
}
