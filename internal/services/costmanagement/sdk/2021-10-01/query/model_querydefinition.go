package query

type QueryDefinition struct {
	Dataset    QueryDataset     `json:"dataset"`
	TimePeriod *QueryTimePeriod `json:"timePeriod,omitempty"`
	Timeframe  TimeframeType    `json:"timeframe"`
	Type       ExportType       `json:"type"`
}
