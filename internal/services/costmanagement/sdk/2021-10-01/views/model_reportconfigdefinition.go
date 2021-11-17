package views

type ReportConfigDefinition struct {
	DataSet                   *ReportConfigDataset    `json:"dataSet,omitempty"`
	IncludeMonetaryCommitment *bool                   `json:"includeMonetaryCommitment,omitempty"`
	TimePeriod                *ReportConfigTimePeriod `json:"timePeriod,omitempty"`
	Timeframe                 ReportTimeframeType     `json:"timeframe"`
	Type                      ReportType              `json:"type"`
}
