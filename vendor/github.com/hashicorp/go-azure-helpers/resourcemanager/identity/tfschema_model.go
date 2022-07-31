package identity

type ModelUserAssigned struct {
	Type        Type     `tfschema:"type"`
	IdentityIds []string `tfschema:"identity_ids"`
}

type ModelSystemAssigned struct {
	Type        Type   `tfschema:"type"`
	PrincipalId string `tfschema:"principal_id"`
	TenantId    string `tfschema:"tenant_id"`
}

type ModelSystemAssignedUserAssigned struct {
	Type        Type     `tfschema:"type"`
	PrincipalId string   `tfschema:"principal_id"`
	TenantId    string   `tfschema:"tenant_id"`
	IdentityIds []string `tfschema:"identity_ids"`
}
