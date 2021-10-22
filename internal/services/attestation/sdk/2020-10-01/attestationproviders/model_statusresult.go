package attestationproviders

type StatusResult struct {
	AttestUri                  *string                      `json:"attestUri,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	Status                     *AttestationServiceStatus    `json:"status,omitempty"`
	TrustModel                 *string                      `json:"trustModel,omitempty"`
}
