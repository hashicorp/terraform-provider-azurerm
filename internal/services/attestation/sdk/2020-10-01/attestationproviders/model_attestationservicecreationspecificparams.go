package attestationproviders

type AttestationServiceCreationSpecificParams struct {
	PolicySigningCertificates *JsonWebKeySet `json:"policySigningCertificates,omitempty"`
}
