package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/atotto/clipboard"
)

type items struct {
	Items []*item `json:"items"`
}

type item struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

func main() {
	apiKey := os.Getenv("GOOGLE_APIKEY")
	cseID := os.Getenv("GOOGLE_CSE_ID")
	companyName := os.Args[1]
	maekabuPattern := "株式会社" + companyName
	mae := 0
	ushirokabuPattern := companyName + "株式会社"
	ushiro := 0

	response, err := http.Get("https://www.googleapis.com/customsearch/v1?key=" + apiKey + "&cx=" + cseID + "&q=" + companyName + "+株式会社")
	if err != nil {
		log.Fatalf("Failed to search : %v\n", err)
	}

	var results items
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		log.Fatalf("Failed to decode search result : %v\n", err)
	}
	for _, result := range results.Items {
		check(maekabuPattern, result.Title, &mae)
		check(maekabuPattern, result.Snippet, &mae)
		check(ushirokabuPattern, result.Title, &ushiro)
		check(ushirokabuPattern, result.Snippet, &ushiro)
	}

	fmt.Printf("前株マッチ数:%d\n", mae)
	fmt.Printf("後株マッチ数:%d\n", ushiro)
	if mae > ushiro {
		fmt.Println("前株です！")
		fmt.Println(maekabuPattern)
		if err := clipboard.WriteAll(maekabuPattern); err != nil {
			log.Fatalf("Failed to copy to clipboard : %v\n", err)
		}
	} else if ushiro > mae {
		fmt.Println("後株です！")
		fmt.Println(ushirokabuPattern)
		if err := clipboard.WriteAll(ushirokabuPattern); err != nil {
			log.Fatalf("Failed to copy to clipboard : %v\n", err)
		}
	} else {
		fmt.Println("わかりません！")
	}
}

func check(pattern string, str string, cnt *int) {
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		log.Fatalf("Failed to regexp.MatchString : %v\n", err)
	}
	if matched {
		*cnt++
	}
}
