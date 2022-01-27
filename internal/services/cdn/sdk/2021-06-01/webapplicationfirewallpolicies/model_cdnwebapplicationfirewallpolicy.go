package webapplicationfirewallpolicies

type CdnWebApplicationFirewallPolicy struct {
	Etag       *string                                    `json:"etag,omitempty"`
	Id         *string                                    `json:"id,omitempty"`
	Location   string                                     `json:"location"`
	Name       *string                                    `json:"name,omitempty"`
	Properties *CdnWebApplicationFirewallPolicyProperties `json:"properties,omitempty"`
	Sku        Sku                                        `json:"sku"`
	SystemData *SystemData                                `json:"systemData,omitempty"`
	Tags       *map[string]string                         `json:"tags,omitempty"`
	Type       *string                                    `json:"type,omitempty"`
}
