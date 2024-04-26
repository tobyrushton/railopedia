package scrapes

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/tobyrushton/railopedia/packages/functions/src/utils"
)

var raileasyUrl string = "https://new.raileasy.co.uk/"

func ScrapeRaileasy(req Request) (ScrapeResultsConditional, error) {
	// open browser
	page := rod.New().MustConnect().MustPage(raileasyUrl)
	defer page.MustClose()

	// input stations
	page.MustElement("#station-autocomplete-from").MustInput(req.Origin)
	page.MustElement("#station-autocomplete-to").MustInput(req.Destination)

	//enter outbound date and time
	out, err := utils.GetTime(req.Departure)
	if err != nil {
		return nil, errors.New("invalid date")
	}
	selectRaileasyDate(page, out, true)
	if req.Return == "" {
		setRaileasySingle(page)
	} else {
		//enter return date and time
		in, err := utils.GetTime(req.Return)
		if err != nil {
			return nil, errors.New("invalid date")
		}
		selectRaileasyDate(page, in, false)
	}

	wait := page.MustWaitRequestIdle()
	// submit form and wait for the results
	page.MustElement("#cookie-banner-accept").MustClick()
	page.MustElement("#search-button").MustClick()
	wait()

	// get outbound journeys
	outboundJourneys := page.MustElement("div.grid").MustElement("div.grid").MustElements("div")[0].MustElements(`div[tabindex="0"]`)
	results := make(ScrapeResultsConditional, 0)

	removeLaterAndEarlier(&outboundJourneys)
	for _, journey := range outboundJourneys {
		if req.Return != "" {
			in, err := utils.GetTime(req.Return)
			if err != nil {
				return nil, errors.New("invalid date")
			}

			results = append(results, getRaileasyPrice(page, journey, out, in))
		} else {
			price := getRaileasyBasePrice(journey)
			departTime, arrivalTime := getRaileasyJourneyTimes(journey)
			departTimeISO := utils.HourStringToISO(departTime, out)
			arrivalTimeISO := utils.HourStringToISO(arrivalTime, out)
			key := departTimeISO + "," + arrivalTimeISO
			results = append(results, ScrapeResultConditional{
				DepartureTime: departTime,
				ArrivalTime:   arrivalTime,
				Price:         map[string]float32{key: price},
			})
		}
	}

	return results, nil
}

func setRaileasySingle(page *rod.Page) {
	journeyType := page.MustElement("#journeyType")
	journeyType.MustElementR("button", "Single").MustClick()
}

func selectRaileasyDate(page *rod.Page, date time.Time, single bool) {
	// activate date picker
	if single {
		page.MustElement(`input[aria-label="Outbound date"]`).MustClick()
	} else {
		page.MustElement(`input[aria-label="Return date"]`).MustClick()
	}

	// select correct month
	monthString := date.Month().String()[0:3]
	currMonth := page.MustElement(`button[aria-label="Open months overlay"]`)

	for monthString != currMonth.MustText() {
		page.MustElement(`button[aria-label="Next month"]`).MustClick()
	}

	// select correct day
	day := date.Day()
	month := date.Month()
	year := date.Year()
	page.MustElement(fmt.Sprintf(`[id="%d-%02d-%d"]`, year, month, day)).MustClick()

	selectRaileasyTime(page, date, single)

}

func selectRaileasyTime(page *rod.Page, date time.Time, outbound bool) {
	// gets the correct time and hours ids
	hourId := "#"
	minuteId := "#"
	if outbound {
		hourId += "outbound-depart-hour"
		minuteId += "outbound-depart-minute"
	} else {
		hourId += "return-depart-hour"
		minuteId += "return-depart-minute"
	}
	// select the correct hour
	hour := date.Hour()
	minute := utils.RoundToNextFifteen(date.Minute())
	if minute == 60 {
		minute = 00
		if hour == 23 {
			hour = 00
		} else {
			hour++
		}
	}
	page.MustElement(hourId).MustSelect(fmt.Sprintf("%d", hour))

	// select the correct minute
	page.MustElement(minuteId).MustSelect(fmt.Sprintf("%d", minute))
}

func getRaileasyBasePrice(journey *rod.Element) float32 {
	// get the base price
	basePriceRaw := journey.MustElement("p.block").MustElement("span").MustText()
	if strings.Contains(basePriceRaw, "from") {
		_, basePriceRaw = utils.SplitString(basePriceRaw, " ")
	}
	basePrice := utils.PriceToFloat(utils.SanitisePrice(basePriceRaw))

	return basePrice
}

func getRaileasyJourneyTimes(journey *rod.Element) (string, string) {
	timeRaw := journey.MustElement("#journey-time").MustText()
	departureTime, arrivalTime := utils.SplitString(timeRaw, " -> ")

	return departureTime, arrivalTime
}

func getRaileasyPrice(page *rod.Page, journey *rod.Element, out time.Time, in time.Time) ScrapeResultConditional {
	// get journey times
	departureTime, arrivalTime := getRaileasyJourneyTimes(journey)
	departureTimeISO := utils.HourStringToISO(departureTime, out)
	arrivalTimeISO := utils.HourStringToISO(arrivalTime, out)

	basePrice := getRaileasyBasePrice(journey)

	// get journey prices
	journey.MustClick()

	// get the prices
	returnJourneys := page.MustElement("div.grid").MustElement("div.grid").MustElements("div.overflow-hidden")[1].MustElements(`div[tabindex="0"]`)
	removeLaterAndEarlier(&returnJourneys)

	price := make(map[string]float32)

	for _, returnJourney := range returnJourneys {
		returnDepartTime, returnArrivalTime := getRaileasyJourneyTimes(returnJourney)
		key := utils.HourStringToISO(returnDepartTime, in) + "," + utils.HourStringToISO(returnArrivalTime, in)

		if returnJourney.MustHas("p.block") {
			// get the price
			priceRaw := returnJourney.MustElement("p.block").MustElement("span").MustText()
			priceFloat := utils.PriceToFloat(utils.SanitisePrice(priceRaw)) + basePrice

			price[key] = priceFloat

		} else { // is the cheapest option
			price[key] = basePrice
		}

	}

	return ScrapeResultConditional{
		DepartureTime: departureTimeISO,
		ArrivalTime:   arrivalTimeISO,
		Price:         price,
	}
}

// removes the later and earlier buttons from list of journeys
func removeLaterAndEarlier(list *rod.Elements) {
	*list = (*list)[1 : len(*list)-1]
}
