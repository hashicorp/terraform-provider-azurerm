package nodetype

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskType string

const (
	DiskTypePremiumLRS     DiskType = "Premium_LRS"
	DiskTypeStandardLRS    DiskType = "Standard_LRS"
	DiskTypeStandardSSDLRS DiskType = "StandardSSD_LRS"
)

func PossibleValuesForDiskType() []string {
	return []string{
		string(DiskTypePremiumLRS),
		string(DiskTypeStandardLRS),
		string(DiskTypeStandardSSDLRS),
	}
}

func parseDiskType(input string) (*DiskType, error) {
	vals := map[string]DiskType{
		"premium_lrs":     DiskTypePremiumLRS,
		"standard_lrs":    DiskTypeStandardLRS,
		"standardssd_lrs": DiskTypeStandardSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskType(input)
	return &out, nil
}

type ManagedResourceProvisioningState string

const (
	ManagedResourceProvisioningStateCanceled  ManagedResourceProvisioningState = "Canceled"
	ManagedResourceProvisioningStateCreated   ManagedResourceProvisioningState = "Created"
	ManagedResourceProvisioningStateCreating  ManagedResourceProvisioningState = "Creating"
	ManagedResourceProvisioningStateDeleted   ManagedResourceProvisioningState = "Deleted"
	ManagedResourceProvisioningStateDeleting  ManagedResourceProvisioningState = "Deleting"
	ManagedResourceProvisioningStateFailed    ManagedResourceProvisioningState = "Failed"
	ManagedResourceProvisioningStateNone      ManagedResourceProvisioningState = "None"
	ManagedResourceProvisioningStateOther     ManagedResourceProvisioningState = "Other"
	ManagedResourceProvisioningStateSucceeded ManagedResourceProvisioningState = "Succeeded"
	ManagedResourceProvisioningStateUpdating  ManagedResourceProvisioningState = "Updating"
)

func PossibleValuesForManagedResourceProvisioningState() []string {
	return []string{
		string(ManagedResourceProvisioningStateCanceled),
		string(ManagedResourceProvisioningStateCreated),
		string(ManagedResourceProvisioningStateCreating),
		string(ManagedResourceProvisioningStateDeleted),
		string(ManagedResourceProvisioningStateDeleting),
		string(ManagedResourceProvisioningStateFailed),
		string(ManagedResourceProvisioningStateNone),
		string(ManagedResourceProvisioningStateOther),
		string(ManagedResourceProvisioningStateSucceeded),
		string(ManagedResourceProvisioningStateUpdating),
	}
}

func parseManagedResourceProvisioningState(input string) (*ManagedResourceProvisioningState, error) {
	vals := map[string]ManagedResourceProvisioningState{
		"canceled":  ManagedResourceProvisioningStateCanceled,
		"created":   ManagedResourceProvisioningStateCreated,
		"creating":  ManagedResourceProvisioningStateCreating,
		"deleted":   ManagedResourceProvisioningStateDeleted,
		"deleting":  ManagedResourceProvisioningStateDeleting,
		"failed":    ManagedResourceProvisioningStateFailed,
		"none":      ManagedResourceProvisioningStateNone,
		"other":     ManagedResourceProvisioningStateOther,
		"succeeded": ManagedResourceProvisioningStateSucceeded,
		"updating":  ManagedResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedResourceProvisioningState(input)
	return &out, nil
}
