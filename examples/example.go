// Alexa Scraper — Scrapeless LLM Chat Scraper (Go example)
//
// Docs:  https://docs.scrapeless.com/en/llm-chat-scraper/quickstart/introduction/
// Token: https://app.scrapeless.com/passport/login?redirect=/quick-start
//
// Run:
//
//	export SCRAPELESS_API_TOKEN="your_api_token"
//	go run example.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const apiURL = "https://api.scrapeless.com/api/v2/scraper/execute"

func main() {
	apiToken := os.Getenv("SCRAPELESS_API_TOKEN")
	if apiToken == "" {
		apiToken = "YOUR_API_TOKEN"
	}

	payload := map[string]any{
		"actor": "scraper.alexa",
		"input": map[string]any{
			"prompt":  "Recommended attractions in New York",
			"country": "US",
		},
		// Optional: receive the result via webhook instead of the sync response.
		// "webhook": map[string]any{"url": "https://www.your-webhook.com"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-token", apiToken)

	client := &http.Client{Timeout: 180 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= 300 {
		panic(fmt.Sprintf("request failed: %d %s", resp.StatusCode, raw))
	}

	var data struct {
		Status     string `json:"status"`
		TaskID     string `json:"task_id"`
		TaskResult struct {
			UserText   string `json:"user_text"`
			MdText     string `json:"md_text"`
			RawText    string `json:"raw_text"`
			References []struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"references"`
			Products []struct {
				Title string `json:"title"`
				Price string `json:"price"`
				URL   string `json:"url"`
			} `json:"products"`
		} `json:"task_result"`
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		panic(err)
	}

	answer := data.TaskResult.RawText
	if answer == "" {
		answer = data.TaskResult.MdText
	}

	fmt.Println("Status: ", data.Status)
	fmt.Println("Task ID:", data.TaskID)
	fmt.Println("\nAnswer:\n", answer)

	for _, ref := range data.TaskResult.References {
		fmt.Printf("- %s -> %s\n", ref.Title, ref.URL)
	}

	for _, product := range data.TaskResult.Products {
		fmt.Printf("* %s (%s) -> %s\n", product.Title, product.Price, product.URL)
	}

	var pretty bytes.Buffer
	_ = json.Indent(&pretty, raw, "", "  ")
	fmt.Println("\nRaw response:\n", pretty.String())
}
