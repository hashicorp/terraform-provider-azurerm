package privatelinkresources

type PrivateLinkResourceDescription struct {
	Id         *string                        `json:"id,omitempty"`
	Name       *string                        `json:"name,omitempty"`
	Properties *PrivateLinkResourceProperties `json:"properties,omitempty"`
	SystemData *SystemData                    `json:"systemData,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
