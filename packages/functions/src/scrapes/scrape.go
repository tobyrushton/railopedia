package scrapes

type Price struct {
	Provider string
	Price    float32
}

type Journey struct {
	DepartureTime string
	ArrivalTime   string
	Prices        []Price
}

type JourneyWithPrices struct {
	DepartureTime string
	ArrivalTime   string
	Prices        []Journey
}

func Scrape(req Request) []JourneyWithPrices {
	trainlineChannel := make(chan ScrapeResultNonConditional)
	trainpalChannel := make(chan ScrapeResultNonConditional)
	// trainticketsChannel := make(chan ScrapeResultsConditional)
	// raileasyChannel := make(chan ScrapeResultsConditional)

	// scrape each site concurrently
	go func() {
		val, _ := ScrapeTrainline(req)
		trainlineChannel <- val
	}()

	go func() {
		val, _ := ScrapeTrainpal(req)
		trainpalChannel <- val
	}()

	// go func() {
	// 	val, _ := ScrapeTraintickets(req)
	// 	trainticketsChannel <- val
	// }()

	// go func() {
	// 	val, _ := ScrapeRaileasy(req)
	// 	raileasyChannel <- val
	// }()

	// key for map should be [depart iso],[arrive iso]
	journeys := make(map[string]JourneyWithPrices)

	// wait for all channels to return
	for i := 0; i < 2; i++ {
		select {
		case trainline := <-trainlineChannel:
			aggregateNonConditionalScrapeResults(trainline, &journeys, "trainline")
		case trainpal := <-trainpalChannel:
			aggregateNonConditionalScrapeResults(trainpal, &journeys, "trainpal")
			// case traintickets := <-trainticketsChannel:
			// 	// do something with traintickets
			// case raileasy := <-raileasyChannel:
			// 	// do something with raileasy
		}
	}

	// convert map to slice
	journeySlice := make([]JourneyWithPrices, 0, len(journeys))

	for _, journey := range journeys {
		journeySlice = append(journeySlice, journey)
	}

	return journeySlice
}
