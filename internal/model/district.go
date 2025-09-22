package model

type District struct {
	Gid        int
	Name       string `json:"name"`
	TimeZone   string `json:"time_zone"`
	ISO3166_2  string `json:"iso_3166_2"`
	NameRU     string `json:"name_ru"`
	NameEn     string `json:"name_en"`
	AdminLevel int    `json:"admin_level"`
	Boundary   string `json:"boundary"`
}
