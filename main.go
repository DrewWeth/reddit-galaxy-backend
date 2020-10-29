package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

func readCsv() [][]string {
	// file, err := os.Open("./subreddits-1000.csv")
	file, err := os.Open("./subreddits.csv")
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
	Link  string  `json:"link"`
}

// Match responses
type Response struct {
	Input string  `json:"input"`
	Index int     `json:"index"`
	Res   []Match `json:"res"`
}

func handler(w http.ResponseWriter, r *http.Request, records [][]string, rows map[string]int, cols map[int]string) {
	path := r.URL.Path[5:]

	rowThatMatches := rows[path]
	if rowThatMatches == 0 {
		json.NewEncoder(w).Encode(nil)
		return
	}
	subRedditMatches := records[rowThatMatches]

	matchList := make([]Match, len(subRedditMatches))
	for i, v := range subRedditMatches {
		converted, _ := strconv.ParseFloat(v, 64)
		matchList[i] = Match{
			Name:  cols[i],
			Value: converted,
			Link:  fmt.Sprintf("http://localhost:8080/sub/%s", cols[i]),
		}
	}

	sort.Slice(matchList, func(i, j int) bool { return matchList[i].Value > matchList[j].Value })

	_res := Response{
		Input: path,
		Index: rowThatMatches,
		Res:   matchList,
	}
	json.NewEncoder(w).Encode(_res)
}

func handlerIndex(w http.ResponseWriter, r *http.Request, records [][]string) {
	fmt.Fprintf(w, "%s", "Oh, Hello")
}

var cols map[int]string
var rows map[string]int

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	records := readCsv()
	fmt.Println("Done parsing")

	cols = make(map[int]string)
	rows = make(map[string]int)

	// Find index for sub name
	for index, elements := range records {
		rows[elements[0]] = index
	}

	// Find name of column from index
	for index, header := range records[0] {
		cols[index] = header

	}

	http.HandleFunc("/sample/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		json.NewEncoder(w).Encode(records[:1])
	})

	http.HandleFunc("/rows/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		json.NewEncoder(w).Encode(rows)
	})

	http.HandleFunc("/cols/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		json.NewEncoder(w).Encode(cols)
	})

	http.HandleFunc("/sub/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		handler(w, r, records, rows, cols)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		handlerIndex(w, r, records)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
