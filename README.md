# Web Crawler in Go

## Steps so far
* Created a basic function that gets the html doc of a webpage
* used golang.org/x/net/html to parse the html page
* filtered out all the URLs
* -
* Stored the URLS in a queue
* Stored only the absolute URLs and not the dynamic ones
* Separated get call in a method
* Implemented BFS