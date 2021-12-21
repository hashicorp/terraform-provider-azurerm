package accounts

type CreateVirtualNetworkRuleWithAccountParameters struct {
	Name       string                                     `json:"name"`
	Properties CreateOrUpdateVirtualNetworkRuleProperties `json:"properties"`
}
