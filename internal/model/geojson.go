package model

import (
	"github.com/paulmach/orb/geojson"
)

type DistrictGeoJSON struct {
	District District                   `json:"district"`
	Features *geojson.FeatureCollection `json:"features"`
}

func NewGeoJSON(district District, features *geojson.FeatureCollection) *DistrictGeoJSON {
	return &DistrictGeoJSON{
		District: district,
		Features: features,
	}
}
