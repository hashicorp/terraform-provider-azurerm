package authorizations

import "strings"

type ExpressRouteAuthorizationProvisioningState string

const (
	ExpressRouteAuthorizationProvisioningStateFailed    ExpressRouteAuthorizationProvisioningState = "Failed"
	ExpressRouteAuthorizationProvisioningStateSucceeded ExpressRouteAuthorizationProvisioningState = "Succeeded"
	ExpressRouteAuthorizationProvisioningStateUpdating  ExpressRouteAuthorizationProvisioningState = "Updating"
)

func PossibleValuesForExpressRouteAuthorizationProvisioningState() []string {
	return []string{
		"Failed",
		"Succeeded",
		"Updating",
	}
}

func parseExpressRouteAuthorizationProvisioningState(input string) (*ExpressRouteAuthorizationProvisioningState, error) {
	vals := map[string]ExpressRouteAuthorizationProvisioningState{
		"failed":    "Failed",
		"succeeded": "Succeeded",
		"updating":  "Updating",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := ExpressRouteAuthorizationProvisioningState(v)
	return &out, nil
}
