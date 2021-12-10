package workspaces

type ManagedIdentityConfiguration struct {
	PrincipalId *string `json:"principalId,omitempty"`
	TenantId    *string `json:"tenantId,omitempty"`
	Type        *string `json:"type,omitempty"`
}
