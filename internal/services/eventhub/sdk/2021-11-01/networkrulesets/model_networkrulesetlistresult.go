package networkrulesets

type NetworkRuleSetListResult struct {
	NextLink *string           `json:"nextLink,omitempty"`
	Value    *[]NetworkRuleSet `json:"value,omitempty"`
}
