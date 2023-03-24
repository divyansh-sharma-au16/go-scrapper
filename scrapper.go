package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type BlogPost struct {
	url  string
	name string
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func scrapePage(pageToScrape string, wg *sync.WaitGroup, blogChan chan<- []BlogPost) {
	defer wg.Done()

	c := colly.NewCollector(
		colly.Async(true),
	)

	var blogs []BlogPost
	c.OnHTML("div.al-post-item", func(e *colly.HTMLElement) {
		blog := BlogPost{}
		blog.url = e.ChildAttr("a", "href")
		blog.name = e.ChildText(".auth-data")
		blogs = append(blogs, blog)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(pageToScrape)

	c.Wait()

	blogChan <- blogs
}

func main() {
	start := time.Now()
	fmt.Println("Starting the scraper...")

	var blogs []BlogPost
	var pagesToScrape []string
	pagesToScrape = append(pagesToScrape, "https://www.arkoselabs.com/blog/page/1/")

	var pagesDiscovered []string
	pagesDiscovered = append(pagesDiscovered, pagesToScrape...)

	var wg sync.WaitGroup
	blogChan := make(chan []BlogPost)

	for len(pagesToScrape) > 0 {
		pageToScrape := pagesToScrape[0]
		pagesToScrape = pagesToScrape[1:]

		wg.Add(1)
		go scrapePage(pageToScrape, &wg, blogChan)

		c := colly.NewCollector(
			colly.Async(true),
		)
		c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
			newPaginationLink := e.Attr("href")
			if !contains(pagesToScrape, newPaginationLink) && !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
				pagesDiscovered = append(pagesDiscovered, newPaginationLink)
			}
		})
		c.Visit(pageToScrape)
		c.Wait()
	}

	go func() {
		wg.Wait()
		close(blogChan)
	}()

	for blogList := range blogChan {
		blogs = append(blogs, blogList...)
	}

	file, err := os.Create("blogs.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"url",
		"name",
	}
	writer.Write(headers)

	for _, blog := range blogs {
		record := []string{
			blog.url,
			blog.name,
		}
		writer.Write(record)
	}
	writer.Flush()

	fmt.Println("Scraping completed. Check blogs.csv for output.")

	elapsed := time.Since(start)
	log.Printf("time taken to complete %s", elapsed)
}
