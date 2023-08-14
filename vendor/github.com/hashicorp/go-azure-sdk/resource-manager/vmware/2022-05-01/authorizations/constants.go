package authorizations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteAuthorizationProvisioningState string

const (
	ExpressRouteAuthorizationProvisioningStateCanceled  ExpressRouteAuthorizationProvisioningState = "Canceled"
	ExpressRouteAuthorizationProvisioningStateFailed    ExpressRouteAuthorizationProvisioningState = "Failed"
	ExpressRouteAuthorizationProvisioningStateSucceeded ExpressRouteAuthorizationProvisioningState = "Succeeded"
	ExpressRouteAuthorizationProvisioningStateUpdating  ExpressRouteAuthorizationProvisioningState = "Updating"
)

func PossibleValuesForExpressRouteAuthorizationProvisioningState() []string {
	return []string{
		string(ExpressRouteAuthorizationProvisioningStateCanceled),
		string(ExpressRouteAuthorizationProvisioningStateFailed),
		string(ExpressRouteAuthorizationProvisioningStateSucceeded),
		string(ExpressRouteAuthorizationProvisioningStateUpdating),
	}
}

func (s *ExpressRouteAuthorizationProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRouteAuthorizationProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRouteAuthorizationProvisioningState(input string) (*ExpressRouteAuthorizationProvisioningState, error) {
	vals := map[string]ExpressRouteAuthorizationProvisioningState{
		"canceled":  ExpressRouteAuthorizationProvisioningStateCanceled,
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
