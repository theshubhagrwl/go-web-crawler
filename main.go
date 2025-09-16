package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Queue struct {
	size     int
	elements []string
}

func (q *Queue) getSize() int {
	return q.size
}

func (q *Queue) enqueue(url string) {
	q.elements = append(q.elements, url)
	q.size++
}

func (q *Queue) dequeue() string {
	url := q.elements[0]
	q.elements = q.elements[1:]
	q.size--
	return url
}

type Crawled struct {
	data map[string]bool
	size int
}

func (c *Crawled) getCrawledSetSize() int {
	return c.size
}

func (c *Crawled) add(url string) {
	c.data[url] = true
	c.size++
}
func (c *Crawled) has(url string) bool {
	return c.data[url]
}

func parseHTML(url string, q *Queue) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error occured %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading body : %v", err)
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Error parsing body: %v", err)
	}
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				//storing only the absolute urls
				if a.Key == "href" && strings.HasPrefix(a.Val, "https://") {
					q.enqueue(a.Val)
				}
			}
		}
	}
}

func main() {
	fmt.Println("Starting!")

	queue := Queue{size: 0, elements: make([]string, 0)}
	crawled := Crawled{data: make(map[string]bool), size: 0}

	const initialUrl = "https://www.theshubhagrwl.in/"
	// const initialUrl = "https://www.example.com/"
	queue.enqueue(initialUrl)
	crawled.add(initialUrl)
	queue.dequeue()
	parseHTML(initialUrl, &queue)

	//bfs
	count := 0
	for queue.getSize() > 0 && count < 10 {
		curr := queue.dequeue()
		if !crawled.has(curr) {
			parseHTML(curr, &queue)
			fmt.Printf("visiting: %s\n", curr)
			crawled.add(curr)
			count++
		}
	}
	fmt.Println(crawled)

}
