package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	fmt.Println("Starting!")

	resp, err := http.Get("https://www.example.com")
	// resp, err := http.Get("https://www.iana.org/help/example-domains")
	if err != nil {
		fmt.Printf("Error occured %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// fmt.Printf("Body: %s", body)

	doc, err := html.Parse(bytes.NewReader(body))
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					break
				}
			}
		}
	}

}
