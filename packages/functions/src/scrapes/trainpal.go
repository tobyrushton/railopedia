package scrapes

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/tobyrushton/railopedia/packages/functions/src/utils"
)

var trainpalUrl string = "https://www.mytrainpal.com/"

func ScrapeTrainpal(req Request) (ScrapeResults, error) {
	out, err := time.Parse(iso8601Layout, req.Departure)
	if err != nil {
		return nil, errors.New("invalid date")
	}

	// open browser
	page := rod.New().MustConnect().MustPage(trainpalUrl)

	// input stations
	page.MustElement("#fromStation").MustInput(req.Origin)
	page.MustElement("div.el-station_cdf6f").MustClick()
	page.MustElement("#toStation").MustInput(req.Destination)
	time.Sleep(time.Second) // delay to allow dropdown to update
	page.MustElement("div.el-station_cdf6f").MustClick()

	// select the dates
	selectTrainpalDate(page, out, true)
	if req.Return != "" {
		in, err := time.Parse(iso8601Layout, req.Return)
		if err != nil {
			return nil, errors.New("invalid date")
		}
		selectTrainpalDate(page, in, false)
	}

	//submit form
	page.MustElement("button.search-btn_db7b7").MustClick()

	// gets journeys from page
	outboundJourneys := page.MustWaitStable().MustElement("div.left-inner_ac0c4").MustElements("div.journey-section_d201d")
	inboundJourneys := page.MustWaitStable().MustElement("div.right-inner_cf7d7").MustElements("div.journey-section_d201d")

	res := ScrapeResults{}

	for _, journey := range outboundJourneys {
		res = append(res, getTrainpalJourneyDetails(journey))
	}
	for _, journey := range inboundJourneys {
		res = append(res, getTrainpalJourneyDetails(journey))
	}

	return res, nil
}

func selectTrainpalDate(page *rod.Page, date time.Time, single bool) {
	// activate the date picker
	if !single {
		page.MustElement("div.add-return_df7cf").MustClick()
		inputList := page.MustElements(`input[placeholder="Date and time"]`)
		inputList[len(inputList)-1].MustClick()
	} else {
		page.MustElement(`input[placeholder="Date and time"]`).MustClick()
	}

	monthYearString := fmt.Sprintf("%s %d", date.Month().String(), date.Year())

	// select the element displaying `month year`
	currMonthYear := page.MustElement("div.date-head_f1c3f")
	currMonthYearString, _ := currMonthYear.Text()

	// continue until on correct month
	for !strings.Contains(currMonthYearString, monthYearString) {
		page.MustElement("div.right-btn_fca5f").MustClick()
		currMonthYearString, _ = currMonthYear.Text()
	}

	// click correct day
	page.MustElementR(`span[role="button"]`, strconv.Itoa(date.Day())).MustClick()

	// active time selector
	page.MustElement("div.time-picker_e4ca4").MustClick()
	// select correct hour and time
	hourList := page.MustElement("ul.hour-wrap").MustElements("li")

	for _, hour := range hourList {
		hourText, _ := hour.Text()
		if strings.Contains(hourText, strconv.Itoa(date.Hour())) {
			hour.MustClick()
			break
		}
	}

	minsList := page.MustElement("ul.mins-wrap").MustElements("li")

	for _, min := range minsList {
		minText, _ := min.Text()
		if strings.Contains(minText, strconv.Itoa(utils.RoundToNextFive(date.Minute()))) {
			min.MustClick()
			break
		}
	}

}

func getTrainpalJourneyDetails(journey *rod.Element) ScrapeResult {
	// get departure time
	departureTime, _ := journey.MustElement("div.from_fa71c").Text()
	// get arrival time
	arrivalTime, _ := journey.MustElement("div.to_cc86d").Text()
	// get price
	price, _ := journey.MustElement("div.price_f360e").MustElement("div").Text()

	return ScrapeResult{
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		Price:         utils.PriceToFloat(utils.SanitisePrice(price)),
	}
}
