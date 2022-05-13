package confidentialledger

type CertBasedSecurityPrincipal struct {
	Cert           *string         `json:"cert,omitempty"`
	LedgerRoleName *LedgerRoleName `json:"ledgerRoleName,omitempty"`
}
