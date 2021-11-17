package forecast

type ForecastDefinition struct {
	Dataset                 ForecastDataset       `json:"dataset"`
	IncludeActualCost       *bool                 `json:"includeActualCost,omitempty"`
	IncludeFreshPartialCost *bool                 `json:"includeFreshPartialCost,omitempty"`
	TimePeriod              *QueryTimePeriod      `json:"timePeriod,omitempty"`
	Timeframe               ForecastTimeframeType `json:"timeframe"`
	Type                    ForecastType          `json:"type"`
}
