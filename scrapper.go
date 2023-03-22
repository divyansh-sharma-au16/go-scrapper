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
 
func main() { 
    fmt.Println("hii everyone....")
    // initializing the slice of structs to store the data to scrape 
    var blogs []BlogPost
 
    // creating a new Colly instance 
    c := colly.NewCollector() 
 
 
    // scraping logic 
    c.OnHTML("div.al-post-item", func(e *colly.HTMLElement) { 
        blog := BlogPost{} 
 
        blog.url = e.ChildAttr("a", "href")  
        blog.name = e.ChildText(".auth-data") 
 
        blogs = append(blogs, blog) 
        fmt.Println("hi there/./././../.")
    }) 

	c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

	// visiting the target page 
    c.Visit("https://www.arkoselabs.com/blog/") 

 
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
