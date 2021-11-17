package exports

type ExportExecution struct {
	ETag       *string                    `json:"eTag,omitempty"`
	Id         *string                    `json:"id,omitempty"`
	Name       *string                    `json:"name,omitempty"`
	Properties *ExportExecutionProperties `json:"properties,omitempty"`
	Type       *string                    `json:"type,omitempty"`
}
