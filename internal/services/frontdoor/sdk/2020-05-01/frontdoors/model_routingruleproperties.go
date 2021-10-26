package frontdoors

type RoutingRuleProperties struct {
	AcceptedProtocols                *[]FrontDoorProtocol                                         `json:"acceptedProtocols,omitempty"`
	EnabledState                     *RoutingRuleEnabledState                                     `json:"enabledState,omitempty"`
	FrontendEndpoints                *[]SubResource                                               `json:"frontendEndpoints,omitempty"`
	PatternsToMatch                  *[]string                                                    `json:"patternsToMatch,omitempty"`
	ResourceState                    *FrontDoorResourceState                                      `json:"resourceState,omitempty"`
	RouteConfiguration               *RouteConfiguration                                          `json:"routeConfiguration,omitempty"`
	RulesEngine                      *SubResource                                                 `json:"rulesEngine,omitempty"`
	WebApplicationFirewallPolicyLink *RoutingRuleUpdateParametersWebApplicationFirewallPolicyLink `json:"webApplicationFirewallPolicyLink,omitempty"`
}
