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
		string(ExpressRouteAuthorizationProvisioningStateFailed),
		string(ExpressRouteAuthorizationProvisioningStateSucceeded),
		string(ExpressRouteAuthorizationProvisioningStateUpdating),
	}
}

func parseExpressRouteAuthorizationProvisioningState(input string) (*ExpressRouteAuthorizationProvisioningState, error) {
	vals := map[string]ExpressRouteAuthorizationProvisioningState{
		"failed":    ExpressRouteAuthorizationProvisioningStateFailed,
		"succeeded": ExpressRouteAuthorizationProvisioningStateSucceeded,
		"updating":  ExpressRouteAuthorizationProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRouteAuthorizationProvisioningState(input)
	return &out, nil
}
