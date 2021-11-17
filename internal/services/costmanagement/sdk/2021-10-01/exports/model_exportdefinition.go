package exports

type ExportDefinition struct {
	DataSet    *ExportDataset    `json:"dataSet,omitempty"`
	TimePeriod *ExportTimePeriod `json:"timePeriod,omitempty"`
	Timeframe  TimeframeType     `json:"timeframe"`
	Type       ExportType        `json:"type"`
}
