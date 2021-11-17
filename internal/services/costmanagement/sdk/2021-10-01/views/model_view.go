package views

type View struct {
	ETag       *string         `json:"eTag,omitempty"`
	Id         *string         `json:"id,omitempty"`
	Name       *string         `json:"name,omitempty"`
	Properties *ViewProperties `json:"properties,omitempty"`
	Type       *string         `json:"type,omitempty"`
}
