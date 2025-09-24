package model

type Table struct {
	Header     Header     `json:"header"`
	HeaderRows HeaderRows `json:"header_rows"`
	Rows       Rows       `json:"rows"`
}

type Header struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type HeaderCell struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CanSort     bool   `json:"can_sort"`
}

type HeaderRows []HeaderCell
type Cell struct {
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
}

type Row []Cell
type Rows []Row
