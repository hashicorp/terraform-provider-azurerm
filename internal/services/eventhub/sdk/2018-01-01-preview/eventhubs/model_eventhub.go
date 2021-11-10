package eventhubs

type Eventhub struct {
	Id         *string             `json:"id,omitempty"`
	Name       *string             `json:"name,omitempty"`
	Properties *EventhubProperties `json:"properties,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
