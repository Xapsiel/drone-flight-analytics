
CREATE TABLE IF NOT EXISTS flight_metrics(
    id SERIAL PRIMARY KEY ,
    region_code INT  ,
    region_name varchar(80),
    total_flight INT DEFAULT 0,
    avg_duration_minutes DECIMAL(10,2) DEFAULT  0.00,
    total_distance_km   DECIMAL(10,2) DEFAULT  0.00,
    peak_load INT DEFAULT 0,
    avg_daily_flights DECIMAL(10,2) DEFAULT 0.00,
    median_daily_flights DECIMAL(10,2) DEFAULT 0.00,
    monthly_growth jsonb,
    flight_density DECIMAL(10,2),
    morning_flights INTEGER,
    day_flights INTEGER,
    evening_flights INTEGER,
    night_flights INTEGER,
    zero_flight_days INTEGER,
    date int,
    UNIQUE (region_code,date)
)