package monitors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingPreference string

const (
	RoutingPreferenceDefault  RoutingPreference = "Default"
	RoutingPreferenceRouteAll RoutingPreference = "RouteAll"
)

func PossibleValuesForRoutingPreference() []string {
	return []string{
		string(RoutingPreferenceDefault),
		string(RoutingPreferenceRouteAll),
	}
}

func (s *RoutingPreference) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoutingPreference(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoutingPreference(input string) (*RoutingPreference, error) {
	vals := map[string]RoutingPreference{
		"default":  RoutingPreferenceDefault,
		"routeall": RoutingPreferenceRouteAll,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingPreference(input)
	return &out, nil
}

type WorkloadMonitorProvisioningState string

const (
	WorkloadMonitorProvisioningStateAccepted  WorkloadMonitorProvisioningState = "Accepted"
	WorkloadMonitorProvisioningStateCreating  WorkloadMonitorProvisioningState = "Creating"
	WorkloadMonitorProvisioningStateDeleting  WorkloadMonitorProvisioningState = "Deleting"
	WorkloadMonitorProvisioningStateFailed    WorkloadMonitorProvisioningState = "Failed"
	WorkloadMonitorProvisioningStateMigrating WorkloadMonitorProvisioningState = "Migrating"
	WorkloadMonitorProvisioningStateSucceeded WorkloadMonitorProvisioningState = "Succeeded"
	WorkloadMonitorProvisioningStateUpdating  WorkloadMonitorProvisioningState = "Updating"
)

func PossibleValuesForWorkloadMonitorProvisioningState() []string {
	return []string{
		string(WorkloadMonitorProvisioningStateAccepted),
		string(WorkloadMonitorProvisioningStateCreating),
		string(WorkloadMonitorProvisioningStateDeleting),
		string(WorkloadMonitorProvisioningStateFailed),
		string(WorkloadMonitorProvisioningStateMigrating),
		string(WorkloadMonitorProvisioningStateSucceeded),
		string(WorkloadMonitorProvisioningStateUpdating),
	}
}

func (s *WorkloadMonitorProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkloadMonitorProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkloadMonitorProvisioningState(input string) (*WorkloadMonitorProvisioningState, error) {
	vals := map[string]WorkloadMonitorProvisioningState{
		"accepted":  WorkloadMonitorProvisioningStateAccepted,
		"creating":  WorkloadMonitorProvisioningStateCreating,
		"deleting":  WorkloadMonitorProvisioningStateDeleting,
		"failed":    WorkloadMonitorProvisioningStateFailed,
		"migrating": WorkloadMonitorProvisioningStateMigrating,
		"succeeded": WorkloadMonitorProvisioningStateSucceeded,
		"updating":  WorkloadMonitorProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkloadMonitorProvisioningState(input)
	return &out, nil
}
