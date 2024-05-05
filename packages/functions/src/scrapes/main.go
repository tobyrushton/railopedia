package scrapes

import (
	"github.com/tobyrushton/railopedia/packages/functions/src/utils"
)

type ScrapeResult struct {
	DepartureTime string // ISO
	ArrivalTime   string // ISO
	Price         float32
}

type ScrapeResultNonConditional struct {
	Outbound []ScrapeResult
	Return   []ScrapeResult
	Link     string
}

type ScrapeResultConditional struct {
	DepartureTime string // ISO
	ArrivalTime   string // ISO
	Link          string
	Price         map[string]float32 // format iso8601:price
}

type ScrapeResults []ScrapeResult
type ScrapeResultsConditional []ScrapeResultConditional

type Request struct {
	Origin      string
	Destination string
	Departure   string // ISO
	Return      string // ISO
}

var iso8601Layout string = "2006-01-02T15:04:05Z0700"

// func getStationByCode(code string) (string, error) {
// 	jsonFile, err := os.Open("../../../../data/station-list.json")

// 	if err != nil {
// 		return "", err
// 	}

// 	defer jsonFile.Close()

// 	byteValue, _ := io.ReadAll(jsonFile)

// 	var stations []Station

// 	json.Unmarshal(byteValue, &stations)

// 	for _, station := range stations {
// 		if station.Code == code {
// 			return strings.ReplaceAll(station.Name, "-", " "), nil
// 		}
// 	}

// 	return "", errors.New("Station not found")
// }

func aggregateNonConditionalScrapeResultsSingle(results ScrapeResultNonConditional, journeys *map[string]Journey, provider string) {
	for _, result := range results.Outbound {
		key := result.DepartureTime + "," + result.ArrivalTime
		if journey, ok := (*journeys)[key]; ok {
			journey.Prices = append(journey.Prices, Price{Provider: provider, Price: result.Price, Link: results.Link})
			(*journeys)[key] = journey
		} else {
			(*journeys)[key] = Journey{
				DepartureTime: result.DepartureTime,
				ArrivalTime:   result.ArrivalTime,
				Prices:        []Price{{Provider: provider, Price: result.Price, Link: results.Link}},
			}
		}
	}
}

func aggregateNonConditionalScrapeResultsReturn(results ScrapeResultNonConditional, journeys *map[string]JourneyWithPrices, provider string) {
	for _, result := range results.Outbound {
		key := result.DepartureTime + "," + result.ArrivalTime
		if journey, ok := (*journeys)[key]; ok {
			price := journey.Prices

			for _, returnJourney := range results.Return {
				index, found := 0, false
				priceItem := Journey{}
				for index < len(price) && !found {
					if price[index].DepartureTime == returnJourney.DepartureTime && price[index].ArrivalTime == returnJourney.ArrivalTime {
						priceItem = price[index]
						found = true
					} else {
						index++
					}
				}
				if !found {
					priceItem = Journey{
						DepartureTime: returnJourney.DepartureTime,
						ArrivalTime:   returnJourney.ArrivalTime,
						Prices:        []Price{{Provider: provider, Price: result.Price + returnJourney.Price, Link: results.Link}},
					}
					price = append(price, priceItem)
				} else {
					priceItem.Prices = append(priceItem.Prices, Price{Provider: provider, Price: result.Price + returnJourney.Price, Link: results.Link})
					price[index] = priceItem
				}
			}

			journey.Prices = price
			(*journeys)[key] = journey
		} else {
			(*journeys)[key] = JourneyWithPrices{
				DepartureTime: result.DepartureTime,
				ArrivalTime:   result.ArrivalTime,
				Prices:        make([]Journey, 0),
			}

			journey = (*journeys)[key]

			price := make([]Journey, 0)

			for _, returnJourney := range results.Return {
				price = append(price, Journey{
					DepartureTime: returnJourney.DepartureTime,
					ArrivalTime:   returnJourney.ArrivalTime,
					Prices:        []Price{{Provider: provider, Price: result.Price + returnJourney.Price, Link: results.Link}},
				})
			}

			journey.Prices = price
			(*journeys)[key] = journey
		}
	}
}

func aggregateConditionalScrapeResultsSingle(results ScrapeResultsConditional, journeys *map[string]Journey, provider string) {
	for _, result := range results {
		key := result.DepartureTime + "," + result.ArrivalTime
		if journey, ok := (*journeys)[key]; ok {
			journey.Prices = append(journey.Prices, Price{Provider: provider, Price: result.Price[key], Link: result.Link})
			(*journeys)[key] = journey
		} else {
			(*journeys)[key] = Journey{
				DepartureTime: result.DepartureTime,
				ArrivalTime:   result.ArrivalTime,
				Prices:        []Price{{Provider: provider, Price: result.Price[key], Link: result.Link}},
			}
		}
	}
}

func aggregateConditionalScrapeResultsReturn(results ScrapeResultsConditional, journeys *map[string]JourneyWithPrices, provider string) {
	for _, result := range results {
		key := result.DepartureTime + "," + result.ArrivalTime
		if journey, ok := (*journeys)[key]; ok {
			price := journey.Prices

			for key, val := range result.Price {
				departureTime, arrivalTime := utils.SplitString(key, ",")

				if journey, i, found := findJourney(price, departureTime, arrivalTime); found {
					journey.Prices = append(journey.Prices, Price{Provider: provider, Price: val, Link: result.Link})
					price[i] = journey
				} else {
					price = append(price, Journey{
						DepartureTime: departureTime,
						ArrivalTime:   arrivalTime,
						Prices:        []Price{{Provider: provider, Price: val, Link: result.Link}},
					})
				}
			}

		} else {
			(*journeys)[key] = JourneyWithPrices{
				DepartureTime: result.DepartureTime,
				ArrivalTime:   result.ArrivalTime,
				Prices:        make([]Journey, 0),
			}

			journey = (*journeys)[key]

			price := make([]Journey, 0)

			for key, val := range result.Price {
				departureTime, arrivalTime := utils.SplitString(key, ",")
				price = append(price, Journey{
					DepartureTime: departureTime,
					ArrivalTime:   arrivalTime,
					Prices:        []Price{{Provider: provider, Price: val, Link: result.Link}},
				})
			}

			journey.Prices = price
			(*journeys)[key] = journey
		}
	}
}

func findJourney(journeys []Journey, departureTime string, arrivalTime string) (Journey, int, bool) {
	for i, journey := range journeys {
		if journey.DepartureTime == departureTime && journey.ArrivalTime == arrivalTime {
			return journey, i, true
		}
	}

	return Journey{}, -1, false
}
