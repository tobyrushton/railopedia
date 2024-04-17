package scrapes

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/tobyrushton/railopedia/packages/functions/src/utils"
)

var raileasyUrl string = "https://new.raileasy.co.uk/"

func ScrapeRaileasy(req Request) (ScrapeResults, error) {
	// open browser
	page := rod.New().MustConnect().MustPage(raileasyUrl)

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

	// submit form
	page.MustElement("#cookie-banner-accept").MustClick()
	page.MustElement("#search-button").MustClick()

	page.MustWaitStable().MustScreenshot("screenshot.png")

	return nil, nil
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
	fmt.Println(monthString)
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
	fmt.Println(hourId, minuteId)
	// select the correct hour
	hour := date.Hour()
	page.MustElement(hourId).MustSelect(fmt.Sprintf("%d", hour))

	// select the correct minute
	minute := utils.RoundToNextFifteen(date.Minute())
	page.MustElement(minuteId).MustSelect(fmt.Sprintf("%d", minute))
}
