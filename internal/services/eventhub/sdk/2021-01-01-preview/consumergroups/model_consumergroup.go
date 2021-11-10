package consumergroups

type ConsumerGroup struct {
	Id         *string                  `json:"id,omitempty"`
	Name       *string                  `json:"name,omitempty"`
	Properties *ConsumerGroupProperties `json:"properties,omitempty"`
	SystemData *SystemData              `json:"systemData,omitempty"`
	Type       *string                  `json:"type,omitempty"`
}
