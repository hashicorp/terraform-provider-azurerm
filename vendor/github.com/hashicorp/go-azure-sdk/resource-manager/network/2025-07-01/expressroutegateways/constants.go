package expressroutegateways

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailoverConnectionStatus string

const (
	FailoverConnectionStatusConnected    FailoverConnectionStatus = "Connected"
	FailoverConnectionStatusDisconnected FailoverConnectionStatus = "Disconnected"
)

func PossibleValuesForFailoverConnectionStatus() []string {
	return []string{
		string(FailoverConnectionStatusConnected),
		string(FailoverConnectionStatusDisconnected),
	}
}

func (s *FailoverConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFailoverConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFailoverConnectionStatus(input string) (*FailoverConnectionStatus, error) {
	vals := map[string]FailoverConnectionStatus{
		"connected":    FailoverConnectionStatusConnected,
		"disconnected": FailoverConnectionStatusDisconnected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FailoverConnectionStatus(input)
	return &out, nil
}

type FailoverTestStatus string

const (
	FailoverTestStatusCompleted   FailoverTestStatus = "Completed"
	FailoverTestStatusExpired     FailoverTestStatus = "Expired"
	FailoverTestStatusInvalid     FailoverTestStatus = "Invalid"
	FailoverTestStatusNotStarted  FailoverTestStatus = "NotStarted"
	FailoverTestStatusRunning     FailoverTestStatus = "Running"
	FailoverTestStatusStartFailed FailoverTestStatus = "StartFailed"
	FailoverTestStatusStarting    FailoverTestStatus = "Starting"
	FailoverTestStatusStopFailed  FailoverTestStatus = "StopFailed"
	FailoverTestStatusStopping    FailoverTestStatus = "Stopping"
)

func PossibleValuesForFailoverTestStatus() []string {
	return []string{
		string(FailoverTestStatusCompleted),
		string(FailoverTestStatusExpired),
		string(FailoverTestStatusInvalid),
		string(FailoverTestStatusNotStarted),
		string(FailoverTestStatusRunning),
		string(FailoverTestStatusStartFailed),
		string(FailoverTestStatusStarting),
		string(FailoverTestStatusStopFailed),
		string(FailoverTestStatusStopping),
	}
}

func (s *FailoverTestStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFailoverTestStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFailoverTestStatus(input string) (*FailoverTestStatus, error) {
	vals := map[string]FailoverTestStatus{
		"completed":   FailoverTestStatusCompleted,
		"expired":     FailoverTestStatusExpired,
		"invalid":     FailoverTestStatusInvalid,
		"notstarted":  FailoverTestStatusNotStarted,
		"running":     FailoverTestStatusRunning,
		"startfailed": FailoverTestStatusStartFailed,
		"starting":    FailoverTestStatusStarting,
		"stopfailed":  FailoverTestStatusStopFailed,
		"stopping":    FailoverTestStatusStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FailoverTestStatus(input)
	return &out, nil
}

type FailoverTestStatusForSingleTest string

const (
	FailoverTestStatusForSingleTestCompleted   FailoverTestStatusForSingleTest = "Completed"
	FailoverTestStatusForSingleTestExpired     FailoverTestStatusForSingleTest = "Expired"
	FailoverTestStatusForSingleTestInvalid     FailoverTestStatusForSingleTest = "Invalid"
	FailoverTestStatusForSingleTestNotStarted  FailoverTestStatusForSingleTest = "NotStarted"
	FailoverTestStatusForSingleTestRunning     FailoverTestStatusForSingleTest = "Running"
	FailoverTestStatusForSingleTestStartFailed FailoverTestStatusForSingleTest = "StartFailed"
	FailoverTestStatusForSingleTestStarting    FailoverTestStatusForSingleTest = "Starting"
	FailoverTestStatusForSingleTestStopFailed  FailoverTestStatusForSingleTest = "StopFailed"
	FailoverTestStatusForSingleTestStopping    FailoverTestStatusForSingleTest = "Stopping"
)

func PossibleValuesForFailoverTestStatusForSingleTest() []string {
	return []string{
		string(FailoverTestStatusForSingleTestCompleted),
		string(FailoverTestStatusForSingleTestExpired),
		string(FailoverTestStatusForSingleTestInvalid),
		string(FailoverTestStatusForSingleTestNotStarted),
		string(FailoverTestStatusForSingleTestRunning),
		string(FailoverTestStatusForSingleTestStartFailed),
		string(FailoverTestStatusForSingleTestStarting),
		string(FailoverTestStatusForSingleTestStopFailed),
		string(FailoverTestStatusForSingleTestStopping),
	}
}

func (s *FailoverTestStatusForSingleTest) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFailoverTestStatusForSingleTest(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFailoverTestStatusForSingleTest(input string) (*FailoverTestStatusForSingleTest, error) {
	vals := map[string]FailoverTestStatusForSingleTest{
		"completed":   FailoverTestStatusForSingleTestCompleted,
		"expired":     FailoverTestStatusForSingleTestExpired,
		"invalid":     FailoverTestStatusForSingleTestInvalid,
		"notstarted":  FailoverTestStatusForSingleTestNotStarted,
		"running":     FailoverTestStatusForSingleTestRunning,
		"startfailed": FailoverTestStatusForSingleTestStartFailed,
		"starting":    FailoverTestStatusForSingleTestStarting,
		"stopfailed":  FailoverTestStatusForSingleTestStopFailed,
		"stopping":    FailoverTestStatusForSingleTestStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FailoverTestStatusForSingleTest(input)
	return &out, nil
}

type FailoverTestType string

const (
	FailoverTestTypeAll                FailoverTestType = "All"
	FailoverTestTypeMultiSiteFailover  FailoverTestType = "MultiSiteFailover"
	FailoverTestTypeSingleSiteFailover FailoverTestType = "SingleSiteFailover"
)

func PossibleValuesForFailoverTestType() []string {
	return []string{
		string(FailoverTestTypeAll),
		string(FailoverTestTypeMultiSiteFailover),
		string(FailoverTestTypeSingleSiteFailover),
	}
}

func (s *FailoverTestType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFailoverTestType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFailoverTestType(input string) (*FailoverTestType, error) {
	vals := map[string]FailoverTestType{
		"all":                FailoverTestTypeAll,
		"multisitefailover":  FailoverTestTypeMultiSiteFailover,
		"singlesitefailover": FailoverTestTypeSingleSiteFailover,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FailoverTestType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type VnetLocalRouteOverrideCriteria string

const (
	VnetLocalRouteOverrideCriteriaContains VnetLocalRouteOverrideCriteria = "Contains"
	VnetLocalRouteOverrideCriteriaEqual    VnetLocalRouteOverrideCriteria = "Equal"
)

func PossibleValuesForVnetLocalRouteOverrideCriteria() []string {
	return []string{
		string(VnetLocalRouteOverrideCriteriaContains),
		string(VnetLocalRouteOverrideCriteriaEqual),
	}
}

func (s *VnetLocalRouteOverrideCriteria) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVnetLocalRouteOverrideCriteria(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVnetLocalRouteOverrideCriteria(input string) (*VnetLocalRouteOverrideCriteria, error) {
	vals := map[string]VnetLocalRouteOverrideCriteria{
		"contains": VnetLocalRouteOverrideCriteriaContains,
		"equal":    VnetLocalRouteOverrideCriteriaEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VnetLocalRouteOverrideCriteria(input)
	return &out, nil
}
