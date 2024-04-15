package scrapes

type ScrapeResult struct {
	DepartureTime string // ISO
	ArrivalTime   string // ISO
	// Link          string
	Price float32
}

type ScrapeResults []ScrapeResult

type Request struct {
	Origin      string
	Destination string
	Departure   string // ISO
	Return      string // ISO
}
