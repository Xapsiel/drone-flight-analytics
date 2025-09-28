package model

type File struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Metadata []byte `json:"metadata"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}
