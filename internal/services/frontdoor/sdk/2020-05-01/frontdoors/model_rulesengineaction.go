package frontdoors

type RulesEngineAction struct {
	RequestHeaderActions       *[]HeaderAction     `json:"requestHeaderActions,omitempty"`
	ResponseHeaderActions      *[]HeaderAction     `json:"responseHeaderActions,omitempty"`
	RouteConfigurationOverride *RouteConfiguration `json:"routeConfigurationOverride,omitempty"`
}
