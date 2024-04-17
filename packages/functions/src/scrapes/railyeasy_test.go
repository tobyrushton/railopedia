package scrapes

import (
	"fmt"
	"testing"
)

func TestRaileasy(t *testing.T) {
	req := Request{
		Origin:      "SAC",
		Destination: "STP",
		Departure:   "2024-05-19T17:26:25Z",
		Return:      "",
	}

	res, err := ScrapeRaileasy(req)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
