package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

type MetricsService struct {
	repo    Repository
	metrics chan *model.Metrics
}

func NewMetricsService(repo Repository) *MetricsService {
	return &MetricsService{repo: repo, metrics: make(chan *model.Metrics, 5)}
}

func (s *MetricsService) Update(ctx context.Context) error {
	reg := s.repo.GetRegions(ctx)
	years := s.repo.GetFlightYears(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for _, region := range reg {
			for _, year := range years {
				s.getMetrics(ctx, region, year)
			}
		}
	}(wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for _, year := range years {
			s.getMetricsAllRussia(ctx, year)
		}
	}(wg)
	go func(wg *sync.WaitGroup) {
		err := s.repo.UpdateMetrics(context.Background(), s.metrics)
		if err != nil {
			slog.Error(fmt.Sprintf("error update metrics: %v", err))
		}
	}(wg)
	wg.Wait()
	return nil

}

func (s *MetricsService) getMetrics(ctx context.Context, region model.District, year int) {
	metrics := &model.Metrics{
		RegionId:   *region.Gid,
		RegionName: *region.Name,
		Year:       year,
	}
	peakLoad, err := s.repo.GetPeakLoad(ctx, *region.Gid, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get peak load: %v", err))
	} else {
		if peakLoad != nil {
			metrics.PeakLoad = peakLoad[0].PeakLoad
		}
	}
	total_flight_and_avg_dur, err := s.repo.TotalFlightAndAVGDuration(ctx, *region.Gid, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get total flight and avg duration: %v", err))
	} else {
		if total_flight_and_avg_dur != nil {
			metrics.TotalFlight = total_flight_and_avg_dur[0].TotalFlight
			metrics.AvgDurationMinutes = total_flight_and_avg_dur[0].AvgDurationMinutes
		}
	}
	monthlyGrowth, err := s.repo.GetMonthlyGrowth(ctx, *region.Gid, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get monthly growth: %v", err))
	} else {
		if monthlyGrowth != nil {
			metrics.MonthlyGrowth = monthlyGrowth[0].MonthlyGrowth
		}
	}
	flightTimes, err := s.repo.GetFlightTimes(ctx, *region.Gid, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get flight times: %v", err))
	} else {
		if flightTimes != nil {
			metrics.MorningFlights = flightTimes[0].MorningFlights

			metrics.DayFlights = flightTimes[0].DayFlights
			metrics.EveningFlights = flightTimes[0].EveningFlights
			metrics.NightFlights = flightTimes[0].NightFlights
		}
	}
	flightDensity, err := s.repo.GetFlightDensity(ctx, *region.Gid, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get flight density: %v", err))
	} else {
		if flightDensity != nil {

			metrics.FlightDensity = flightDensity[0].FlightDensity
		}
	}
	dailyFlight, err := s.repo.GetDailyFlightMetrics(ctx, *region.Gid, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get daily flight metrics: %v", err))
	} else {
		if dailyFlight != nil {

			metrics.AvgDailyFlights = dailyFlight[0].AvgDailyFlights
			metrics.MedianDailyFlights = dailyFlight[0].MedianDailyFlights
		}
	}

	zeroFlight, err := s.repo.GetZeroFlightDays(ctx, *region.Gid, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get zero flight days: %v", err))
	} else {
		if zeroFlight != nil {

			metrics.ZeroFlightDays = zeroFlight[0].ZeroFlightDays
		}
	}
	total_distance, err := s.repo.GetTotalDistance(ctx, *region.Gid, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get total distance: %v", err))
	} else {
		if total_distance != nil {

			metrics.TotalDistance = total_distance[0].TotalDistanceKm
		}

	}
	s.metrics <- metrics
}

func (s *MetricsService) getMetricsAllRussia(ctx context.Context, year int) {
	metrics := &model.Metrics{
		RegionName: "Российская Федерация",
		Year:       year,
	}
	peakLoad, err := s.repo.GetPeakLoadAllRussia(ctx, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get peak load: %v", err))
	} else {
		metrics.PeakLoad = peakLoad
	}
	total_flight, avg_dur, err := s.repo.TotalFlightAndAVGDurationAllRussia(ctx, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get total flight and avg duration: %v", err))
	} else {
		metrics.TotalFlight = total_flight
		metrics.AvgDurationMinutes = avg_dur
	}
	monthlyGrowth, err := s.repo.GetRussiaMonthlyGrowth(ctx, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get monthly growth: %v", err))
	} else {
		metrics.MonthlyGrowth = monthlyGrowth
	}
	m, d, e, n, err := s.repo.GetFlightTimesAllRussia(ctx, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get flight times: %v", err))
	} else {
		metrics.MorningFlights = m

		metrics.DayFlights = d
		metrics.EveningFlights = e
		metrics.NightFlights = n
	}
	flightDensity, err := s.repo.GetFlightDensityAllRussia(ctx, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get flight density: %v", err))
	} else {

		metrics.FlightDensity = flightDensity
	}
	avgDailyFlights, medianDailyFlights, err := s.repo.GetDailyFlightMetricsAllRussia(ctx, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get daily flight metrics: %v", err))
	} else {

		metrics.AvgDailyFlights = avgDailyFlights
		metrics.MedianDailyFlights = medianDailyFlights
	}

	zeroFlight, err := s.repo.GetZeroFlightDaysAllRussia(ctx, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get zero flight days: %v", err))
	} else {

		metrics.ZeroFlightDays = zeroFlight
	}
	total_distance, err := s.repo.GetTotalDistanceAllRussia(ctx, year)
	if err != nil {
		slog.Error(fmt.Sprintf("error get total distance: %v", err))
	} else {
		metrics.TotalDistance = total_distance
	}

	s.metrics <- metrics
}
