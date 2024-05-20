package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Category uint32

const (
	Category_User Category = iota
	Category_Article
)

type User struct {
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	Category Category `json:"category"`
}

type Article struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Author   string   `json:"author"`
	Category Category `json:"category"`
}

const IndexName = "mix_index"

func main() {

	var info map[string]interface{}

	cfg := elastic.Config{
		Addresses: []string{"http://elasticsearch-elasticsearch-1:9200"},
	}
	es, err := elastic.NewClient(cfg)
	if err != nil {
		log.Fatalln("can't connect to elasticsearch")
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elastic.Version)
	log.Printf("Server: %s", info["version"].(map[string]interface{})["number"])

	// Check response status
	// Deserialize the response into a map.
	// Check response status
	// Deserialize the response into a map.
	// Print client and server version numbers.
	// reBuildIndex(es)
	// indexUser(es)
	// indexArticle(es)

	query := `{
  "query": {
     "term": {
      "category": {
        "value": "0",
        "boost": 1.0
      }
    }
  }
}`
	res, err = es.Search(
		es.Search.WithIndex(IndexName),
		es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		log.Fatalf("Error create index response: %s", err)
	}

	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	if res.IsError() {
		log.Printf("[%s] Error search document", res.Status())
	} else {
		if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {

			log.Printf("[%s] %s", res.Status(), info["hits"].(map[string]interface{})["hits"])
		}
	}
}

func reBuildIndex(es *elastic.Client) {
	var info map[string]interface{}

	deleteRequest := esapi.IndicesDeleteRequest{
		Index: []string{IndexName},
	}

	res, err := deleteRequest.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error create index response: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	mapping := `{
		"mappings": {
			"properties": {
				"category": {
					"type": "short"
				}
			}
		}
	}`

	request := esapi.IndicesCreateRequest{
		Index: IndexName,
		Body:  strings.NewReader(mapping),
	}

	res, err = request.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error create index response: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	log.Printf("Client: %s", elastic.Version)
	log.Printf("Server: %s", info["version"].(map[string]interface{})["number"])
}

func indexUser(es *elastic.Client) {
	var info map[string]interface{}

	user := User{
		Name:     "John Doe",
		Age:      30,
		Category: Category_User,
	}
	data, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Error marshaling user: %s", err)
	}

	var id = 1
	req := esapi.IndexRequest{
		Index:      IndexName,
		DocumentID: strconv.Itoa(id),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), id)
	} else {
		if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {

			log.Printf("[%s] %s; version=%d", res.Status(), info["result"], int(info["_version"].(float64)))
		}
	}
}

func indexArticle(es *elastic.Client) {
	var info map[string]interface{}

	article := Article{
		Title:    "Golang",
		Content:  "Golang is good",
		Author:   "Kindy Wu",
		Category: Category_Article,
	}
	data, err := json.Marshal(article)
	if err != nil {
		log.Fatalf("Error marshaling article: %s", article)
	}

	var id = 2
	req := esapi.IndexRequest{
		Index:      IndexName,
		DocumentID: strconv.Itoa(id),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), id)
	} else {
		if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {

			log.Printf("[%s] %s; version=%d", res.Status(), info["result"], int(info["_version"].(float64)))
		}
	}
}
