package model

type Page struct {
	Domain      string
	Country     string
	TileServer  string
	Title       string
	Keywords    string
	Description string
	Year        int
	HasMap      bool

	Error string
}
