package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"


)

func readCsv() [][]string {
	file, err := os.Open("./subreddits.csv")
	// file, err := os.Open("./subreddits-1000.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	return records
}

// Match value with name
type Match struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// Match responses
type Response struct {
	Input string  `json:"input"`
	Res   []Match `json:"res"`
}

const SUBREDDIT_FILENAME = "subreddit-related-build.json"


func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	lookup := useOrCreateLookupFile()

	http.HandleFunc("/subv2/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		path := r.URL.Path[7:]
		json.NewEncoder(w).Encode(Response{
			Input: path,
			Res: lookup[path],
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		fmt.Fprintf(w, "%s", "Oh, Hello")
	})

	fmt.Println("Running on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
