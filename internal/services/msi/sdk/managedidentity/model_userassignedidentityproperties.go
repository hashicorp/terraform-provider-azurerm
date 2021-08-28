package managedidentity

type UserAssignedIdentityProperties struct {
	ClientId    *string `json:"clientId,omitempty"`
	PrincipalId *string `json:"principalId,omitempty"`
	TenantId    *string `json:"tenantId,omitempty"`
}
