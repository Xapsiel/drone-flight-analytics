package model

type Metrics struct {
	RegionId           int
	RegionName         string
	PeakLoad           int
	TotalFlight        int
	AvgDurationMinutes float32
	MonthlyGrowth      map[int]float64
	MorningFlights     int
	DayFlights         int
	EveningFlights     int
	NightFlights       int
	FlightDensity      float64
	AvgDailyFlights    float64
	MedianDailyFlights float64
	ZeroFlightDays     int
	Year               int
	TotalDistance      float64
}
