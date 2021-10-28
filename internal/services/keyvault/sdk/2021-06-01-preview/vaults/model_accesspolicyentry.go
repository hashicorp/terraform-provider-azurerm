package vaults

type AccessPolicyEntry struct {
	ApplicationId *string     `json:"applicationId,omitempty"`
	ObjectId      string      `json:"objectId"`
	Permissions   Permissions `json:"permissions"`
	TenantId      string      `json:"tenantId"`
}
