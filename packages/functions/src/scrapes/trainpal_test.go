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

	if len(res) == 0 {
		t.Error("No results")
	}

	for _, r := range res {
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

	if len(res) == 0 {
		t.Error("No results")
	}

	for _, r := range res {
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
