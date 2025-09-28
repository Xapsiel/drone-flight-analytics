package repository

import (
	"context"
	"encoding/json"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

// GetDistrictGeoJSON возвращает FeatureCollection с одной фичей для района по gid
func (r *Repository) GetDistrictGeoJSON(ctx context.Context, id int) ([]byte, error) {
	// Простой упрощающий конвейер на стороне БД; при необходимости параметризуйте толерантность
	query := `
WITH features AS (
    SELECT jsonb_build_object(
        'type','Feature',
        'geometry', ST_AsGeoJSON(geom, 6, 2)::jsonb,
        'properties', jsonb_build_object(
            'gid', gid,
            'district', name,
            'name_ru', name_ru,
            'name_en', name_en,
            'boundary', boundary,
            'admin_level', admin_leve,
            'timezone', timezone
        )
    ) AS feature
    FROM district_shapes
    WHERE gid = $1
)
SELECT jsonb_build_object('type','FeatureCollection','features', jsonb_agg(feature))
FROM features
`
	var raw json.RawMessage
	if err := r.db.QueryRow(ctx, query, id).Scan(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// GetAllDistrictsGeoJSONHandler возвращает FeatureCollection по всем районам в сыром JSON
func (r *Repository) GetAllDistrictsGeoJSONHandler(ctx context.Context) ([]byte, error) {
	query := `
WITH features AS (
    SELECT jsonb_build_object(
        'type','Feature',
        'geometry', ST_AsGeoJSON(geom, 6, 2)::jsonb,
        'properties', jsonb_build_object(
            'gid', gid,
            'district', name,
            'name_ru', name_ru,
            'name_en', name_en,
            'boundary', boundary,
            'admin_level', admin_leve,
            'timezone', timezone
        )
    ) AS feature
    FROM district_shapes
)
SELECT jsonb_build_object('type','FeatureCollection','features', jsonb_agg(feature))
FROM features
`
	var raw json.RawMessage
	if err := r.db.QueryRow(ctx, query).Scan(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// GetDistrictsMVT возвращает векторный тайл MVT (protobuf) по z/x/y
func (r *Repository) GetDistrictsMVT(ctx context.Context, z, x, y int) ([]byte, error) {
	// Используем ST_TileEnvelope для bbox тайла в 3857 и строим MVT
	// Предполагаем, что district_shapes.geom в SRID 4326
	query := `
WITH tile AS (
    SELECT ST_TileEnvelope($1, $2, $3) AS geom
)
SELECT ST_AsMVT(mvtq, 'districts', 4096, 'mvt_geom')
FROM (
    SELECT 
        gid,
        name_ru,
        name_en,
        admin_leve,
        timezone,
        ST_AsMVTGeom(
            ST_Transform(
                CASE WHEN ST_SRID(geom) = 0 THEN ST_SetSRID(geom, 4326) ELSE geom END,
                3857
            ),
            (SELECT geom FROM tile),
            4096, 64, true
        ) AS mvt_geom
    FROM district_shapes
    WHERE ST_Intersects(
        ST_Transform(
            CASE WHEN ST_SRID(geom) = 0 THEN ST_SetSRID(geom, 4326) ELSE geom END,
            3857
        ),
        (SELECT geom FROM tile)
    )
) AS mvtq;`

	var tileBytes []byte
	if err := r.db.QueryRow(ctx, query, z, x, y).Scan(&tileBytes); err != nil {
		return nil, err
	}
	return tileBytes, nil
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

// GetDistrictsMVT возвращает векторный тайл MVT (protobuf) по z/x/y
func (r *Repository) GetDistrictsMVT(ctx context.Context, z, x, y int) ([]byte, error) {
	query := `
WITH tile AS (
    SELECT ST_TileEnvelope($1, $2, $3) AS geom
)
SELECT ST_AsMVT(mvtq, 'districts', 4096, 'mvt_geom')
FROM (
    SELECT 
        gid,
        name_ru,
        name_en,
        admin_leve,
        timezone,
        ST_AsMVTGeom(
            ST_Transform(
                CASE WHEN ST_SRID(geom) = 0 THEN ST_SetSRID(geom, 4326) ELSE geom END,
                3857
            ),
            (SELECT geom FROM tile),
            4096, 64, true
        ) AS mvt_geom
    FROM district_shapes
    WHERE ST_Intersects(
        ST_Transform(
            CASE WHEN ST_SRID(geom) = 0 THEN ST_SetSRID(geom, 4326) ELSE geom END,
            3857
        ),
        (SELECT geom FROM tile)
    )
) AS mvtq;`

	var tileBytes []byte
	if err := r.db.QueryRow(ctx, query, z, x, y).Scan(&tileBytes); err != nil {
		return nil, err
	}
	return tileBytes, nil
}
