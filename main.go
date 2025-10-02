package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"

	"log"

	_ "github.com/mattn/go-sqlite3"

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

func parseHTML(db *DbConnection, url string, q *Queue) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error occured %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	db.addWebPageToDb(url, string(body))

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

	db := DbConnection{db: "sqlite3", fileName: "./webpages.db", conn: nil}
	db.connect()

	queue := Queue{size: 0, elements: make([]string, 0)}
	crawled := Crawled{data: make(map[string]bool), size: 0}

	const initialUrl = "https://www.theshubhagrwl.in/"

	queue.enqueue(initialUrl)
	crawled.add(initialUrl)
	queue.dequeue()
	parseHTML(&db, initialUrl, &queue)

	//bfs
	count := 0
	for queue.getSize() > 0 && count < 10 {
		curr := queue.dequeue()
		if !crawled.has(curr) {
			parseHTML(&db, curr, &queue)
			fmt.Printf("visiting: %s\n", curr)
			crawled.add(curr)
			count++
		}
	}
	fmt.Println(crawled)

	db.disconnect()

}
