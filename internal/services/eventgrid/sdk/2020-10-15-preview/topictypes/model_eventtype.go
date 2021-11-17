package topictypes

type EventType struct {
	Id         *string              `json:"id,omitempty"`
	Name       *string              `json:"name,omitempty"`
	Properties *EventTypeProperties `json:"properties,omitempty"`
	Type       *string              `json:"type,omitempty"`
}
