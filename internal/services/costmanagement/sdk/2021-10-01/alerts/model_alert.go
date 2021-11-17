package alerts

type Alert struct {
	ETag       *string          `json:"eTag,omitempty"`
	Id         *string          `json:"id,omitempty"`
	Name       *string          `json:"name,omitempty"`
	Properties *AlertProperties `json:"properties,omitempty"`
	Type       *string          `json:"type,omitempty"`
}
