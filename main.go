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
	res, err := http.Get("https://api.github.com/search/repositories?q=stars:>=5000+language:c&sort=stars&order=desc")
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

	printData(data)
}

func printData(data JSONData) {
	log.Printf("Repositories found: %d", data.Count)
	const format = "%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)

	fmt.Fprintf(tw, format, "Repo", "Stars", "Created at", "Description")

	for _, i := range data.Items {
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