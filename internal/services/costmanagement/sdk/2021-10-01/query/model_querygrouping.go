package query

type QueryGrouping struct {
	Name string          `json:"name"`
	Type QueryColumnType `json:"type"`
}
