package scrapes

import (
	"encoding/json"
	"os"
	"testing"
)

func TestScrape(t *testing.T) {
	skipCI(t)

	res := Scrape(TestRequestReturn)

	if len(res) == 0 {
		t.Error("No results")
	}

	jsonData, _ := json.Marshal(res)

	os.WriteFile("test.json", jsonData, 0644)
}
