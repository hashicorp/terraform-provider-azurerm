package exports

type Export struct {
	ETag       *string           `json:"eTag,omitempty"`
	Id         *string           `json:"id,omitempty"`
	Name       *string           `json:"name,omitempty"`
	Properties *ExportProperties `json:"properties,omitempty"`
	Type       *string           `json:"type,omitempty"`
}
