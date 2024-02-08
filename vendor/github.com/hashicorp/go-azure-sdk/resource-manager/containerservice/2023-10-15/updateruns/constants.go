package updateruns

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterUpgradeType string

const (
	ManagedClusterUpgradeTypeFull          ManagedClusterUpgradeType = "Full"
	ManagedClusterUpgradeTypeNodeImageOnly ManagedClusterUpgradeType = "NodeImageOnly"
)

func PossibleValuesForManagedClusterUpgradeType() []string {
	return []string{
		string(ManagedClusterUpgradeTypeFull),
		string(ManagedClusterUpgradeTypeNodeImageOnly),
	}
}

func parseManagedClusterUpgradeType(input string) (*ManagedClusterUpgradeType, error) {
	vals := map[string]ManagedClusterUpgradeType{
		"full":          ManagedClusterUpgradeTypeFull,
		"nodeimageonly": ManagedClusterUpgradeTypeNodeImageOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedClusterUpgradeType(input)
	return &out, nil
}

type NodeImageSelectionType string

const (
	NodeImageSelectionTypeConsistent NodeImageSelectionType = "Consistent"
	NodeImageSelectionTypeLatest     NodeImageSelectionType = "Latest"
)

func PossibleValuesForNodeImageSelectionType() []string {
	return []string{
		string(NodeImageSelectionTypeConsistent),
		string(NodeImageSelectionTypeLatest),
	}
}

func parseNodeImageSelectionType(input string) (*NodeImageSelectionType, error) {
	vals := map[string]NodeImageSelectionType{
		"consistent": NodeImageSelectionTypeConsistent,
		"latest":     NodeImageSelectionTypeLatest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodeImageSelectionType(input)
	return &out, nil
}

type UpdateRunProvisioningState string

const (
	UpdateRunProvisioningStateCanceled  UpdateRunProvisioningState = "Canceled"
	UpdateRunProvisioningStateFailed    UpdateRunProvisioningState = "Failed"
	UpdateRunProvisioningStateSucceeded UpdateRunProvisioningState = "Succeeded"
)

func PossibleValuesForUpdateRunProvisioningState() []string {
	return []string{
		string(UpdateRunProvisioningStateCanceled),
		string(UpdateRunProvisioningStateFailed),
		string(UpdateRunProvisioningStateSucceeded),
	}
}

func parseUpdateRunProvisioningState(input string) (*UpdateRunProvisioningState, error) {
	vals := map[string]UpdateRunProvisioningState{
		"canceled":  UpdateRunProvisioningStateCanceled,
		"failed":    UpdateRunProvisioningStateFailed,
		"succeeded": UpdateRunProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateRunProvisioningState(input)
	return &out, nil
}

type UpdateState string

const (
	UpdateStateCompleted  UpdateState = "Completed"
	UpdateStateFailed     UpdateState = "Failed"
	UpdateStateNotStarted UpdateState = "NotStarted"
	UpdateStateRunning    UpdateState = "Running"
	UpdateStateSkipped    UpdateState = "Skipped"
	UpdateStateStopped    UpdateState = "Stopped"
	UpdateStateStopping   UpdateState = "Stopping"
)

func PossibleValuesForUpdateState() []string {
	return []string{
		string(UpdateStateCompleted),
		string(UpdateStateFailed),
		string(UpdateStateNotStarted),
		string(UpdateStateRunning),
		string(UpdateStateSkipped),
		string(UpdateStateStopped),
		string(UpdateStateStopping),
	}
}

func parseUpdateState(input string) (*UpdateState, error) {
	vals := map[string]UpdateState{
		"completed":  UpdateStateCompleted,
		"failed":     UpdateStateFailed,
		"notstarted": UpdateStateNotStarted,
		"running":    UpdateStateRunning,
		"skipped":    UpdateStateSkipped,
		"stopped":    UpdateStateStopped,
		"stopping":   UpdateStateStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateState(input)
	return &out, nil
}
