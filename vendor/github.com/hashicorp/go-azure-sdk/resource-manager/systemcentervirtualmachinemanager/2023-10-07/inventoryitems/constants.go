package inventoryitems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InventoryType string

const (
	InventoryTypeCloud                  InventoryType = "Cloud"
	InventoryTypeVirtualMachine         InventoryType = "VirtualMachine"
	InventoryTypeVirtualMachineTemplate InventoryType = "VirtualMachineTemplate"
	InventoryTypeVirtualNetwork         InventoryType = "VirtualNetwork"
)

func PossibleValuesForInventoryType() []string {
	return []string{
		string(InventoryTypeCloud),
		string(InventoryTypeVirtualMachine),
		string(InventoryTypeVirtualMachineTemplate),
		string(InventoryTypeVirtualNetwork),
	}
}

func (s *InventoryType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInventoryType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInventoryType(input string) (*InventoryType, error) {
	vals := map[string]InventoryType{
		"cloud":                  InventoryTypeCloud,
		"virtualmachine":         InventoryTypeVirtualMachine,
		"virtualmachinetemplate": InventoryTypeVirtualMachineTemplate,
		"virtualnetwork":         InventoryTypeVirtualNetwork,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InventoryType(input)
	return &out, nil
}

type OsType string

const (
	OsTypeLinux   OsType = "Linux"
	OsTypeOther   OsType = "Other"
	OsTypeWindows OsType = "Windows"
)

func PossibleValuesForOsType() []string {
	return []string{
		string(OsTypeLinux),
		string(OsTypeOther),
		string(OsTypeWindows),
	}
}

func (s *OsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOsType(input string) (*OsType, error) {
	vals := map[string]OsType{
		"linux":   OsTypeLinux,
		"other":   OsTypeOther,
		"windows": OsTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OsType(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateAccepted     ResourceProvisioningState = "Accepted"
	ResourceProvisioningStateCanceled     ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateCreated      ResourceProvisioningState = "Created"
	ResourceProvisioningStateDeleting     ResourceProvisioningState = "Deleting"
	ResourceProvisioningStateFailed       ResourceProvisioningState = "Failed"
	ResourceProvisioningStateProvisioning ResourceProvisioningState = "Provisioning"
	ResourceProvisioningStateSucceeded    ResourceProvisioningState = "Succeeded"
	ResourceProvisioningStateUpdating     ResourceProvisioningState = "Updating"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateAccepted),
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateCreated),
		string(ResourceProvisioningStateDeleting),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateProvisioning),
		string(ResourceProvisioningStateSucceeded),
		string(ResourceProvisioningStateUpdating),
	}
}

func (s *ResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceProvisioningState(input string) (*ResourceProvisioningState, error) {
	vals := map[string]ResourceProvisioningState{
		"accepted":     ResourceProvisioningStateAccepted,
		"canceled":     ResourceProvisioningStateCanceled,
		"created":      ResourceProvisioningStateCreated,
		"deleting":     ResourceProvisioningStateDeleting,
		"failed":       ResourceProvisioningStateFailed,
		"provisioning": ResourceProvisioningStateProvisioning,
		"succeeded":    ResourceProvisioningStateSucceeded,
		"updating":     ResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}
