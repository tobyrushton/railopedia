package scrapes

import (
	"fmt"
	"sort"
)

type Price struct {
	Provider string
	Price    float32
	Link     string
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
	raileasyChannel := make(chan ScrapeResultsConditional)

	errChannel := make(chan error)

	// catches any panics and sends them to the error channel
	catchPanic := func(t string) {
		fmt.Println("finished: ", t)
		if r := recover(); r != nil {
			fmt.Println("panic: ", r)
			errChannel <- fmt.Errorf("panic: %v", r)
		}
	}

	// scrape each site concurrently
	go func() {
		defer catchPanic("tl")
		val, err := ScrapeTrainline(req)
		if err != nil {
			errChannel <- err
		} else {
			trainlineChannel <- val
		}
	}()

	go func() {
		defer catchPanic("tp")
		val, err := ScrapeTrainpal(req)
		if err != nil {
			errChannel <- err
		} else {
			trainpalChannel <- val
		}
	}()

	go func() {
		defer catchPanic("tt")
		val, err := ScrapeTraintickets(req)
		trainticketsChannel <- val
		if err != nil {
			errChannel <- err
		} else {
			trainticketsChannel <- val
		}
	}()

	go func() {
		defer catchPanic("re")
		fmt.Println("test")
		val, err := ScrapeRaileasy(req)
		if err != nil {
			errChannel <- err
		} else {
			raileasyChannel <- val
		}
	}()

	// key for map should be [depart iso],[arrive iso]
	journeys := make(map[string]JourneyWithPrices)

	// wait for all channels to return
	for i := 0; i < 4; i++ {
		select {
		case trainline := <-trainlineChannel:
			aggregateNonConditionalScrapeResults(trainline, &journeys, "trainline")
		case trainpal := <-trainpalChannel:
			aggregateNonConditionalScrapeResults(trainpal, &journeys, "trainpal")
		case traintickets := <-trainticketsChannel:
			aggregateConditionalScrapeResults(traintickets, &journeys, "traintickets")
		case raileasy := <-raileasyChannel:
			fmt.Println("raileasy")
			aggregateConditionalScrapeResults(raileasy, &journeys, "raileasy")
		case <-errChannel:
			continue
		}
	}

	// convert map to slice
	journeySlice := make([]JourneyWithPrices, 0, len(journeys))

	for _, journey := range journeys {
		journeySlice = append(journeySlice, journey)
	}

	// sort by departure time
	sort.Slice(journeySlice, func(i, j int) bool {
		return journeySlice[i].DepartureTime < journeySlice[j].DepartureTime
	})

	return journeySlice, nil
}
