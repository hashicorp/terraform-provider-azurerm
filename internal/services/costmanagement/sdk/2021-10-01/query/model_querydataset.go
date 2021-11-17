package query

type QueryDataset struct {
	Aggregation   *map[string]QueryAggregation `json:"aggregation,omitempty"`
	Configuration *QueryDatasetConfiguration   `json:"configuration,omitempty"`
	Filter        *QueryFilter                 `json:"filter,omitempty"`
	Granularity   *GranularityType             `json:"granularity,omitempty"`
	Grouping      *[]QueryGrouping             `json:"grouping,omitempty"`
}
