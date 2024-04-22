package scrapes

import (
	"testing"
)

func TestTrainticketsSingle(t *testing.T) {
	skipCI(t)

	res, err := ScrapeTraintickets(TestRequest)

	if err != nil {
		t.Error(err)
	}

	if len(res) == 0 {
		t.Error("No results returned")
	}

	for _, r := range res {
		if !timeRegex.MatchString(r.DepartureTime) {
			t.Errorf("Invalid departure time: %s", r.DepartureTime)
		}
		if !timeRegex.MatchString(r.ArrivalTime) {
			t.Errorf("Invalid arrival time: %s", r.ArrivalTime)
		}

		if len(r.Price) != 1 {
			t.Errorf("Invalid price: %v", r.Price)
		}

		for time, price := range r.Price {
			if !isoRegex.MatchString(time) {
				t.Errorf("Invalid time: %s", time)
			}
			if price < 0 {
				t.Errorf("Invalid price: %f", price)
			}
		}
	}
}

func TestTrainticketsReturn(t *testing.T) {
	skipCI(t)

	res, err := ScrapeTraintickets(TestRequestReturn)
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

		for time, price := range r.Price {
			if !isoRegex.MatchString(time) {
				t.Errorf("Invalid time: %s", time)
			}
			if price < 0 {
				t.Errorf("Invalid price: %f", price)
			}
		}
	}
}
