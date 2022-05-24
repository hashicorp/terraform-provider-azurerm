package confidentialledger

type AADBasedSecurityPrincipal struct {
	LedgerRoleName *LedgerRoleName `json:"ledgerRoleName,omitempty"`
	PrincipalId    *string         `json:"principalId,omitempty"`
	TenantId       *string         `json:"tenantId,omitempty"`
}
