package signalr

type PrivateLinkResource struct {
	Id         *string                        `json:"id,omitempty"`
	Name       *string                        `json:"name,omitempty"`
	Properties *PrivateLinkResourceProperties `json:"properties,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
