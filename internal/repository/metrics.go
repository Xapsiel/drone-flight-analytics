package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

func (r *Repository) GetMetrics(ctx context.Context, id int, year int) (model.Metrics, error) {
	query := `
				SELECT 
					region_code,region_name,
					total_flight,avg_duration_minutes,
					total_distance_km,peak_load,
					avg_daily_flights,median_daily_flights,
					monthly_growth,flight_density,
					morning_flights,day_flights,
					evening_flights,night_flights,
					zero_flight_days,date
				FROM flight_metrics
				WHERE region_code = $1 AND date = $2
			 `
	res := model.Metrics{}
	jsonDate := []byte{}
	row := r.db.QueryRow(ctx, query, id, year)
	err := row.Scan(
		&res.RegionId, &res.RegionName,
		&res.TotalFlight, &res.AvgDurationMinutes,
		&res.TotalDistance, &res.PeakLoad,
		&res.AvgDailyFlights, &res.MedianDailyFlights,
		&jsonDate, &res.FlightDensity,
		&res.MorningFlights, &res.DayFlights,
		&res.EveningFlights, &res.NightFlights,
		&res.ZeroFlightDays, &res.Year,
	)
	if err != nil {
		slog.Error("error with scan metrics: %v", err)
		return res, err
	}
	err = json.Unmarshal(jsonDate, &res.MonthlyGrowth)
	if err != nil {
		slog.Error("error with unmarshal metrics: %v", err)
		return res, err
	}
	return res, nil

}

func (r *Repository) TotalFlightAndAVGDuration(ctx context.Context, regID int, year int) ([]struct {
	RegionCode         int
	RegionName         string
	TotalFlight        int
	AvgDurationMinutes float32
}, error) {

	query := `
		SELECT
			m.region AS region_code,
			ds.name AS region_name,
			COUNT(DISTINCT m.sid) AS total_flight,
			AVG(EXTRACT(EPOCH FROM (
				CASE
					WHEN m.atd > m.ata THEN (m.dof + INTERVAL '1 day' + m.ata) - (m.dof + m.atd)
					ELSE (m.dof + m.ata) - (m.dof + m.atd)
				END
			)) / 60) AS avg_duration_minutes
		FROM messages m
		LEFT JOIN district_shapes ds ON m.region = ds.gid
		WHERE m.ata IS NOT NULL AND m.arr_coordinate IS NOT NULL
	`

	args := []interface{}{}
	if regID != 0 {
		query += " AND m.region = $1"
		args = append(args, regID)
	}
	if year != 0 {
		paramPos := len(args) + 1
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM m.dof) = $%d", paramPos)
		args = append(args, year)
	}

	query += " GROUP BY m.region, ds.name"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query flight metrics: %w", err)
	}
	defer rows.Close()

	var results []struct {
		RegionCode         int
		RegionName         string
		TotalFlight        int
		AvgDurationMinutes float32
	}

	for rows.Next() {
		var r struct {
			RegionCode         int
			RegionName         string
			TotalFlight        int
			AvgDurationMinutes float32
		}
		if err := rows.Scan(&r.RegionCode, &r.RegionName, &r.TotalFlight, &r.AvgDurationMinutes); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (r *Repository) GetPeakLoad(ctx context.Context, regID int, year int) ([]struct {
	RegionCode int
	RegionName string
	PeakLoad   int
}, error) {
	query := `
		WITH hourly_load AS (
			SELECT
				ds.gid AS region_code,
				ds.name AS region_name,
				DATE_TRUNC('hour', m.dof + m.atd) AS hour,
				COUNT(*) AS hourly_count
			FROM district_shapes ds
			LEFT JOIN messages m ON m.region = ds.gid AND m.ata IS NOT NULL
			WHERE 1=1
	`
	args := []interface{}{}
	if regID != 0 {
		query += " AND m.region = $1"
		args = append(args, regID)
	}
	if year != 0 {
		paramPos := len(args) + 1
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM m.dof) = $%d", paramPos)
		args = append(args, year)
	}
	query += `
			GROUP BY ds.gid, ds.name, DATE_TRUNC('hour', m.dof + m.atd)
		),
		peak_hours AS (
			SELECT
				region_code,
				region_name,
				hourly_count,
				ROW_NUMBER() OVER (PARTITION BY region_code ORDER BY hourly_count DESC, hour ASC) AS rn
			FROM hourly_load
		)
		SELECT
			region_code,
			region_name,
			hourly_count AS peak_load
		FROM peak_hours
		WHERE rn = 1
		ORDER BY region_code;
	`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query peak load: %w", err)
	}
	defer rows.Close()

	var results []struct {
		RegionCode int
		RegionName string
		PeakLoad   int
	}
	for rows.Next() {
		var r struct {
			RegionCode int
			RegionName string
			PeakLoad   int
		}
		if err := rows.Scan(&r.RegionCode, &r.RegionName, &r.PeakLoad); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, r)
	}
	return results, nil
}

