package scrapes

import (
	"fmt"
	"testing"
)

func TestTrainline(t *testing.T) {
	req := Request{
		Origin:      "SAC",
		Destination: "STP",
		Departure:   "2024-04-20T12:26:25Z",
		Return:      "",
	}

	res, err := ScrapeTrainline(req)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
