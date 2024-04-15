package scrapes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

var searchUrl string = "https://www.thetrainline.com/buytickets/"
var iso8601Layout string = "2006-01-02T15:04:05Z0700"

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

func getStationByCode(code string) (string, error) {
	jsonFile, err := os.Open("../../../../data/station-list.json")

	if err != nil {
		return "", err
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var stations []Station

	json.Unmarshal(byteValue, &stations)

	for _, station := range stations {
		if station.Code == code {
			return strings.ReplaceAll(station.Name, "-", " "), nil
		}
	}

	return "", errors.New("Station not found")
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
				price64, _ := strconv.ParseFloat(priceVal, 32)
				result.Price = float32(price64)
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

	c.Post(searchUrl, form)

	return res, nil
}
