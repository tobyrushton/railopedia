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

func Scrape(req Request) ([]JourneyWithPrices, error) {
	trainlineChannel := make(chan ScrapeResultNonConditional)
	trainpalChannel := make(chan ScrapeResultNonConditional)
	trainticketsChannel := make(chan ScrapeResultsConditional)
	// raileasyChannel := make(chan ScrapeResultsConditional)

	errChannel := make(chan error)

	// scrape each site concurrently
	go func() {
		val, err := ScrapeTrainline(req)
		if err != nil {
			errChannel <- err
		} else {
			trainlineChannel <- val
		}
	}()

	go func() {
		val, err := ScrapeTrainpal(req)
		if err != nil {
			errChannel <- err
		} else {
			trainpalChannel <- val
		}
	}()

	go func() {
		val, err := ScrapeTraintickets(req)
		trainticketsChannel <- val
		if err != nil {
			errChannel <- err
		} else {
			trainticketsChannel <- val
		}
	}()

	// go func() {
	// 	val, _ := ScrapeRaileasy(req)
	// 	raileasyChannel <- val
	// }()

	// key for map should be [depart iso],[arrive iso]
	journeys := make(map[string]JourneyWithPrices)

	// wait for all channels to return
	for i := 0; i < 3; i++ {
		select {
		case trainline := <-trainlineChannel:
			aggregateNonConditionalScrapeResults(trainline, &journeys, "trainline")
		case trainpal := <-trainpalChannel:
			aggregateNonConditionalScrapeResults(trainpal, &journeys, "trainpal")
		case traintickets := <-trainticketsChannel:
			aggregateConditionalScrapeResults(traintickets, &journeys, "traintickets")
			// case raileasy := <-raileasyChannel:
			// 	addRaileasyResults(raileasy, &journeys, "raileasy") <- broken
		case err := <-errChannel:
			return nil, err
		}
	}

	// convert map to slice
	journeySlice := make([]JourneyWithPrices, 0, len(journeys))

	for _, journey := range journeys {
		journeySlice = append(journeySlice, journey)
	}

	return journeySlice, nil
}
