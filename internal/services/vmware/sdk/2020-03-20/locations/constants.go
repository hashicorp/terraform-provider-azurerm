package locations

import "strings"

type QuotaEnabled string

const (
	QuotaEnabledDisabled QuotaEnabled = "Disabled"
	QuotaEnabledEnabled  QuotaEnabled = "Enabled"
)

func PossibleValuesForQuotaEnabled() []string {
	return []string{
		"Disabled",
		"Enabled",
	}
}

func parseQuotaEnabled(input string) (*QuotaEnabled, error) {
	vals := map[string]QuotaEnabled{
		"disabled": "Disabled",
		"enabled":  "Enabled",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := QuotaEnabled(v)
	return &out, nil
}

type TrialStatus string

const (
	TrialStatusTrialAvailable TrialStatus = "TrialAvailable"
	TrialStatusTrialDisabled  TrialStatus = "TrialDisabled"
	TrialStatusTrialUsed      TrialStatus = "TrialUsed"
)

func PossibleValuesForTrialStatus() []string {
	return []string{
		"TrialAvailable",
		"TrialDisabled",
		"TrialUsed",
	}
}

func parseTrialStatus(input string) (*TrialStatus, error) {
	vals := map[string]TrialStatus{
		"trialavailable": "TrialAvailable",
		"trialdisabled":  "TrialDisabled",
		"trialused":      "TrialUsed",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := TrialStatus(v)
	return &out, nil
}
