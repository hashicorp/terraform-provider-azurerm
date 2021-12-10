package managedidentity

type SystemAssignedIdentityProperties struct {
	ClientId        *string `json:"clientId,omitempty"`
	ClientSecretUrl *string `json:"clientSecretUrl,omitempty"`
	PrincipalId     *string `json:"principalId,omitempty"`
	TenantId        *string `json:"tenantId,omitempty"`
}
