package scrapes

import (
	"encoding/json"
	"os"
	"testing"
)

func TestScrapeReturn(t *testing.T) {
	skipCI(t)

	res, _, _ := Scrape(TestRequestReturn)

	if len(res) == 0 {
		t.Error("No results")
	}

	jsonData, _ := json.Marshal(res)

	os.WriteFile("test.json", jsonData, 0644)
}

func TestScrapeSingle(t *testing.T) {
	skipCI(t)

	_, res, _ := Scrape(TestRequest)

	if len(res) == 0 {
		t.Error("No results")
	}

	jsonData, _ := json.Marshal(res)

	os.WriteFile("test.json", jsonData, 0644)
}
