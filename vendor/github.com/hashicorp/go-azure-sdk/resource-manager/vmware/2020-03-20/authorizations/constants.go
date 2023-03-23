package authorizations

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
