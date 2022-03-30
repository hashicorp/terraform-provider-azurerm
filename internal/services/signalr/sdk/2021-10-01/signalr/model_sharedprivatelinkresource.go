package signalr

type SharedPrivateLinkResource struct {
	Id         *string                              `json:"id,omitempty"`
	Name       *string                              `json:"name,omitempty"`
	Properties *SharedPrivateLinkResourceProperties `json:"properties,omitempty"`
	SystemData *SystemData                          `json:"systemData,omitempty"`
	Type       *string                              `json:"type,omitempty"`
}
