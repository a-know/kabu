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

type outlineItems struct {
	OutlineItems []*outlineItem `json:"items"`
}

type item struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

type outlineItem struct {
	Title        string `json:"title"`
	FormattedURL string `json:"formattedUrl"`
}

func main() {
	apiKey := os.Getenv("GOOGLE_APIKEY")
	cseID := os.Getenv("GOOGLE_CSE_ID")

	var companyName string
	if len(os.Args) < 2 {
		fmt.Println("引数が省略されたため、クリップボードの内容を用いて判定します。")
		text, err := clipboard.ReadAll()
		if err != nil {
			log.Fatalf("Failed to read from clipboard : %v\n", err)
		} else {
			companyName = text
		}
	} else {
		companyName = os.Args[1]
	}

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

	// 会社概要URLの取得
	outlineResponse, err := http.Get("https://www.googleapis.com/customsearch/v1?key=" + apiKey + "&cx=" + cseID + "&q=" + companyName + "+株式会社+概要")
	if err != nil {
		log.Fatalf("Failed to search : %v\n", err)
	}

	var outlineResults outlineItems
	defer outlineResponse.Body.Close()
	if err := json.NewDecoder(outlineResponse.Body).Decode(&outlineResults); err != nil {
		log.Fatalf("Failed to decode search result : %v\n", err)
	}
	var outlineURL string
	for _, outlineItem := range outlineResults.OutlineItems {
		matched, err := regexp.MatchString("概要", outlineItem.Title)
		if err != nil {
			log.Fatalf("Failed to regexp.MatchString : %v\n", err)
			break
		}
		if matched {
			outlineURL = outlineItem.FormattedURL
			break
		}
	}
	if outlineURL != "" {
		fmt.Println(fmt.Sprintf("会社概要URL: %s", outlineURL))
		if err := clipboard.WriteAll(outlineURL); err != nil {
			log.Fatalf("Failed to copy to clipboard : %v\n", err)
		}
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
