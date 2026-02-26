package fleetupdatestrategies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetUpdateStrategyProvisioningState string

const (
	FleetUpdateStrategyProvisioningStateCanceled  FleetUpdateStrategyProvisioningState = "Canceled"
	FleetUpdateStrategyProvisioningStateFailed    FleetUpdateStrategyProvisioningState = "Failed"
	FleetUpdateStrategyProvisioningStateSucceeded FleetUpdateStrategyProvisioningState = "Succeeded"
)

func PossibleValuesForFleetUpdateStrategyProvisioningState() []string {
	return []string{
		string(FleetUpdateStrategyProvisioningStateCanceled),
		string(FleetUpdateStrategyProvisioningStateFailed),
		string(FleetUpdateStrategyProvisioningStateSucceeded),
	}
}

func (s *FleetUpdateStrategyProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFleetUpdateStrategyProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFleetUpdateStrategyProvisioningState(input string) (*FleetUpdateStrategyProvisioningState, error) {
	vals := map[string]FleetUpdateStrategyProvisioningState{
		"canceled":  FleetUpdateStrategyProvisioningStateCanceled,
		"failed":    FleetUpdateStrategyProvisioningStateFailed,
		"succeeded": FleetUpdateStrategyProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FleetUpdateStrategyProvisioningState(input)
	return &out, nil
}
