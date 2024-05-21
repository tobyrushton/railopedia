package scrapes

import (
	"testing"
)

func TestTrainlineReturn(t *testing.T) {
	res, err := ScrapeTrainline(TestRequestReturn)

	if err != nil {
		t.Error(err)
	}

	if len(res.Outbound) == 0 {
		t.Error("No results")
	}
	if len(res.Return) == 0 {
		t.Error("No results")
	}

	vals := make(ScrapeResults, 0)
	vals = append(vals, res.Outbound...)
	vals = append(vals, res.Return...)

	for _, r := range vals {
		if !isoRegex.MatchString(r.DepartureTime) {
			t.Errorf("Invalid departure time: %s", r.DepartureTime)
		}
		if !isoRegex.MatchString(r.ArrivalTime) {
			t.Errorf("Invalid arrival time: %s", r.ArrivalTime)
		}
		if r.Price < 0 {
			t.Errorf("Invalid price: %f", r.Price)
		}
	}
}

func TestTrainlineSingle(t *testing.T) {
	res, err := ScrapeTrainline(TestRequest)
	if err != nil {
		t.Error(err)
	}

	if len(res.Outbound) == 0 {
		t.Error("No results")
	}
	if len(res.Return) != 0 {
		t.Error("Results returned for single journey")
	}

	for _, r := range res.Outbound {
		if !isoRegex.MatchString(r.DepartureTime) {
			t.Errorf("Invalid departure time: %s", r.DepartureTime)
		}
		if !isoRegex.MatchString(r.ArrivalTime) {
			t.Errorf("Invalid arrival time: %s", r.ArrivalTime)
		}
		if r.Price < 0 {
			t.Errorf("Invalid price: %f", r.Price)
		}
	}
}
