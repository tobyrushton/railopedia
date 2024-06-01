package scrapes

import (
	"fmt"
	"sort"
	"time"
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

func Scrape(req Request) ([]JourneyWithPrices, []Journey, error) {
	trainlineChannel := make(chan ScrapeResultNonConditional)
	trainpalChannel := make(chan ScrapeResultNonConditional)
	trainticketsChannel := make(chan ScrapeResultsConditional)
	raileasyChannel := make(chan ScrapeResultsConditional)

	errChannel := make(chan error)
	timeoutChannel := make(chan struct{})

	// catches any panics and sends them to the error channel
	catchPanic := func(t string) {
		fmt.Println("finished: ", t)
		if r := recover(); r != nil {
			fmt.Println("panic: ", r)
			errChannel <- fmt.Errorf("panic: %v", r)
		}
	}

	go func() {
		time.Sleep(20 * time.Second)
		timeoutChannel <- struct{}{}
		timeoutChannel <- struct{}{}
		timeoutChannel <- struct{}{}
		timeoutChannel <- struct{}{}
	}()

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
		if err != nil {
			errChannel <- err
		} else {
			trainticketsChannel <- val
		}
	}()

	go func() {
		defer catchPanic("re")
		val, err := ScrapeRaileasy(req)
		if err != nil {
			errChannel <- err
		} else {
			raileasyChannel <- val
		}
	}()

	// key for map should be [depart iso],[arrive iso]
	journeysReturn := make(map[string]JourneyWithPrices)
	journeysSingle := make(map[string]Journey)

	isReturn := req.Return != ""

	// wait for all channels to return
	for i := 0; i < 4; i++ {
		select {
		case trainline := <-trainlineChannel:
			if isReturn {
				aggregateNonConditionalScrapeResultsReturn(trainline, &journeysReturn, "Trainline")
			} else {
				aggregateNonConditionalScrapeResultsSingle(trainline, &journeysSingle, "Trainline")
			}
		case trainpal := <-trainpalChannel:
			if isReturn {
				aggregateNonConditionalScrapeResultsReturn(trainpal, &journeysReturn, "Trainpal")
			} else {
				aggregateNonConditionalScrapeResultsSingle(trainpal, &journeysSingle, "Trainpal")
			}
		case traintickets := <-trainticketsChannel:
			if isReturn {
				fmt.Println("traintickets")
				aggregateConditionalScrapeResultsReturn(traintickets, &journeysReturn, "Traintickets")
			} else {
				aggregateConditionalScrapeResultsSingle(traintickets, &journeysSingle, "Traintickets")
			}
		case raileasy := <-raileasyChannel:
			if isReturn {
				aggregateConditionalScrapeResultsReturn(raileasy, &journeysReturn, "Raileasy")
			} else {
				aggregateConditionalScrapeResultsSingle(raileasy, &journeysSingle, "Raileasy")
			}
		case <-errChannel:
			continue
		case <-timeoutChannel:
			continue
		}
	}

	// convert map to slice
	journeyReturnSlice := make([]JourneyWithPrices, 0, len(journeysReturn))
	journeySingleSlice := make([]Journey, 0, len(journeysSingle))

	for _, journey := range journeysReturn {
		journeyReturnSlice = append(journeyReturnSlice, journey)
	}

	for _, journey := range journeysSingle {
		journeySingleSlice = append(journeySingleSlice, journey)
	}

	// sort by departure time
	sort.Slice(journeyReturnSlice, func(i, j int) bool {
		return journeyReturnSlice[i].DepartureTime < journeyReturnSlice[j].DepartureTime
	})
	sort.Slice(journeySingleSlice, func(i, j int) bool {
		return journeySingleSlice[i].DepartureTime < journeySingleSlice[j].DepartureTime
	})

	return journeyReturnSlice, journeySingleSlice, nil
}
