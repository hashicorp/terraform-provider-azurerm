package connections

type ConfirmConsentCodeDefinition struct {
	Code     *string `json:"code,omitempty"`
	ObjectId *string `json:"objectId,omitempty"`
	TenantId *string `json:"tenantId,omitempty"`
}
