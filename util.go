package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"io/ioutil"
	"encoding/json"
)

func makeSubLookup() map[string][]Match {
	records := readCsv()
	getColumnName := make(map[int]string)
	getRowFromName := make(map[string]int)
	
	lookup := make(map[string][]Match)

	// Find index for sub name
	for index, elements := range records {
		getRowFromName[elements[0]] = index
	}
	// Find name of column from index
	for index, header := range records[0] {
		getColumnName[index] = header
	}
	
	listOfRelated := make([]Match, len(records) - 1)
	for rowI, row := range records{
		if(rowI != 0){
			for i, value := range row {
				if(i != 0){
					converted, _ := strconv.ParseFloat(value, 64)
					listOfRelated[i - 1] = Match{
						Name:  getColumnName[i],
						Value: converted,
					}
				}
			}	
			// Sort by value
			sort.Slice(listOfRelated, func(i, j int) bool { return listOfRelated[i].Value > listOfRelated[j].Value })
			// Make local copy 
			s := make([]Match, 5)
			copy(s[:], listOfRelated[:5])
			lookup[row[0]] = s
		}
	}
	return lookup
}

func readFile() map[string][]Match {
	jsonFile, err := os.Open(SUBREDDIT_FILENAME)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	fmt.Println("Successfully opened " + SUBREDDIT_FILENAME)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	
	var matchValues map[string][]Match
	json.Unmarshal(byteValue, &matchValues)
	fmt.Println("Successfully parsed " + SUBREDDIT_FILENAME)
	return matchValues
}

func createFile() {
	lookup := makeSubLookup()

	file, _ := json.MarshalIndent(lookup, "", " ")
	_ = ioutil.WriteFile(SUBREDDIT_FILENAME, file, 0644)
}

func useOrCreateLookupFile()map[string][]Match {
	if _, err := os.Stat(SUBREDDIT_FILENAME); os.IsNotExist(err) {
		fmt.Println("File not found, creating...")
		createFile()
	}
	return readFile()
}


