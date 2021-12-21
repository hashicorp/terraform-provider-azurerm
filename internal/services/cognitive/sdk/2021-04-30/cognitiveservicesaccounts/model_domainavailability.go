package cognitiveservicesaccounts

type DomainAvailability struct {
	IsSubdomainAvailable *bool   `json:"isSubdomainAvailable,omitempty"`
	Reason               *string `json:"reason,omitempty"`
	SubdomainName        *string `json:"subdomainName,omitempty"`
	Type                 *string `json:"type,omitempty"`
}
