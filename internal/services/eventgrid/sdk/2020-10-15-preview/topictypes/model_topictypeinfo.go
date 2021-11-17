package topictypes

type TopicTypeInfo struct {
	Id         *string              `json:"id,omitempty"`
	Name       *string              `json:"name,omitempty"`
	Properties *TopicTypeProperties `json:"properties,omitempty"`
	Type       *string              `json:"type,omitempty"`
}
