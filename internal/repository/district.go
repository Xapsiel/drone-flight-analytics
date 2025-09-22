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
				    timezone,geom   FROM district_shapes
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
	feature.Properties["district"] = name
	feature.Properties["name_ru"] = name
	feature.Properties["name_en"] = name
	feature.Properties["boundary"] = name
	feature.Properties["admin_level"] = name
	feature.Properties["timezone"] = name
	result.Features.Features = append(result.Features.Features, feature)
	return result, nil
}
