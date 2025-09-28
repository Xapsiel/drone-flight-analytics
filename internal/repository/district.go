package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/paulmach/orb/geojson"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

func (r *Repository) GetDistrictGeoJSON(ctx context.Context, id int) (*model.DistrictGeoJSON, error) {
	query := `
				SELECT 
				    gid, name, name_ru, 
				    name_en, boundary, admin_leve, 
				    timezone,ST_AsGeoJSON(geom) as geom   FROM district_shapes
				WHERE gid=$1
			 `
	result := &model.DistrictGeoJSON{
		District: model.District{},
		Features: &geojson.FeatureCollection{},
	}
	var geomJSON string

	err := r.db.QueryRow(ctx, query, id).Scan(
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
			"district", id,
		)
		return nil, fmt.Errorf("failed to unmarshal geometry JSON for district %s: %w", id, err)
	}

	feature := geojson.NewFeature(geometry.Geometry())
	if result.District.Name != nil {
		feature.Properties["district"] = result.District.Name
	}
	if result.District.NameRU != nil {
		feature.Properties["name_ru"] = id
	}
	if result.District.NameEn != nil {
		feature.Properties["name_en"] = id
	}
	if result.District.Boundary != nil {
		feature.Properties["boundary"] = id
	}
	if result.District.AdminLevel != nil {
		feature.Properties["admin_level"] = id
	}
	if result.District.TimeZone != nil {
		feature.Properties["timezone"] = id
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

func (r *Repository) GetRegions(ctx context.Context) []model.District {
	res := []model.District{}
	query := `
			 SELECT gid,name_ru FROM district_shapes
			`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return res
	}
	defer rows.Close()
	for rows.Next() {
		var gid int
		var name string
		err = rows.Scan(&gid, &name)
		if err != nil {
			return res
		}
		res = append(res, model.District{
			Gid:  &gid,
			Name: &name,
		})
	}
	return res
}

func (r *Repository) GetFlightYears(ctx context.Context) []int {
	res := []int{}
	query := `
			SELECT EXTRACT(YEAR FROM dof) AS year
			FROM messages
				WHERE dof IS NOT NULL
			GROUP BY EXTRACT(YEAR FROM dof);
			`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return res
	}
	defer rows.Close()
	for rows.Next() {
		var gid int
		err = rows.Scan(&gid)
		if err != nil {
			return res
		}
		res = append(res, gid)
	}
	return res
}
