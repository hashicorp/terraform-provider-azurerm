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
		"Available",
		"Consumed",
		"Deactivated",
		"Deleted",
	}
}

func parseHcxEnterpriseSiteStatus(input string) (*HcxEnterpriseSiteStatus, error) {
	vals := map[string]HcxEnterpriseSiteStatus{
		"available":   "Available",
		"consumed":    "Consumed",
		"deactivated": "Deactivated",
		"deleted":     "Deleted",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := HcxEnterpriseSiteStatus(v)
	return &out, nil
}
