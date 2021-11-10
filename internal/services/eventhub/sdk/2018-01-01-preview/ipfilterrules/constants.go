package ipfilterrules

import "strings"

type IPAction string

const (
	IPActionAccept IPAction = "Accept"
	IPActionReject IPAction = "Reject"
)

func PossibleValuesForIPAction() []string {
	return []string{
		"Accept",
		"Reject",
	}
}

func parseIPAction(input string) (*IPAction, error) {
	vals := map[string]IPAction{
		"accept": "Accept",
		"reject": "Reject",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := IPAction(v)
	return &out, nil
}
