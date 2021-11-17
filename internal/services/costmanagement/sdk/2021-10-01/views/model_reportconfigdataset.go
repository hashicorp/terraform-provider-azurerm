package views

type ReportConfigDataset struct {
	Aggregation   *map[string]ReportConfigAggregation `json:"aggregation,omitempty"`
	Configuration *ReportConfigDatasetConfiguration   `json:"configuration,omitempty"`
	Filter        *ReportConfigFilter                 `json:"filter,omitempty"`
	Granularity   *ReportGranularityType              `json:"granularity,omitempty"`
	Grouping      *[]ReportConfigGrouping             `json:"grouping,omitempty"`
	Sorting       *[]ReportConfigSorting              `json:"sorting,omitempty"`
}
