package creators

type CreatorUpdateParameters struct {
	Properties *CreatorProperties `json:"properties,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
