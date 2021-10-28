package vaults

type PrivateEndpointConnectionItem struct {
	Etag       *string                              `json:"etag,omitempty"`
	Id         *string                              `json:"id,omitempty"`
	Properties *PrivateEndpointConnectionProperties `json:"properties,omitempty"`
}