func (r *Repository) GetMonthlyGrowth(ctx context.Context, regID int, year int) ([]struct {
	RegionCode    int
	RegionName    string
	MonthlyGrowth map[int]float64
}, error) {
	query := `
		WITH monthly_counts AS (
			SELECT
				m.region AS region_code,
				ds.name AS region_name,
				DATE_TRUNC('month', m.dof) AS month_date,
				COUNT(m.sid) AS monthly_flight_count
			FROM messages m
			JOIN district_shapes ds ON m.region = ds.gid
			WHERE m.ata IS NOT NULL
	`
	args := []interface{}{}
	if regID != 0 {
		query += " AND m.region = $1"
		args = append(args, regID)
	}
	if year != 0 {
		query += " AND EXTRACT(YEAR FROM m.dof) = $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, year)
	}
	query += `
			GROUP BY m.region, ds.name, DATE_TRUNC('month', m.dof)
		),
		growth_calc AS (
			SELECT
				region_code,
				region_name,
				month_date,
				monthly_flight_count,
				LAG(monthly_flight_count) OVER (PARTITION BY region_code ORDER BY month_date) AS prev_month_count
			FROM monthly_counts
		),
		growth_data AS (
			SELECT
				region_code,
				region_name,
				month_date,
				EXTRACT(MONTH FROM month_date) AS month_number,
				CASE
					WHEN prev_month_count > 0
					THEN (monthly_flight_count::NUMERIC - prev_month_count) / prev_month_count * 100
					ELSE 0
				END AS growth_percentage,
				ROW_NUMBER() OVER (PARTITION BY region_code ORDER BY month_date DESC) AS rn
			FROM growth_calc
		),
		growth_json AS (
			SELECT
				region_code,
				region_name,
				json_object_agg(month_number, growth_percentage) AS monthly_growth_json
			FROM growth_data
			GROUP BY region_code, region_name
		)
		SELECT
			region_code,
			region_name,
			monthly_growth_json AS monthly_growth
		FROM growth_json
		WHERE EXISTS (
			SELECT 1
			FROM growth_data gd
			WHERE gd.region_code = growth_json.region_code
			AND gd.rn = 1
		)
	`

	var rows pgx.Rows
	var err error
	if regID != 0 {
		rows, err = r.db.Query(ctx, query, regID, year)
	} else {
		rows, err = r.db.Query(ctx, query, year)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly growth: %w", err)
	}
	defer rows.Close()

	var results []struct {
		RegionCode    int
		RegionName    string
		MonthlyGrowth map[int]float64
	}
	for rows.Next() {
		var regionCode int
		var regionName string
		var monthlyGrowthJSON string // Принимаем JSON как строку
		if err := rows.Scan(&regionCode, &regionName, &monthlyGrowthJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var growthMap map[int]float64
		if err := json.Unmarshal([]byte(monthlyGrowthJSON), &growthMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		monthlyGrowth := make(map[int]float64, 12)
		for month, growth := range growthMap {
			if month >= 1 && month <= 12 {
				monthlyGrowth[month-1] = growth
			}
		}

		results = append(results, struct {
			RegionCode    int
			RegionName    string
			MonthlyGrowth map[int]float64
		}{
			RegionCode:    regionCode,
			RegionName:    regionName,
			MonthlyGrowth: monthlyGrowth,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}
func (r *Repository) GetDailyFlightMetrics(ctx context.Context, regID int, year int) ([]struct {
	RegionCode         int
	RegionName         string
	AvgDailyFlights    float64
	MedianDailyFlights float64
}, error) {

	query := `
		WITH daily_flights AS (
			SELECT
				m.region AS region_code,
				ds.name AS region_name,
				CASE
					WHEN m.atd > m.ata THEN m.dof + INTERVAL '1 day'
					ELSE m.dof
				END AS flight_date,
				COUNT(m.sid) AS daily_flight_count
			FROM messages m
			JOIN district_shapes ds ON m.region = ds.gid
			WHERE m.ata IS NOT NULL
	`

	args := []interface{}{}
	if regID != 0 {
		query += " AND m.region = $1"
		args = append(args, regID)
	}
	if year != 0 {
		paramPos := len(args) + 1
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM m.dof) = $%d", paramPos)
		args = append(args, year)
	}

	query += `
			GROUP BY m.region, ds.name, CASE
				WHEN m.atd > m.ata THEN m.dof + INTERVAL '1 day'
				ELSE m.dof
			END
		)
		SELECT
			region_code,
			region_name,
			AVG(daily_flight_count) AS avg_daily_flights,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY daily_flight_count) AS median_daily_flights
		FROM daily_flights
		GROUP BY region_code, region_name
	`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily flight metrics: %w", err)
	}
	defer rows.Close()

	var results []struct {
		RegionCode         int
		RegionName         string
		AvgDailyFlights    float64
		MedianDailyFlights float64
	}

	for rows.Next() {
		var r struct {
			RegionCode         int
			RegionName         string
			AvgDailyFlights    float64
			MedianDailyFlights float64
		}
		if err := rows.Scan(&r.RegionCode, &r.RegionName, &r.AvgDailyFlights, &r.MedianDailyFlights); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}
func (r *Repository) GetFlightDensity(ctx context.Context, regID int, year int) ([]struct {
	RegionCode    int
	RegionName    string
	FlightDensity float64
}, error) {
	query := `
		SELECT
			m.region AS region_code,
			ds.name AS region_name,
			(COUNT(DISTINCT m.sid)::NUMERIC / (ds.area_km2 / 1000)) AS flight_density
		FROM messages m
		JOIN district_shapes ds ON m.region = ds.gid
		WHERE m.ata IS NOT NULL AND m.arr_coordinate IS NOT NULL
	`

	args := []interface{}{}
	if regID != 0 {
		query += " AND m.region = $1"
		args = append(args, regID)
	}
	if year != 0 {
		paramPos := len(args) + 1
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM m.dof) = $%d", paramPos)
		args = append(args, year)
	}

	query += `
		GROUP BY m.region, ds.name, ds.area_km2
	`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query flight density: %w", err)
	}
	defer rows.Close()

	var results []struct {
		RegionCode    int
		RegionName    string
		FlightDensity float64
	}

	for rows.Next() {
		var r struct {
			RegionCode    int
			RegionName    string
			FlightDensity float64
		}
		if err := rows.Scan(&r.RegionCode, &r.RegionName, &r.FlightDensity); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (r *Repository) GetFlightTimes(ctx context.Context, regID int, year int) ([]struct {
	RegionCode     int
	RegionName     string
	MorningFlights int
	DayFlights     int
	EveningFlights int
	NightFlights   int
}, error) {
	query := `
		SELECT
			m.region AS region_code,
			ds.name AS region_name,
			SUM(CASE WHEN EXTRACT(HOUR FROM m.atd) BETWEEN 6 AND 11 THEN 1 ELSE 0 END) AS morning_flights,
			SUM(CASE WHEN EXTRACT(HOUR FROM m.atd) BETWEEN 12 AND 17 THEN 1 ELSE 0 END) AS day_flights,
			SUM(CASE WHEN EXTRACT(HOUR FROM m.atd) BETWEEN 18 AND 23 THEN 1 ELSE 0 END) AS evening_flights,
			SUM(CASE WHEN EXTRACT(HOUR FROM m.atd) BETWEEN 0 AND 5 THEN 1 ELSE 0 END) AS night_flights
		FROM messages m
		JOIN district_shapes ds ON m.region = ds.gid
		WHERE m.ata IS NOT NULL
	`
	args := []interface{}{}
	if regID != 0 {
		query += " AND m.region = $1"
		args = append(args, regID)
	}
	if year != 0 {
		paramPos := len(args) + 1
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM m.dof) = $%d", paramPos)
		args = append(args, year)
	}
	query += " GROUP BY m.region, ds.name"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query flight times: %w", err)
	}
	defer rows.Close()

	var results []struct {
		RegionCode     int
		RegionName     string
		MorningFlights int
		DayFlights     int
		EveningFlights int
		NightFlights   int
	}
	for rows.Next() {
		var r struct {
			RegionCode     int
			RegionName     string
			MorningFlights int
			DayFlights     int
			EveningFlights int
			NightFlights   int
		}
		if err := rows.Scan(&r.RegionCode, &r.RegionName, &r.MorningFlights, &r.DayFlights, &r.EveningFlights, &r.NightFlights); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, r)
	}
	return results, nil
}

func (r *Repository) GetZeroFlightDays(ctx context.Context, regID int, year int) ([]struct {
	RegionCode     int
	RegionName     string
	ZeroFlightDays int
}, error) {
	query := `
		WITH flight_days AS (
			SELECT
				m.region AS region_code,
				ds.name AS region_name,
				COUNT(DISTINCT m.dof) AS flight_days
			FROM messages m
			JOIN district_shapes ds ON m.region = ds.gid
			WHERE m.ata IS NOT NULL
	`
	args := []interface{}{}
	if regID != 0 {
		query += " AND m.region = $1"
		args = append(args, regID)
	}
	if year != 0 {
		query += " AND EXTRACT(YEAR FROM m.dof) = $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, year)
	}
	query += `
			GROUP BY m.region, ds.name
		),
		total_days AS (
			SELECT
				ds.gid AS region_code,
				ds.name AS region_name,
				365 AS total_days_in_year
			FROM district_shapes ds
		)
		SELECT
			td.region_code,
			td.region_name,
			td.total_days_in_year - COALESCE(fd.flight_days, 0) AS zero_flight_days
		FROM total_days td
		LEFT JOIN flight_days fd ON td.region_code = fd.region_code AND td.region_name = fd.region_name
	`

	var rows pgx.Rows
	var err error
	if regID != 0 {
		rows, err = r.db.Query(ctx, query, regID, year)
	} else {
		rows, err = r.db.Query(ctx, query, year)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query zero flight days: %w", err)
	}
	defer rows.Close()

	var results []struct {
		RegionCode     int
		RegionName     string
		ZeroFlightDays int
	}
	for rows.Next() {
		var r struct {
			RegionCode     int
			RegionName     string
			ZeroFlightDays int
		}
		if err := rows.Scan(&r.RegionCode, &r.RegionName, &r.ZeroFlightDays); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if regID != 0 && r.RegionCode == regID {
			results = append(results, r)
		} else if regID == 0 {
			results = append(results, r)

		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return results, nil
}
func (r *Repository) GetTotalDistance(ctx context.Context, regID int, year int) ([]struct {
	RegionCode      int
	RegionName      string
	TotalDistanceKm float64
}, error) {
	query := `
		WITH all_coordinates AS (
			SELECT
				m.sid,
				m.region AS region_code,
				ds.name AS region_name,
				m.dep_coordinate AS coordinate,
				0 AS coord_order,
				m.dof
			FROM messages m
			JOIN district_shapes ds ON m.region = ds.gid
			WHERE m.ata IS NOT NULL AND m.arr_coordinate IS NOT NULL

			UNION ALL

			SELECT
				fc.sid,
				m.region AS region_code,
				ds.name AS region_name,
				fc.coordinate,
				fc.id AS coord_order,
				m.dof
			FROM flight_coordinates fc
			JOIN messages m ON fc.sid = m.sid
			JOIN district_shapes ds ON m.region = ds.gid

			UNION ALL

			SELECT
				m.sid,
				m.region AS region_code,
				ds.name AS region_name,
				m.arr_coordinate AS coordinate,
				999999999 AS coord_order,
				m.dof
			FROM messages m
			JOIN district_shapes ds ON m.region = ds.gid
			WHERE m.ata IS NOT NULL AND m.arr_coordinate IS NOT NULL
		),
		coordinate_pairs AS (
			SELECT
				ac1.sid,
				ac1.region_code,
				ac1.region_name,
				ac1.coordinate AS start_coord,
				LEAD(ac1.coordinate) OVER (PARTITION BY ac1.sid ORDER BY ac1.coord_order) AS end_coord,
				ac1.dof
			FROM all_coordinates ac1
		),
		flight_distances AS (
			SELECT
				sid,
				region_code,
				region_name,
				SUM(
					ST_Distance(
						geography(start_coord),
						geography(end_coord)
					) / 1000
				) AS total_distance_km,
				MIN(dof) AS flight_date -- берем дату для фильтрации по году
			FROM coordinate_pairs
			WHERE end_coord IS NOT NULL
			GROUP BY sid, region_code, region_name
		)
		SELECT
			region_code,
			region_name,
			SUM(total_distance_km) AS total_distance_km
		FROM flight_distances
		WHERE 1=1
	`

	args := []interface{}{}
	if regID != 0 {
		query += " AND region_code = $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, regID)
	}
	if year != 0 {
		query += " AND EXTRACT(YEAR FROM flight_date) = $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, year)
	}

	query += " GROUP BY region_code, region_name"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query total distance: %w", err)
	}
	defer rows.Close()

	var results []struct {
		RegionCode      int
		RegionName      string
		TotalDistanceKm float64
	}

	for rows.Next() {
		var r struct {
			RegionCode      int
			RegionName      string
			TotalDistanceKm float64
		}
		if err := rows.Scan(&r.RegionCode, &r.RegionName, &r.TotalDistanceKm); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return results, nil
}
func (r *Repository) UpdateMetrics(ctx context.Context, metrics chan *model.Metrics) error {
	query := `
		INSERT INTO flight_metrics (
			region_code, region_name, total_flight, avg_duration_minutes,
			total_distance_km, peak_load, avg_daily_flights, median_daily_flights,
			monthly_growth, flight_density, morning_flights, day_flights,
			evening_flights, night_flights, zero_flight_days, date
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16
		)
		ON CONFLICT (region_code,date) 
		DO UPDATE SET
			region_name = EXCLUDED.region_name,
			total_flight = EXCLUDED.total_flight,
			avg_duration_minutes = EXCLUDED.avg_duration_minutes,
			total_distance_km = EXCLUDED.total_distance_km,
			peak_load = EXCLUDED.peak_load,
			avg_daily_flights = EXCLUDED.avg_daily_flights,
			median_daily_flights = EXCLUDED.median_daily_flights,
			monthly_growth = EXCLUDED.monthly_growth,
			flight_density = EXCLUDED.flight_density,
			morning_flights = EXCLUDED.morning_flights,
			day_flights = EXCLUDED.day_flights,
			evening_flights = EXCLUDED.evening_flights,
			night_flights = EXCLUDED.night_flights,
			zero_flight_days = EXCLUDED.zero_flight_days,
			date = EXCLUDED.date
	`

	for m := range metrics {
		if m.RegionName == "" {
			continue
		}
		jsonData, err := json.Marshal(m.MonthlyGrowth)
		if err != nil {
			return fmt.Errorf("failed to marshal monthly growth: %w", err)
		}
		if jsonData == nil || m.RegionId == 17 {
			fmt.Println("a")
		}
		_, err = r.db.Exec(ctx, query,
			m.RegionId, m.RegionName,
			m.TotalFlight, m.AvgDurationMinutes,
			m.TotalDistance, m.PeakLoad,
			m.AvgDailyFlights, m.MedianDailyFlights,
			jsonData, m.FlightDensity,
			m.MorningFlights, m.DayFlights,
			m.EveningFlights, m.NightFlights,
			m.ZeroFlightDays, m.Year,
		)
		if err != nil {
			slog.Info("error upserting flight metrics: ", err)
			continue
		}
	}
	return nil
}

func (r *Repository) TotalFlightAndAVGDurationAllRussia(ctx context.Context, year int) (int, float32, error) {
	query := `
		SELECT
			COUNT(DISTINCT m.sid) AS total_flight,
			AVG(EXTRACT(EPOCH FROM (
				CASE
					WHEN m.atd > m.ata THEN (m.dof + INTERVAL '1 day' + m.ata) - (m.dof + m.atd)
					ELSE (m.dof + m.ata) - (m.dof + m.atd)
				END
			)) / 60) AS avg_duration_minutes
		FROM messages m
		WHERE m.ata IS NOT NULL AND m.arr_coordinate IS NOT NULL
	`
	args := []interface{}{}
	if year != 0 {
		query += " AND EXTRACT(YEAR FROM m.dof) = $1"
		args = append(args, year)
	}

	row := r.db.QueryRow(ctx, query, args...)
	var totalFlight int
	var avgDuration float32
	if err := row.Scan(&totalFlight, &avgDuration); err != nil {
		return 0, 0, fmt.Errorf("failed to query total flight and avg duration: %w", err)
	}
	return totalFlight, avgDuration, nil
}

func (r *Repository) GetDailyFlightMetricsAllRussia(ctx context.Context, year int) (float64, float64, error) {
	query := `
		WITH daily_flights AS (
			SELECT
				CASE WHEN m.atd > m.ata THEN m.dof + INTERVAL '1 day' ELSE m.dof END AS flight_date,
				COUNT(m.sid) AS daily_flight_count
			FROM messages m
			WHERE m.ata IS NOT NULL
	`
	args := []interface{}{}
	if year != 0 {
		query += " AND EXTRACT(YEAR FROM m.dof) = $1"
		args = append(args, year)
	}
	query += `
			GROUP BY CASE WHEN m.atd > m.ata THEN m.dof + INTERVAL '1 day' ELSE m.dof END
		)
		SELECT
			AVG(daily_flight_count) AS avg_daily_flights,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY daily_flight_count) AS median_daily_flights
		FROM daily_flights
	`

	row := r.db.QueryRow(ctx, query, args...)
	var avgDaily, medianDaily float64
	if err := row.Scan(&avgDaily, &medianDaily); err != nil {
		return 0, 0, fmt.Errorf("failed to query daily flight metrics: %w", err)
	}
	return avgDaily, medianDaily, nil
}

func (r *Repository) GetRussiaMonthlyGrowth(ctx context.Context, year int) (map[int]float64, error) {
	query := `
WITH monthly_counts AS (
    SELECT
        DATE_TRUNC('month', m.dof) AS month_date,
        COUNT(m.sid) AS monthly_flight_count
    FROM messages m
    WHERE m.ata IS NOT NULL AND EXTRACT(YEAR FROM m.dof) = $1
    GROUP BY DATE_TRUNC('month', m.dof)
),
growth_calc AS (
    SELECT
        month_date,
        monthly_flight_count,
        LAG(monthly_flight_count) OVER (ORDER BY month_date) AS prev_month_count
    FROM monthly_counts
),
growth_data AS (
    SELECT
        EXTRACT(MONTH FROM month_date)::INT AS month_number,
        CASE
            WHEN prev_month_count > 0
                THEN (monthly_flight_count::NUMERIC - prev_month_count) / prev_month_count * 100
            ELSE 0
        END AS growth_percentage,
        ROW_NUMBER() OVER (ORDER BY month_date DESC) AS rn
    FROM growth_calc
),
growth_json AS (
    SELECT json_object_agg(month_number, growth_percentage) AS monthly_growth_json
    FROM growth_data
)
SELECT monthly_growth_json
FROM growth_json
WHERE EXISTS (
    SELECT 1
    FROM growth_data gd
    WHERE gd.rn = 1
);
`

	var monthlyGrowth map[int]float64
	err := r.db.QueryRow(ctx, query, year).Scan(&monthlyGrowth)
	if err != nil {
		return nil, fmt.Errorf("failed to query Russia monthly growth for year %d: %w", year, err)
	}

	return monthlyGrowth, nil
}

func (r *Repository) GetFlightDensityAllRussia(ctx context.Context, year int) (float64, error) {
	query := `
		SELECT
			COUNT(DISTINCT m.sid)::NUMERIC / SUM(ds.area_km2 / 1000) AS flight_density
		FROM messages m
		JOIN district_shapes ds ON m.region = ds.gid
		WHERE m.ata IS NOT NULL AND m.arr_coordinate IS NOT NULL
	`
	args := []interface{}{}
	if year != 0 {
		query += " AND EXTRACT(YEAR FROM m.dof) = $1"
		args = append(args, year)
	}

	row := r.db.QueryRow(ctx, query, args...)
	var density float64
	if err := row.Scan(&density); err != nil {
		return 0, fmt.Errorf("failed to query flight density: %w", err)
	}
	return density, nil
}

func (r *Repository) GetFlightTimesAllRussia(ctx context.Context, year int) (int, int, int, int, error) {
	query := `
		SELECT
			SUM(CASE WHEN EXTRACT(HOUR FROM m.atd) BETWEEN 6 AND 11 THEN 1 ELSE 0 END) AS morning_flights,
			SUM(CASE WHEN EXTRACT(HOUR FROM m.atd) BETWEEN 12 AND 17 THEN 1 ELSE 0 END) AS day_flights,
			SUM(CASE WHEN EXTRACT(HOUR FROM m.atd) BETWEEN 18 AND 23 THEN 1 ELSE 0 END) AS evening_flights,
			SUM(CASE WHEN EXTRACT(HOUR FROM m.atd) BETWEEN 0 AND 5 THEN 1 ELSE 0 END) AS night_flights
		FROM messages m
		WHERE m.ata IS NOT NULL
	`
	args := []interface{}{}
	if year != 0 {
		query += " AND EXTRACT(YEAR FROM m.dof) = $1"
		args = append(args, year)
	}

	row := r.db.QueryRow(ctx, query, args...)
	var morning, day, evening, night int
	if err := row.Scan(&morning, &day, &evening, &night); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to query flight times: %w", err)
	}
	return morning, day, evening, night, nil
}

func (r *Repository) GetZeroFlightDaysAllRussia(ctx context.Context, year int) (int, error) {
	query := `
		SELECT
			365 - COUNT(DISTINCT m.dof) AS zero_flight_days
		FROM messages m
	`
	args := []interface{}{}
	if year != 0 {
		query += " WHERE EXTRACT(YEAR FROM m.dof) = $1"
		args = append(args, year)
	}

	row := r.db.QueryRow(ctx, query, args...)
	var zeroDays int
	if err := row.Scan(&zeroDays); err != nil {
		return 0, fmt.Errorf("failed to query zero flight days: %w", err)
	}
	return zeroDays, nil
}

func (r *Repository) GetTotalDistanceAllRussia(ctx context.Context, year int) (float64, error) {
	query := `
		WITH all_coordinates AS (SELECT m.sid, m.dep_coordinate AS coordinate, 0 AS coord_order, m.dof
								 FROM messages m
								 WHERE m.ata IS NOT NULL
								   AND m.arr_coordinate IS NOT NULL
								 UNION ALL
								 SELECT fc.sid, fc.coordinate, fc.id AS coord_order, m.dof
								 FROM flight_coordinates fc
										  JOIN messages m ON fc.sid = m.sid
								 UNION ALL
								 SELECT m.sid, m.arr_coordinate, 999999999 AS coord_order, m.dof
								 FROM messages m
								 WHERE m.ata IS NOT NULL
								   AND m.arr_coordinate IS NOT NULL),
			coordinate_pairs AS (
			SELECT
			ac1.sid,
			ac1.coordinate AS start_coord,
			LEAD(ac1.coordinate) OVER (PARTITION BY ac1.sid ORDER BY ac1.coord_order) AS end_coord,
			ac1.dof
			FROM all_coordinates ac1
			)
		SELECT SUM(ST_Distance(geography(start_coord), geography(end_coord)) / 1000) AS total_distance_km
		FROM coordinate_pairs
		WHERE end_coord IS NOT NULL
		  

	`
	args := []interface{}{}
	if year != 0 {
		query += " AND EXTRACT(YEAR FROM dof) = $1"
		args = append(args, year)
	}

	row := r.db.QueryRow(ctx, query, args...)
	var totalDistance float64
	if err := row.Scan(&totalDistance); err != nil {
		return 0, fmt.Errorf("failed to query total distance: %w", err)
	}
	return totalDistance, nil
}
func (r *Repository) GetPeakLoadAllRussia(ctx context.Context, year int) (int, error) {
	query := `
		WITH hourly_load AS (
    SELECT
        DATE_TRUNC('hour', m.dof + m.atd) AS hour,
        COUNT(*) AS hourly_count
    FROM messages m
    WHERE m.ata IS NOT NULL


	`
	args := []interface{}{}
	if year != 0 {
		query += " AND EXTRACT(YEAR FROM m.dof) = $1"
		args = append(args, year)
	}
	query += `
    GROUP BY DATE_TRUNC('hour', m.dof + m.atd)
		)
		SELECT
			hourly_count AS peak_load
		FROM hourly_load
		ORDER BY hourly_count DESC
		LIMIT 1;
	`

	row := r.db.QueryRow(ctx, query, args...)
	var peakLoad int
	if err := row.Scan(&peakLoad); err != nil {
		return 0, fmt.Errorf("failed to query peak load: %w", err)
	}
	return peakLoad, nil
}
