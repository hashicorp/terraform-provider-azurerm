package attestationproviders

type AttestationProviderListResult struct {
	SystemData *SystemData             `json:"systemData,omitempty"`
	Value      *[]AttestationProviders `json:"value,omitempty"`
}
