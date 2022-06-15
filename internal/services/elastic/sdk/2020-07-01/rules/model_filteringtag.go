package rules

type FilteringTag struct {
	Action *TagAction `json:"action,omitempty"`
	Name   *string    `json:"name,omitempty"`
	Value  *string    `json:"value,omitempty"`
}
