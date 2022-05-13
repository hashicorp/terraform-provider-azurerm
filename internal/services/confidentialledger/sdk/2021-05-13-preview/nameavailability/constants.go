package nameavailability

import "strings"

type CheckNameAvailabilityReason string

const (
	CheckNameAvailabilityReasonAlreadyExists CheckNameAvailabilityReason = "AlreadyExists"
	CheckNameAvailabilityReasonInvalid       CheckNameAvailabilityReason = "Invalid"
)

func PossibleValuesForCheckNameAvailabilityReason() []string {
	return []string{
		string(CheckNameAvailabilityReasonAlreadyExists),
		string(CheckNameAvailabilityReasonInvalid),
	}
}

func parseCheckNameAvailabilityReason(input string) (*CheckNameAvailabilityReason, error) {
	vals := map[string]CheckNameAvailabilityReason{
		"alreadyexists": CheckNameAvailabilityReasonAlreadyExists,
		"invalid":       CheckNameAvailabilityReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CheckNameAvailabilityReason(input)
	return &out, nil
}
