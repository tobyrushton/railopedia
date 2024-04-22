package scrapes

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/tobyrushton/railopedia/packages/functions/src/utils"
)

var trainlineUrl string = "https://www.thetrainline.com/buytickets/"

type Station struct {
	Name string `json:"name"`
	Id   string `json:"trainline_id"`
	Code string `json:"code"`
}

type Tickets struct {
	Tickets []map[string]any `json:"tickets"`
}

type StandardTickets struct {
	StandardTickets []Tickets `json:"standardTickets"`
}

type scrapedJson struct {
	FullJourneys []StandardTickets `json:"fullJourneys"`
}

func ScrapeTrainline(req Request) (ScrapeResultNonConditional, error) {
	out, err := time.Parse(iso8601Layout, req.Departure)
	if err != nil {
		fmt.Print(err)
		return ScrapeResultNonConditional{}, errors.New("invalid date")
	}
	outDay := out.Day()
	outMonth := out.Month()
	outYear := out.Year()
	outHour := out.Hour()
	outMin := out.Minute()

	outStation, err := getStationByCode(req.Origin)
	if err != nil {
		return ScrapeResultNonConditional{}, err
	}

	inStation, err := getStationByCode(req.Destination)
	if err != nil {
		return ScrapeResultNonConditional{}, err
	}

	form := map[string]string{
		"OriginStation":             outStation,
		"DestinationStation":        inStation,
		"RouteRestriction":          "NULL",
		"ViaAvoidStation":           "",
		"outwardDate":               fmt.Sprintf("%d-%d-%d", outYear, outMonth, outDay),
		"OutwardLeaveAfterOrBefore": "A",
		"OutwardHour":               fmt.Sprintf("%d", outHour),
		"OutwardMinute":             fmt.Sprintf("%d", outMin),
		"returnDate":                "",
		"InwardLeaveAfterOrBefore":  "A",
		"ReturnHour":                "",
		"ReturnMinute":              "",
		"AdultsTravelling":          "1",
		"ChildrenTravelling":        "0",
		"ExtendedSearch":            "Get times & tickets",
	}

	if req.Return != "" {
		in, err := time.Parse(iso8601Layout, req.Return)
		if err != nil {
			return ScrapeResultNonConditional{}, errors.New("invalid date")
		}
		inDay := in.Day()
		inMonth := in.Month()
		inYear := in.Year()
		inHour := in.Hour()
		inMin := in.Minute()

		form["returnDate"] = fmt.Sprintf("%d-%d-%d", inYear, inMonth, inDay)
		form["ReturnHour"] = fmt.Sprintf("%d", inHour)
		form["ReturnMinute"] = fmt.Sprintf("%d", inMin)
	}

	c := colly.NewCollector()

	var res ScrapeResultNonConditional

	c.OnHTML("body", func(e *colly.HTMLElement) {
		data := e.ChildAttr("form", "data-defaults")

		var results scrapedJson
		json.Unmarshal([]byte(data), &results)
		res.Outbound = getTrainlinePrices(results, true)

		if req.Return != "" {
			res.Return = getTrainlinePrices(results, false)
		}
	})

	c.Post(trainlineUrl, form)

	return res, nil
}

func getTrainlinePrices(data scrapedJson, outbound bool) ScrapeResults {
	var index int
	if outbound {
		index = 0
	} else {
		index = 1
	}

	var res ScrapeResults

	for _, ticket := range data.FullJourneys[index].StandardTickets[0].Tickets {
		result := ScrapeResult{}
		if priceVal, ok := ticket["price"].(string); ok {
			result.Price = utils.PriceToFloat(priceVal)
		}
		if departureVal, ok := ticket["departureDateTime"].(string); ok {
			result.DepartureTime = formatTrainlineDate(departureVal)
		}
		if arrivalVal, ok := ticket["arrivalDateTime"].(string); ok {
			result.ArrivalTime = formatTrainlineDate(arrivalVal)
		}
		// TODO: Error handling if one is wrong
		// TODO: Add link
		// fmt.Println(ticket)
		res = append(res, result)
	}

	return res
}

// 2006-01-02 15:04:05

func formatTrainlineDate(date string) string {
	t, err := time.Parse("2006-01-02T15:04:05", date)
	if err != nil {
		return ""
	}
	return t.Format(iso8601Layout)
}
