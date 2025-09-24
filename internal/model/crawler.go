package model

import "github.com/paulmach/orb"

type ParsedMessage struct {
	Region      string      `json:"region"`        // Исходный регион (e.g. Ростовский)
	SID         string      `json:"sid"`           // Уникальный ID
	DOF         string      `json:"dof"`           // Нормализованная дата YYYY-MM-DD
	ATD         string      `json:"atd"`           // Время вылета чч:мм
	ATA         string      `json:"ata,omitempty"` // Время прибытия (если есть)
	DepCoords   string      `json:"dep_coords"`    // Норм. координаты вылета ddmmssNdddmmssE
	ArrCoords   string      `json:"arr_coords,omitempty"`
	DepLatLon   orb.Point   `json:"dep_latlon"` // Decimal [ lon,lat]
	ArrLatLon   orb.Point   `json:"arr_latlon,omitempty"`
	ArrRegionRF string      `json:"arr_region_rf,omitempty"`
	ZoneCoords  []string    `json:"zone_coords,omitempty"` // Промежуточные координаты зоны
	ZoneLatLon  []orb.Point `json:"zone_latlon,omitempty"` // Промежуточные координаты в decimal
	OPR         string      `json:"opr,omitempty"`         // Оператор
	REG         string      `json:"reg,omitempty"`         // Регистрация
	TYP         string      `json:"typ,omitempty"`         // Тип (BLA, AER etc.)
	RMK         string      `json:"rmk,omitempty"`         // Замечания
	MinAlt      int         `json:"min_alt"`               // Минимальная высота в метрах
	MaxAlt      int         `json:"max_alt"`               // Максимальная высота в метрах
}
