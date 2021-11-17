package topics

type ExtensionTopic struct {
	Id         *string                   `json:"id,omitempty"`
	Name       *string                   `json:"name,omitempty"`
	Properties *ExtensionTopicProperties `json:"properties,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}
