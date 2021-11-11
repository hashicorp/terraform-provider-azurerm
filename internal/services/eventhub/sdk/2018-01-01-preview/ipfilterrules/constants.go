package ipfilterrules

import "strings"

type IPAction string

const (
	IPActionAccept IPAction = "Accept"
	IPActionReject IPAction = "Reject"
)

func PossibleValuesForIPAction() []string {
	return []string{
		string(IPActionAccept),
		string(IPActionReject),
	}
}

func parseIPAction(input string) (*IPAction, error) {
	vals := map[string]IPAction{
		"accept": IPActionAccept,
		"reject": IPActionReject,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPAction(input)
	return &out, nil
}
