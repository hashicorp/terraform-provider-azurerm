package eventhubs

type Destination struct {
	Name       *string                `json:"name,omitempty"`
	Properties *DestinationProperties `json:"properties,omitempty"`
}
