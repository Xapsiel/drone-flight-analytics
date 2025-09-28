package service

import (
	"context"
	"time"

	"github.com/Xapsiel/bpla_dashboard/internal/config"
	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

type Repository interface {
	SaveMessage(ctx context.Context, mes *model.ParsedMessage, fileID int) error

	TotalFlightAndAVGDuration(ctx context.Context, regID int, year int) ([]struct {
		RegionCode         int
		RegionName         string
		TotalFlight        int
		AvgDurationMinutes float32
	}, error)
	GetDailyFlightMetrics(ctx context.Context, regID int, year int) ([]struct {
		RegionCode         int
		RegionName         string
		AvgDailyFlights    float64
		MedianDailyFlights float64
	}, error)
	GetPeakLoad(ctx context.Context, regID int, year int) ([]struct {
		RegionCode int
		RegionName string
		PeakLoad   int
	}, error)
	GetMonthlyGrowth(ctx context.Context, regID int, year int) ([]struct {
		RegionCode    int
		RegionName    string
		MonthlyGrowth map[int]float64
	}, error)
	GetFlightDensity(ctx context.Context, regID int, year int) ([]struct {
		RegionCode    int
		RegionName    string
		FlightDensity float64
	}, error)
	GetFlightTimes(ctx context.Context, regID int, year int) ([]struct {
		RegionCode     int
		RegionName     string
		MorningFlights int
		DayFlights     int
		EveningFlights int
		NightFlights   int
	}, error)
	GetZeroFlightDays(ctx context.Context, regID int, year int) ([]struct {
		RegionCode     int
		RegionName     string
		ZeroFlightDays []time.Time
	}, error)
	GetTotalDistance(ctx context.Context, regID int, year int) ([]struct {
		RegionCode      int
		RegionName      string
		TotalDistanceKm float64
	}, error)

	TotalFlightAndAVGDurationAllRussia(ctx context.Context, year int) (int, float32, error)
	GetPeakLoadAllRussia(ctx context.Context, year int) (int, error)
	GetDailyFlightMetricsAllRussia(ctx context.Context, year int) (float64, float64, error)
	GetRussiaMonthlyGrowth(ctx context.Context, year int) (map[int]float64, error)
	GetFlightDensityAllRussia(ctx context.Context, year int) (float64, error)
	GetFlightTimesAllRussia(ctx context.Context, year int) (int, int, int, int, error)
	GetZeroFlightDaysAllRussia(ctx context.Context, year int) ([]time.Time, error)
	GetTotalDistanceAllRussia(ctx context.Context, year int) (float64, error)

	GetRegions(ctx context.Context) []model.District
	GetFlightYears(ctx context.Context) []int
	UpdateMetrics(ctx context.Context, metrics chan *model.Metrics) error
}

type Service struct {
	*UserService
	*MetricsService
	*ParserService
}

func New(repo Repository, cfg config.OidcConfig) Service {
	return Service{
		UserService:    NewUserService(repo, cfg),
		ParserService:  NewParserService(repo),
		MetricsService: NewMetricsService(repo),
	}
}
