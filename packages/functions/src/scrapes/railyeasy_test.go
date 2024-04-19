package scrapes

import (
	"regexp"
	"testing"
	"time"
)

var now = time.Now().Add(1 * time.Hour)
var threeHoursLater = now.Add(3 * time.Hour)

var TestRequest = Request{
	Origin:      "SAC",
	Destination: "STP",
	Departure:   now.Format(iso8601Layout),
	Return:      "",
}

var TestRequestReturn = Request{
	Origin:      "SAC",
	Destination: "STP",
	Departure:   now.Format(iso8601Layout),
	Return:      threeHoursLater.Format(iso8601Layout),
}

var timeRegexString = "(?i)[0-9]+:[0-9]+"
var isoRegexString = "(?i)[0-9]+-[0-9]+-[0-9]+T[0-9]{2}:[0-9]{2}:[0-9]{2}(\\.[0-9]{1,3})?Z"

var timeRegex = regexp.MustCompile(timeRegexString)
var isoRegex = regexp.MustCompile(isoRegexString)

func TestRaileasyReturn(t *testing.T) {
	res, err := ScrapeRaileasy(TestRequestReturn)
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

func TestRaileasySingle(t *testing.T) {
	res, err := ScrapeRaileasy(TestRequest)
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
