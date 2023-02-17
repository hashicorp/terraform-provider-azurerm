package service

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PreemptionCapability string

const (
	PreemptionCapabilityMayPreempt PreemptionCapability = "MayPreempt"
	PreemptionCapabilityNotPreempt PreemptionCapability = "NotPreempt"
)

func PossibleValuesForPreemptionCapability() []string {
	return []string{
		string(PreemptionCapabilityMayPreempt),
		string(PreemptionCapabilityNotPreempt),
	}
}

func parsePreemptionCapability(input string) (*PreemptionCapability, error) {
	vals := map[string]PreemptionCapability{
		"maypreempt": PreemptionCapabilityMayPreempt,
		"notpreempt": PreemptionCapabilityNotPreempt,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreemptionCapability(input)
	return &out, nil
}

type PreemptionVulnerability string

const (
	PreemptionVulnerabilityNotPreemptable PreemptionVulnerability = "NotPreemptable"
	PreemptionVulnerabilityPreemptable    PreemptionVulnerability = "Preemptable"
)

func PossibleValuesForPreemptionVulnerability() []string {
	return []string{
		string(PreemptionVulnerabilityNotPreemptable),
		string(PreemptionVulnerabilityPreemptable),
	}
}

func parsePreemptionVulnerability(input string) (*PreemptionVulnerability, error) {
	vals := map[string]PreemptionVulnerability{
		"notpreemptable": PreemptionVulnerabilityNotPreemptable,
		"preemptable":    PreemptionVulnerabilityPreemptable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreemptionVulnerability(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateDeleted   ProvisioningState = "Deleted"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUnknown   ProvisioningState = "Unknown"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
		"canceled":  ProvisioningStateCanceled,
		"deleted":   ProvisioningStateDeleted,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"unknown":   ProvisioningStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type SdfDirection string

const (
	SdfDirectionBidirectional SdfDirection = "Bidirectional"
	SdfDirectionDownlink      SdfDirection = "Downlink"
	SdfDirectionUplink        SdfDirection = "Uplink"
)

func PossibleValuesForSdfDirection() []string {
	return []string{
		string(SdfDirectionBidirectional),
		string(SdfDirectionDownlink),
		string(SdfDirectionUplink),
	}
}

func parseSdfDirection(input string) (*SdfDirection, error) {
	vals := map[string]SdfDirection{
		"bidirectional": SdfDirectionBidirectional,
		"downlink":      SdfDirectionDownlink,
		"uplink":        SdfDirectionUplink,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SdfDirection(input)
	return &out, nil
}

type TrafficControlPermission string

const (
	TrafficControlPermissionBlocked TrafficControlPermission = "Blocked"
	TrafficControlPermissionEnabled TrafficControlPermission = "Enabled"
)

func PossibleValuesForTrafficControlPermission() []string {
	return []string{
		string(TrafficControlPermissionBlocked),
		string(TrafficControlPermissionEnabled),
	}
}

func parseTrafficControlPermission(input string) (*TrafficControlPermission, error) {
	vals := map[string]TrafficControlPermission{
		"blocked": TrafficControlPermissionBlocked,
		"enabled": TrafficControlPermissionEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrafficControlPermission(input)
	return &out, nil
}
