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

func getLinks(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error occured %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading body : %v", err)
	}
	links := make([]string, 0)
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Error parsing body: %v", err)
	}
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				//storing only the absolute urls
				if a.Key == "href" && strings.HasPrefix(a.Val, "https://") {
					links = append(links, a.Val)
					// break
				}
			}
		}
	}
	return links
}

func main() {
	fmt.Println("Starting!")

	links := getLinks("https://www.example.com/")

	//bfs
	visited := make(map[string]bool)
	count := 0
	for len(links) > 0 {
		curr := links[0]
		links = links[1:]

		if !visited[curr] {
			newLinks := getLinks(curr)
			fmt.Printf("accessed: %v \n", curr)
			visited[curr] = true
			links = append(links, newLinks...)
			count++
		}

		if count == 10 {
			fmt.Println(visited)
			break
		}
	}

}
