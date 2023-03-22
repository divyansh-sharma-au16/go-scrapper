package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type BlogPost struct {
	url, name string
}

// it verifies if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	fmt.Println("hii everyone....")
	// initializing the slice of structs to store the data to scrape
	var blogs []BlogPost

	// initializing the list of pages to scrape with an empty slice
	var pagesToScrape []string

	// the first pagination URL to scrape
	pageToScrape := "https://www.arkoselabs.com/blog/page/1/"

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	// current iteration
	i := 1
	// max pages to scrape
	limit := 5

	// creating a new Colly instance
	c := colly.NewCollector()

	// iterating over the list of pagination links to implement the crawling logic
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		// discovering a new page
		newPaginationLink := e.Attr("href")

		// if the page discovered is new
		if !contains(pagesToScrape, newPaginationLink) {
			// if the page discovered should be scraped
			if !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})

	// scraping logic
	c.OnHTML("div.al-post-item", func(e *colly.HTMLElement) {
		blog := BlogPost{}

		blog.url = e.ChildAttr("a", "href")
		blog.name = e.ChildText(".auth-data")

		blogs = append(blogs, blog)
	})

	c.OnScraped(func(response *colly.Response) {
		// until there is still a page to scrape
		if len(pagesToScrape) != 0 && i < limit {
			// getting the current page to scrape and removing it from the list
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			// incrementing the iteration counter
			i++

			// visiting a new page
			c.Visit(pageToScrape)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// visiting the target page
	c.Visit(pageToScrape)

	// opening the CSV file
	file, err := os.Create("blogs.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// writing the CSV headers
	headers := []string{
		"url",
		"name",
	}
	writer.Write(headers)

	// writing each blog as a CSV row
	for _, blog := range blogs {
		// converting a blog to an array of strings
		record := []string{
			blog.url,
			blog.name,
		}

		// adding a CSV record to the output file
		writer.Write(record)
	}
	defer writer.Flush()

}
