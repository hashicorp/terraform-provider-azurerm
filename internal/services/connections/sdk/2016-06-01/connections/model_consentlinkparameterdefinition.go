package connections

type ConsentLinkParameterDefinition struct {
	ObjectId      *string `json:"objectId,omitempty"`
	ParameterName *string `json:"parameterName,omitempty"`
	RedirectUrl   *string `json:"redirectUrl,omitempty"`
	TenantId      *string `json:"tenantId,omitempty"`
}
