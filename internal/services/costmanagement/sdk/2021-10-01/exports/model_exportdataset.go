package exports

type ExportDataset struct {
	Configuration *ExportDatasetConfiguration `json:"configuration,omitempty"`
	Granularity   *GranularityType            `json:"granularity,omitempty"`
}
