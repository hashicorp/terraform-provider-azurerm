package hcxenterprisesites

import "strings"

type HcxEnterpriseSiteStatus string

const (
	HcxEnterpriseSiteStatusAvailable   HcxEnterpriseSiteStatus = "Available"
	HcxEnterpriseSiteStatusConsumed    HcxEnterpriseSiteStatus = "Consumed"
	HcxEnterpriseSiteStatusDeactivated HcxEnterpriseSiteStatus = "Deactivated"
	HcxEnterpriseSiteStatusDeleted     HcxEnterpriseSiteStatus = "Deleted"
)

func PossibleValuesForHcxEnterpriseSiteStatus() []string {
	return []string{
		string(HcxEnterpriseSiteStatusAvailable),
		string(HcxEnterpriseSiteStatusConsumed),
		string(HcxEnterpriseSiteStatusDeactivated),
		string(HcxEnterpriseSiteStatusDeleted),
	}
}

func parseHcxEnterpriseSiteStatus(input string) (*HcxEnterpriseSiteStatus, error) {
	vals := map[string]HcxEnterpriseSiteStatus{
		"available":   HcxEnterpriseSiteStatusAvailable,
		"consumed":    HcxEnterpriseSiteStatusConsumed,
		"deactivated": HcxEnterpriseSiteStatusDeactivated,
		"deleted":     HcxEnterpriseSiteStatusDeleted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HcxEnterpriseSiteStatus(input)
	return &out, nil
}
