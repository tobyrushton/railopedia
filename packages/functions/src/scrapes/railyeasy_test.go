package scrapes

import (
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/tobyrushton/railopedia/packages/functions/src/utils"
)

var now = time.Now().Add(100 * time.Hour)
var threeHoursLater = now.Add(3 * time.Hour)

var TestRequest = Request{
	Origin:      "SAC",
	Destination: "STP",
	Departure:   now.Format(iso8601Layout),
	Return:      "",
	Railcard:    "16-25",
}

var TestRequestReturn = Request{
	Origin:      "SAC",
	Destination: "STP",
	Departure:   now.Format(iso8601Layout),
	Return:      threeHoursLater.Format(iso8601Layout),
	Railcard:    "16-25",
}

var timeRegexString = "(?i)[0-9]+:[0-9]+"
var isoRegexString = "(?i)[0-9]+-[0-9]+-[0-9]+T[0-9]+:[0-9]+:[0-9]+"

var timeRegex = regexp.MustCompile(timeRegexString)
var isoRegex = regexp.MustCompile(isoRegexString)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}
}

func TestRaileasyReturn(t *testing.T) {
	skipCI(t)

	res, err := ScrapeRaileasy(TestRequestReturn)
	if err != nil {
		t.Error(err)
	}

	if len(res) == 0 {
		t.Error("No results")
	}

	for _, r := range res {
		if !isoRegex.MatchString(r.DepartureTime) {
			t.Errorf("Invalid departure time: %s", r.DepartureTime)
		}
		if !isoRegex.MatchString(r.ArrivalTime) {
			t.Errorf("Invalid arrival time: %s", r.ArrivalTime)
		}

		for time, price := range r.Price {
			time1, time2 := utils.SplitString(time, ",")
			if !isoRegex.MatchString(time1) {
				t.Errorf("Invalid time: %s", time)
			}
			if !isoRegex.MatchString(time2) {
				t.Errorf("Invalid time: %s", time)
			}
			if price < 0 {
				t.Errorf("Invalid price: %f", price)
			}
		}
	}
}

func TestRaileasySingle(t *testing.T) {
	skipCI(t)

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
