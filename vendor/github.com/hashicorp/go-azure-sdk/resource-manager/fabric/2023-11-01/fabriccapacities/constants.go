package fabriccapacities

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityReason string

const (
	CheckNameAvailabilityReasonAlreadyExists CheckNameAvailabilityReason = "AlreadyExists"
	CheckNameAvailabilityReasonInvalid       CheckNameAvailabilityReason = "Invalid"
)

func PossibleValuesForCheckNameAvailabilityReason() []string {
	return []string{
		string(CheckNameAvailabilityReasonAlreadyExists),
		string(CheckNameAvailabilityReasonInvalid),
	}
}

func (s *CheckNameAvailabilityReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCheckNameAvailabilityReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCheckNameAvailabilityReason(input string) (*CheckNameAvailabilityReason, error) {
	vals := map[string]CheckNameAvailabilityReason{
		"alreadyexists": CheckNameAvailabilityReasonAlreadyExists,
		"invalid":       CheckNameAvailabilityReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CheckNameAvailabilityReason(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
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
		"canceled":     ProvisioningStateCanceled,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ResourceState string

const (
	ResourceStateActive       ResourceState = "Active"
	ResourceStateDeleting     ResourceState = "Deleting"
	ResourceStateFailed       ResourceState = "Failed"
	ResourceStatePaused       ResourceState = "Paused"
	ResourceStatePausing      ResourceState = "Pausing"
	ResourceStatePreparing    ResourceState = "Preparing"
	ResourceStateProvisioning ResourceState = "Provisioning"
	ResourceStateResuming     ResourceState = "Resuming"
	ResourceStateScaling      ResourceState = "Scaling"
	ResourceStateSuspended    ResourceState = "Suspended"
	ResourceStateSuspending   ResourceState = "Suspending"
	ResourceStateUpdating     ResourceState = "Updating"
)

func PossibleValuesForResourceState() []string {
	return []string{
		string(ResourceStateActive),
		string(ResourceStateDeleting),
		string(ResourceStateFailed),
		string(ResourceStatePaused),
		string(ResourceStatePausing),
		string(ResourceStatePreparing),
		string(ResourceStateProvisioning),
		string(ResourceStateResuming),
		string(ResourceStateScaling),
		string(ResourceStateSuspended),
		string(ResourceStateSuspending),
		string(ResourceStateUpdating),
	}
}

func (s *ResourceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceState(input string) (*ResourceState, error) {
	vals := map[string]ResourceState{
		"active":       ResourceStateActive,
		"deleting":     ResourceStateDeleting,
		"failed":       ResourceStateFailed,
		"paused":       ResourceStatePaused,
		"pausing":      ResourceStatePausing,
		"preparing":    ResourceStatePreparing,
		"provisioning": ResourceStateProvisioning,
		"resuming":     ResourceStateResuming,
		"scaling":      ResourceStateScaling,
		"suspended":    ResourceStateSuspended,
		"suspending":   ResourceStateSuspending,
		"updating":     ResourceStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceState(input)
	return &out, nil
}

type RpSkuTier string

const (
	RpSkuTierFabric RpSkuTier = "Fabric"
)

func PossibleValuesForRpSkuTier() []string {
	return []string{
		string(RpSkuTierFabric),
	}
}

func (s *RpSkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRpSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRpSkuTier(input string) (*RpSkuTier, error) {
	vals := map[string]RpSkuTier{
		"fabric": RpSkuTierFabric,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RpSkuTier(input)
	return &out, nil
}
