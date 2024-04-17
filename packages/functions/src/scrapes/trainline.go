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

func ScrapeTrainline(req Request) (ScrapeResults, error) {
	out, err := time.Parse(iso8601Layout, req.Departure)
	if err != nil {
		fmt.Print(err)
		return nil, errors.New("invalid date")
	}
	outDay := out.Day()
	outMonth := out.Month()
	outYear := out.Year()
	outHour := out.Hour()
	outMin := out.Minute()

	// in, err := time.Parse(iso8601Layout, req.Return)
	// if err != nil {
	// 	return nil, errors.New("Invalid date")
	// }

	outStation, err := getStationByCode(req.Origin)
	if err != nil {
		return nil, err
	}

	inStation, err := getStationByCode(req.Destination)
	if err != nil {
		return nil, err
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

	c := colly.NewCollector()

	var res ScrapeResults

	c.OnHTML("body", func(e *colly.HTMLElement) {
		data := e.ChildAttr("form", "data-defaults")

		var results scrapedJson
		json.Unmarshal([]byte(data), &results)

		for _, ticket := range results.FullJourneys[0].StandardTickets[0].Tickets {
			result := ScrapeResult{}
			if priceVal, ok := ticket["price"].(string); ok {
				result.Price = utils.PriceToFloat(priceVal)
			}
			if departureVal, ok := ticket["departureDateTime"].(string); ok {
				result.DepartureTime = departureVal
			}
			if arrivalVal, ok := ticket["arrivalDateTime"].(string); ok {
				result.ArrivalTime = arrivalVal
			}
			// TODO: Error handling if one is wrong
			// TODO: Add link
			// fmt.Println(ticket)
			res = append(res, result)
		}
	})

	c.Post(trainlineUrl, form)

	return res, nil
}
