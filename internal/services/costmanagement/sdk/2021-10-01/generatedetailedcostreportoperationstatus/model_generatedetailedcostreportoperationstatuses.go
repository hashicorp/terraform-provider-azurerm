package generatedetailedcostreportoperationstatus

type GenerateDetailedCostReportOperationStatuses struct {
	Error      *ErrorDetails `json:"error,omitempty"`
	Id         *string       `json:"id,omitempty"`
	Name       *string       `json:"name,omitempty"`
	Properties *DownloadURL  `json:"properties,omitempty"`
	Status     *Status       `json:"status,omitempty"`
	Type       *string       `json:"type,omitempty"`
}
