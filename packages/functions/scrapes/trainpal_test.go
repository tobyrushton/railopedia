package scrapes

import (
	"os"
	"testing"
)

func TestTrainpalSingle(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	res, err := ScrapeTrainpal(TestRequest)
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
		if !timeRegex.MatchString(r.DepartureTime) {
			t.Errorf("Invalid departure time: %s", r.DepartureTime)
		}
		if !timeRegex.MatchString(r.ArrivalTime) {
			t.Errorf("Invalid arrival time: %s", r.ArrivalTime)
		}
		if r.Price < 0 {
			t.Errorf("Invalid price: %f", r.Price)
		}
	}
}

func TestTrainpalReturn(t *testing.T) {
	skipCI(t)

	res, err := ScrapeTrainpal(TestRequestReturn)
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
		if !timeRegex.MatchString(r.DepartureTime) {
			t.Errorf("Invalid departure time: %s", r.DepartureTime)
		}
		if !timeRegex.MatchString(r.ArrivalTime) {
			t.Errorf("Invalid arrival time: %s", r.ArrivalTime)
		}
		if r.Price < 0 {
			t.Errorf("Invalid price: %f", r.Price)
		}
	}
}
