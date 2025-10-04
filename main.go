package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"log"

	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Queue struct {
	size     int
	elements []string
	mu       sync.Mutex
}

func (q *Queue) getSize() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.size
}

func (q *Queue) enqueue(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.elements = append(q.elements, url)
	q.size++
}

func (q *Queue) dequeue() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.elements) == 0 {
		return "", false
	}
	url := q.elements[0]
	q.elements = q.elements[1:]
	return url, true
}

type Crawled struct {
	data map[string]bool
	size int
	mu   sync.Mutex
}

func (c *Crawled) getCrawledSetSize() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.size
}

func (c *Crawled) add(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[url] = true
	c.size++
}

func (c *Crawled) has(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[url]
}

type DbConnection struct {
	db       string
	fileName string
	conn     *sql.DB
}

func (db *DbConnection) connect() {
	db.db = "sqlite3"
	db.fileName = "./webpages.db"
	connection, err := sql.Open(db.db, db.fileName)
	db.conn = connection
	if err != nil {
		log.Fatal(err)
	}
	initialSqlStatement := `create table if not exists webpages (id integer not null primary key, url TEXT, content TEXT);`
	_, err = db.conn.Exec(initialSqlStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, initialSqlStatement)
		return
	}
}

func (db *DbConnection) disconnect() {
	db.conn.Close()
}

func (db *DbConnection) addWebPageToDb(url string, content string) {
	_, err := db.conn.Exec(
		"INSERT INTO webpages (url, content) VALUES (?, ?)",
		url, content,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type Webpage struct {
	Url     string
	Title   string
	Content string
}

func parseHTML(db *DbConnection, url string, q *Queue, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <-sem }()

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

	var bodyText string

	// Walk the DOM
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.Body {
			bodyText = collectText(n)
		}

		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.HasPrefix(a.Val, "https://") {
					q.enqueue(a.Val)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	fmt.Printf("Added %s \n", url)
	db.addWebPageToDb(url, bodyText)
}

func collectText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data) + " "
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(collectText(c))
	}
	return sb.String()
}

func main() {
	fmt.Println("Starting!")

	db := DbConnection{db: "sqlite3", fileName: "./webpages.db", conn: nil}
	db.connect()

	queue := Queue{size: 0, elements: make([]string, 0)}
	crawled := Crawled{data: make(map[string]bool), size: 0}

	const initialUrl = "https://www.theshubhagrwl.in/"

	queue.enqueue(initialUrl)

	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	//bfs
	count := 0
	for count < 10 {
		url, ok := queue.dequeue()
		if !ok {
			if len(sem) == 0 {
				break
			}
			// if some are still active, wait briefly and check again
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if crawled.has(url) {
			continue
		}
		count++
		sem <- struct{}{}
		wg.Add(1)
		go parseHTML(&db, url, &queue, sem, &wg)
		crawled.add(url)
	}

	wg.Wait()
	// fmt.Println(crawled)

	db.disconnect()
}
