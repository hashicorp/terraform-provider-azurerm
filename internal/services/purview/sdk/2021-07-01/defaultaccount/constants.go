package defaultaccount

import "strings"

type ScopeType string

const (
	ScopeTypeSubscription ScopeType = "Subscription"
	ScopeTypeTenant       ScopeType = "Tenant"
)

func PossibleValuesForScopeType() []string {
	return []string{
		string(ScopeTypeSubscription),
		string(ScopeTypeTenant),
	}
}

func parseScopeType(input string) (*ScopeType, error) {
	vals := map[string]ScopeType{
		"subscription": ScopeTypeSubscription,
		"tenant":       ScopeTypeTenant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScopeType(input)
	return &out, nil
}
