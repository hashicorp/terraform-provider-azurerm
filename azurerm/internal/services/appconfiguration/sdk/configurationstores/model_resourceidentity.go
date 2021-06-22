package configurationstores

type ResourceIdentity struct {
	PrincipalId *string       `json:"principalId,omitempty"`
	TenantId    *string       `json:"tenantId,omitempty"`
	Type        *IdentityType `json:"type,omitempty"`
}
