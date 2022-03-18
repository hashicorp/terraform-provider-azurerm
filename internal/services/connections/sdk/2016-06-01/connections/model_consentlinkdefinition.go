package connections

type ConsentLinkDefinition struct {
	DisplayName        *string    `json:"displayName,omitempty"`
	FirstPartyLoginUri *string    `json:"firstPartyLoginUri,omitempty"`
	Link               *string    `json:"link,omitempty"`
	Status             *LinkState `json:"status,omitempty"`
}
