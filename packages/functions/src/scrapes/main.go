package scrapes

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

type ScrapeResult struct {
	DepartureTime string // ISO
	ArrivalTime   string // ISO
	// Link          string
	Price float32
}

type ScrapeResultConditional struct {
	DepartureTime string // ISO
	ArrivalTime   string // ISO
	// Link          string
	Price map[string]float32 // format iso8601:price
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
