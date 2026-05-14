package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

//repo owner
type Owner struct {
	Login string
}

//repo template
type Item struct {
	ID int
	Name string
	FullName string `json:"full_name"`
	Owner Owner
	Description string
	CreatedAt string `json:"created_at"`
	StarCount int `json:"stargazers_count"`
}

type JSONData struct {
	Count int `json:"total_count"`
	Items []Item
}

func main() {
	var limit int

	fmt.Println("How many repository search results would you like to see?")
	fmt.Scanln(&limit)

	if limit <= 0 || limit > 100 {
		log.Fatal("Limit must be between 1 and 100.")
	}

	url := fmt.Sprintf("https://api.github.com/search/repositories?q=stars:>=5000+language:c&sort=stars&order=desc&per_page=%d", limit)
	
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// if res.StatusCode != 200 {
	// 	log.Fatal("Unexpected status code", res.StatusCode)
	// }

	if res.StatusCode != http.StatusOK {
		log.Fatal("Unexpected status code", res.StatusCode)
	}

	data := JSONData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	printData(data, limit)
}

func printData(data JSONData, limit int) {
	log.Printf("Repositories found: %d", data.Count)
	const format = "%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)

	fmt.Fprintf(tw, format, "Repo", "Stars", "Created at", "Description")

	if limit > len(data.Items) {
		limit = len(data.Items)
	}
	items := data.Items[:limit]

	for _, i := range items {
		desc := i.Description
		if len(desc) > 35 {
			desc = string(desc[:35]) + "..."
		}

		t, err := time.Parse(time.RFC3339, i.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(tw, format, i.FullName, i.StarCount, t.Year(), desc)
	}

	tw.Flush()
}