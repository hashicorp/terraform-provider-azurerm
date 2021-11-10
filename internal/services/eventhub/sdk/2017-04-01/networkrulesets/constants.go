package networkrulesets

import "strings"

type DefaultAction string

const (
	DefaultActionAllow DefaultAction = "Allow"
	DefaultActionDeny  DefaultAction = "Deny"
)

func PossibleValuesForDefaultAction() []string {
	return []string{
		"Allow",
		"Deny",
	}
}

func parseDefaultAction(input string) (*DefaultAction, error) {
	vals := map[string]DefaultAction{
		"allow": "Allow",
		"deny":  "Deny",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := DefaultAction(v)
	return &out, nil
}

type NetworkRuleIPAction string

const (
	NetworkRuleIPActionAllow NetworkRuleIPAction = "Allow"
)

func PossibleValuesForNetworkRuleIPAction() []string {
	return []string{
		"Allow",
	}
}

func parseNetworkRuleIPAction(input string) (*NetworkRuleIPAction, error) {
	vals := map[string]NetworkRuleIPAction{
		"allow": "Allow",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := NetworkRuleIPAction(v)
	return &out, nil
}
