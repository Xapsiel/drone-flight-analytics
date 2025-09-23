package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/paulmach/orb/geojson"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

func (r *Repository) GetDistrictGeoJSON(ctx context.Context, name string) (*model.DistrictGeoJSON, error) {
	query := `
				SELECT 
				    gid, name, name_ru, 
				    name_en, boundary, admin_leve, 
				    timezone,ST_AsGeoJSON(geom) as geom   FROM district_shapes
				WHERE name=$1
			 `
	result := &model.DistrictGeoJSON{
		District: model.District{},
		Features: &geojson.FeatureCollection{},
	}
	var geomJSON string

	err := r.db.QueryRow(ctx, query, name).Scan(
		&result.District.Gid, &result.District.Name,
		&result.District.NameRU, &result.District.NameEn,
		&result.District.Boundary, &result.District.AdminLevel,
		&result.District.TimeZone, &geomJSON,
	)
	if err != nil {
		return nil, err

	}
	geometry, err := geojson.UnmarshalGeometry([]byte(geomJSON))
	if err != nil {
		slog.Error(
			"Failed to unmarshal geometry JSON",
			"error", err,
			"district", name,
		)
		return nil, fmt.Errorf("failed to unmarshal geometry JSON for district %s: %w", name, err)
	}

	feature := geojson.NewFeature(geometry.Geometry())
	if result.District.Name != nil {
		feature.Properties["district"] = result.District.Name
	}
	if result.District.NameRU != nil {
		feature.Properties["name_ru"] = name
	}
	if result.District.NameEn != nil {
		feature.Properties["name_en"] = name
	}
	if result.District.Boundary != nil {
		feature.Properties["boundary"] = name
	}
	if result.District.AdminLevel != nil {
		feature.Properties["admin_level"] = name
	}
	if result.District.TimeZone != nil {
		feature.Properties["timezone"] = name
	}
	result.Features.Features = append(result.Features.Features, feature)
	return result, nil
}
func (r *Repository) GetAllDistrictsGeoJSONHandler(ctx context.Context) ([]model.DistrictGeoJSON, error) {
	query := `
				SELECT 
				    gid, name, name_ru, 
				    name_en, boundary, admin_leve, 
				    timezone,ST_AsGeoJSON(geom) as geom   FROM district_shapes
			 `
	results := []model.DistrictGeoJSON{}
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err

	}
	defer rows.Close()
	for rows.Next() {
		result := &model.DistrictGeoJSON{
			District: model.District{},
			Features: &geojson.FeatureCollection{},
		}

		var geomJSON string
		err = rows.Scan(
			&result.District.Gid, &result.District.Name,
			&result.District.NameRU, &result.District.NameEn,
			&result.District.Boundary, &result.District.AdminLevel,
			&result.District.TimeZone, &geomJSON,
		)
		if err != nil {
			return nil, err

		}
		geometry, err := geojson.UnmarshalGeometry([]byte(geomJSON))
		if err != nil {
			return nil, err
		}

		feature := geojson.NewFeature(geometry.Geometry())
		if result.District.Name != nil {
			feature.Properties["district"] = result.District.Name
		}
		if result.District.NameRU != nil {

			feature.Properties["name_ru"] = result.District.NameRU
		}
		if result.District.NameEn != nil {
			feature.Properties["name_en"] = result.District.NameEn
		}
		if result.District.Boundary != nil {
			feature.Properties["boundary"] = result.District.Boundary
		}
		if result.District.AdminLevel != nil {
			feature.Properties["admin_level"] = result.District.AdminLevel
		}
		if result.District.TimeZone != nil {
			feature.Properties["timezone"] = result.District.TimeZone
		}
		result.Features.Features = append(result.Features.Features, feature)
		results = append(results, *result)
	}
	return results, nil
}
