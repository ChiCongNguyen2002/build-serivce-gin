package resp

type Paging struct {
	Offset     int64       `json:"offset"`
	Limit      int64       `json:"limit"`
	Total      int64       `json:"total,omitempty"`
	NextCursor interface{} `json:"nextCursor,omitempty"`
}
