package systemtopics

type SystemTopicProperties struct {
	MetricResourceId  *string                    `json:"metricResourceId,omitempty"`
	ProvisioningState *ResourceProvisioningState `json:"provisioningState,omitempty"`
	Source            *string                    `json:"source,omitempty"`
	TopicType         *string                    `json:"topicType,omitempty"`
}
