package connections

import "strings"

type LinkState string

const (
	LinkStateAuthenticated   LinkState = "Authenticated"
	LinkStateError           LinkState = "Error"
	LinkStateUnauthenticated LinkState = "Unauthenticated"
)

func PossibleValuesForLinkState() []string {
	return []string{
		string(LinkStateAuthenticated),
		string(LinkStateError),
		string(LinkStateUnauthenticated),
	}
}

func parseLinkState(input string) (*LinkState, error) {
	vals := map[string]LinkState{
		"authenticated":   LinkStateAuthenticated,
		"error":           LinkStateError,
		"unauthenticated": LinkStateUnauthenticated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LinkState(input)
	return &out, nil
}
