package namespaces

type Identity struct {
	PrincipalId            *string                         `json:"principalId,omitempty"`
	TenantId               *string                         `json:"tenantId,omitempty"`
	Type                   *ManagedServiceIdentityType     `json:"type,omitempty"`
	UserAssignedIdentities *UserAssignedIdentityProperties `json:"userAssignedIdentities,omitempty"`
}
