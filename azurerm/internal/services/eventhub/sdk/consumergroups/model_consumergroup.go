package consumergroups

type ConsumerGroup struct {
	Id         *string                  `json:"id,omitempty"`
	Name       *string                  `json:"name,omitempty"`
	Properties *ConsumerGroupProperties `json:"properties,omitempty"`
	Type       *string                  `json:"type,omitempty"`
}
