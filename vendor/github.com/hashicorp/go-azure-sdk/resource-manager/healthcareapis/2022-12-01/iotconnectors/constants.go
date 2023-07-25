package iotconnectors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IotIdentityResolutionType string

const (
	IotIdentityResolutionTypeCreate IotIdentityResolutionType = "Create"
	IotIdentityResolutionTypeLookup IotIdentityResolutionType = "Lookup"
)

func PossibleValuesForIotIdentityResolutionType() []string {
	return []string{
		string(IotIdentityResolutionTypeCreate),
		string(IotIdentityResolutionTypeLookup),
	}
}

func (s *IotIdentityResolutionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIotIdentityResolutionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIotIdentityResolutionType(input string) (*IotIdentityResolutionType, error) {
	vals := map[string]IotIdentityResolutionType{
		"create": IotIdentityResolutionTypeCreate,
		"lookup": IotIdentityResolutionTypeLookup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IotIdentityResolutionType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted          ProvisioningState = "Accepted"
	ProvisioningStateCanceled          ProvisioningState = "Canceled"
	ProvisioningStateCreating          ProvisioningState = "Creating"
	ProvisioningStateDeleting          ProvisioningState = "Deleting"
	ProvisioningStateDeprovisioned     ProvisioningState = "Deprovisioned"
	ProvisioningStateFailed            ProvisioningState = "Failed"
	ProvisioningStateMoving            ProvisioningState = "Moving"
	ProvisioningStateSucceeded         ProvisioningState = "Succeeded"
	ProvisioningStateSuspended         ProvisioningState = "Suspended"
	ProvisioningStateSystemMaintenance ProvisioningState = "SystemMaintenance"
	ProvisioningStateUpdating          ProvisioningState = "Updating"
	ProvisioningStateVerifying         ProvisioningState = "Verifying"
	ProvisioningStateWarned            ProvisioningState = "Warned"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateDeprovisioned),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateSuspended),
		string(ProvisioningStateSystemMaintenance),
		string(ProvisioningStateUpdating),
		string(ProvisioningStateVerifying),
		string(ProvisioningStateWarned),
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
		"accepted":          ProvisioningStateAccepted,
		"canceled":          ProvisioningStateCanceled,
		"creating":          ProvisioningStateCreating,
		"deleting":          ProvisioningStateDeleting,
		"deprovisioned":     ProvisioningStateDeprovisioned,
		"failed":            ProvisioningStateFailed,
		"moving":            ProvisioningStateMoving,
		"succeeded":         ProvisioningStateSucceeded,
		"suspended":         ProvisioningStateSuspended,
		"systemmaintenance": ProvisioningStateSystemMaintenance,
		"updating":          ProvisioningStateUpdating,
		"verifying":         ProvisioningStateVerifying,
		"warned":            ProvisioningStateWarned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
