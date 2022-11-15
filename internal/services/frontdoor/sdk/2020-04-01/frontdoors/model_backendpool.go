package frontdoors

type BackendPool struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties *BackendPoolProperties `json:"properties,omitempty"`
	Type       *string                `json:"type,omitempty"`
}
