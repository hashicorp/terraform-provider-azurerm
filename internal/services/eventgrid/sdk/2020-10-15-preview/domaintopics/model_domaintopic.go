package domaintopics

type DomainTopic struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties *DomainTopicProperties `json:"properties,omitempty"`
	SystemData *SystemData            `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}
