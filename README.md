# Web Crawler in Go

# Dataflow of the web crawler:
<image src="./images/data-flow.png"/>

# Learnings
My main aim behind building this web crawler was to familiarize myself with the concepts of web crawling and also the implementation of golang concurrency.

Via this project I was able to learn new concepts of web crawling like parsing, queuing, politeness. I have also implemented Data structures Queue and Breadth First Search. In go lang I got to learn about Channels, Wait Groups, Semaphores. 

This crawler is not perfect, I have tested it for 1000 sites but there are a lot of improvements that can be made. First is namely the politeness. Currently this crawler stores all the content of the body which can be optimized. Also I can improve on the error handling. Also the feature I want to add is the support for robots.txt

Short summary
## Pros
* Concurrently fetches and parses web pages
* Database inserts are concurrent
* BFS is efficient 
* Performance is decent

## Cons
* Storing the entire body might also include script tags and some css
* Relative links are ignored
* Politeness can be improved