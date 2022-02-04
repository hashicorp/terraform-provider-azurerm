package frontdoors

type RulesEngineProperties struct {
	ResourceState *FrontDoorResourceState `json:"resourceState,omitempty"`
	Rules         *[]RulesEngineRule      `json:"rules,omitempty"`
}
