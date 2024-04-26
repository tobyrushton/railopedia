package scrapes

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/tobyrushton/railopedia/packages/functions/src/utils"
)

var trainticketsUrl = "https://www.traintickets.com/?/"

func ScrapeTraintickets(req Request) (ScrapeResultsConditional, error) {
	//open browser
	page := rod.New().MustConnect().MustPage(trainticketsUrl)
	defer page.MustClose()

	// input stations
	setTrainticketsStation(page, req.Origin, "#origin")
	setTrainticketsStation(page, req.Destination, "#destination")

	// select the dates
	err := setTrainticketsDate(page, req.Departure, true)
	if err != nil {
		return nil, errors.New("invalid date")
	}
	if req.Return != "" {
		err := setTrainticketsDate(page, req.Return, false)
		if err != nil {
			return nil, errors.New("invalid date")
		}
	}

	// submit form
	page.MustElement("#searchButton").MustClick()

	// get outbound journeys
	page.MustElementR("h3", "Choose") // waits for journeys to finish loading
	outboundJourneys := page.MustWaitLoad().MustElement("#outbound").MustElements("li.journey")
	res := make(ScrapeResultsConditional, 0)
	out, _ := time.Parse(iso8601Layout, req.Departure)

	if req.Return == "" {
		for _, journey := range outboundJourneys {
			price := getTrainticketsPrice(journey)
			departTime, arrivalTime := getTrainticketsJourneyTimes(journey)
			isoDepartTime := utils.HourStringToISO(departTime, out)
			isoArrivalTime := utils.HourStringToISO(arrivalTime, out)
			key := isoDepartTime + "," + isoArrivalTime
			res = append(res, ScrapeResultConditional{
				DepartureTime: isoDepartTime,
				ArrivalTime:   isoArrivalTime,
				Price:         map[string]float32{key: price},
			})

		}
	} else {
		in, _ := time.Parse(iso8601Layout, req.Return)
		for _, journey := range outboundJourneys {
			res = append(res, getTrainticketsJourneyPriceReturn(page, journey, out, in))
		}
	}

	return res, nil
}

func setTrainticketsDate(page *rod.Page, date string, outbound bool) error {
	// format date
	dayOfJourney, err := time.Parse(iso8601Layout, date)
	if err != nil {
		return errors.New("invalid date")
	}
	// get correct input
	inputSelector := "#departdate"
	if !outbound {
		inputSelector = "#returndate"
	}

	// activate date selection
	page.MustElement(inputSelector).MustClick()

	var dateBox *rod.Element
	if outbound {
		dateBox = page.MustElement("#depart").MustElement("div.bb-date")
	} else {
		dateBox = page.MustElement("#return").MustElement("div.bb-date")
	}

	// select correct month
	currMonth := dateBox.MustElement("div.bb-date__label")
	month := dayOfJourney.Month().String()
	for currMonth.MustText() != month {
		dateBox.MustElement("button.bb-date__next").MustClick()
	}

	// select correct day
	day := dayOfJourney.Day()
	dateBox.MustElement(fmt.Sprintf("button[data-bb-day='%d']", day)).MustClick()

	hour := dayOfJourney.Hour()
	minute := utils.RoundToNextFifteen(dayOfJourney.Minute())
	if minute == 60 {
		minute = 00
		if hour == 23 {
			hour = 00
		} else {
			hour++
		}
	}

	// select correct time
	dateBox.MustElement(`select[name="hour"]`).MustSelect(fmt.Sprintf("%d", hour))
	dateBox.MustElement(`select[name="min"]`).MustSelect(fmt.Sprintf("%d", minute))

	// close date picker
	dateBox.MustElement("button.btn-ok").MustClick()

	return nil
}

func getTrainticketsPrice(journey *rod.Element) float32 {
	priceString := journey.MustElement("span.cost").MustText()
	_, price := utils.SplitString(priceString, ":")
	return utils.PriceToFloat(utils.SanitisePrice(price[1:8]))
}

func getTrainticketsJourneyTimes(journey *rod.Element) (string, string) {
	times := journey.MustElements("span.cityTime")
	return times[0].MustText(), times[1].MustText()
}

func getTrainticketsJourneyPriceReturn(page *rod.Page, journey *rod.Element, out time.Time, in time.Time) ScrapeResultConditional {
	// select outbound journey
	journey.MustElement("button.go").MustClick()

	// get return journeys
	returnJourneys := page.MustElement("#inbound").MustElements("li.journey")

	price := make(map[string]float32)

	for _, returnJourney := range returnJourneys {
		journeyPrice := getTrainticketsPrice(returnJourney)
		departTime, arrivalTime := getTrainticketsJourneyTimes(returnJourney)
		departTimeISO := utils.HourStringToISO(departTime, in)
		arrivalTimeISO := utils.HourStringToISO(arrivalTime, in)
		key := departTimeISO + "," + arrivalTimeISO
		price[key] = journeyPrice
	}

	departTime, arrivalTime := getTrainticketsJourneyTimes(journey)
	departTimeISO := utils.HourStringToISO(departTime, out)
	arrivalTimeISO := utils.HourStringToISO(arrivalTime, out)

	page.MustElement("#change-outbound").MustClick()

	return ScrapeResultConditional{
		DepartureTime: departTimeISO,
		ArrivalTime:   arrivalTimeISO,
		Price:         price,
	}
}

func setTrainticketsStation(page *rod.Page, station string, selector string) {
	wait := page.MustWaitRequestIdle()
	page.MustElement(selector).MustClick()
	for _, s := range station {
		page.Keyboard.MustType(input.Key(s))
	}
	wait()
}
