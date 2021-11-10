package wcfrelays

type WcfRelay struct {
	Id         *string             `json:"id,omitempty"`
	Name       *string             `json:"name,omitempty"`
	Properties *WcfRelayProperties `json:"properties,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
