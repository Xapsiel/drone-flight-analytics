package model

type File struct {
	Filename   string `json:"filename"`
	Size       int64  `json:"size"`
	Metadata   []byte `json:"metadata"`
	AuthorID   string `json:"author_id"`
	Status     string `json:"status"`
	ValidCount int    `json:"valid_count"`
	ErrorCount int    `json:"error_count"`
}
