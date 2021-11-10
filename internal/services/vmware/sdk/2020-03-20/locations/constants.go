package locations

import "strings"

type QuotaEnabled string

const (
	QuotaEnabledDisabled QuotaEnabled = "Disabled"
	QuotaEnabledEnabled  QuotaEnabled = "Enabled"
)

func PossibleValuesForQuotaEnabled() []string {
	return []string{
		string(QuotaEnabledDisabled),
		string(QuotaEnabledEnabled),
	}
}

func parseQuotaEnabled(input string) (*QuotaEnabled, error) {
	vals := map[string]QuotaEnabled{
		"disabled": QuotaEnabledDisabled,
		"enabled":  QuotaEnabledEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QuotaEnabled(input)
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
		string(TrialStatusTrialAvailable),
		string(TrialStatusTrialDisabled),
		string(TrialStatusTrialUsed),
	}
}

func parseTrialStatus(input string) (*TrialStatus, error) {
	vals := map[string]TrialStatus{
		"trialavailable": TrialStatusTrialAvailable,
		"trialdisabled":  TrialStatusTrialDisabled,
		"trialused":      TrialStatusTrialUsed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrialStatus(input)
	return &out, nil
}
