package virtualmachinetemplates

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllocationMethod string

const (
	AllocationMethodDynamic AllocationMethod = "Dynamic"
	AllocationMethodStatic  AllocationMethod = "Static"
)

func PossibleValuesForAllocationMethod() []string {
	return []string{
		string(AllocationMethodDynamic),
		string(AllocationMethodStatic),
	}
}

func (s *AllocationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAllocationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAllocationMethod(input string) (*AllocationMethod, error) {
	vals := map[string]AllocationMethod{
		"dynamic": AllocationMethodDynamic,
		"static":  AllocationMethodStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AllocationMethod(input)
	return &out, nil
}

type CreateDiffDisk string

const (
	CreateDiffDiskFalse CreateDiffDisk = "false"
	CreateDiffDiskTrue  CreateDiffDisk = "true"
)

func PossibleValuesForCreateDiffDisk() []string {
	return []string{
		string(CreateDiffDiskFalse),
		string(CreateDiffDiskTrue),
	}
}

func (s *CreateDiffDisk) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCreateDiffDisk(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCreateDiffDisk(input string) (*CreateDiffDisk, error) {
	vals := map[string]CreateDiffDisk{
		"false": CreateDiffDiskFalse,
		"true":  CreateDiffDiskTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateDiffDisk(input)
	return &out, nil
}

type DynamicMemoryEnabled string

const (
	DynamicMemoryEnabledFalse DynamicMemoryEnabled = "false"
	DynamicMemoryEnabledTrue  DynamicMemoryEnabled = "true"
)

func PossibleValuesForDynamicMemoryEnabled() []string {
	return []string{
		string(DynamicMemoryEnabledFalse),
		string(DynamicMemoryEnabledTrue),
	}
}

func (s *DynamicMemoryEnabled) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDynamicMemoryEnabled(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDynamicMemoryEnabled(input string) (*DynamicMemoryEnabled, error) {
	vals := map[string]DynamicMemoryEnabled{
		"false": DynamicMemoryEnabledFalse,
		"true":  DynamicMemoryEnabledTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DynamicMemoryEnabled(input)
	return &out, nil
}

type ForceDelete string

const (
	ForceDeleteFalse ForceDelete = "false"
	ForceDeleteTrue  ForceDelete = "true"
)

func PossibleValuesForForceDelete() []string {
	return []string{
		string(ForceDeleteFalse),
		string(ForceDeleteTrue),
	}
}

func (s *ForceDelete) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseForceDelete(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseForceDelete(input string) (*ForceDelete, error) {
	vals := map[string]ForceDelete{
		"false": ForceDeleteFalse,
		"true":  ForceDeleteTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForceDelete(input)
	return &out, nil
}

type IsCustomizable string

const (
	IsCustomizableFalse IsCustomizable = "false"
	IsCustomizableTrue  IsCustomizable = "true"
)

func PossibleValuesForIsCustomizable() []string {
	return []string{
		string(IsCustomizableFalse),
		string(IsCustomizableTrue),
	}
}

func (s *IsCustomizable) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsCustomizable(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIsCustomizable(input string) (*IsCustomizable, error) {
	vals := map[string]IsCustomizable{
		"false": IsCustomizableFalse,
		"true":  IsCustomizableTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsCustomizable(input)
	return &out, nil
}

type IsHighlyAvailable string

const (
	IsHighlyAvailableFalse IsHighlyAvailable = "false"
	IsHighlyAvailableTrue  IsHighlyAvailable = "true"
)

func PossibleValuesForIsHighlyAvailable() []string {
	return []string{
		string(IsHighlyAvailableFalse),
		string(IsHighlyAvailableTrue),
	}
}

func (s *IsHighlyAvailable) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsHighlyAvailable(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIsHighlyAvailable(input string) (*IsHighlyAvailable, error) {
	vals := map[string]IsHighlyAvailable{
		"false": IsHighlyAvailableFalse,
		"true":  IsHighlyAvailableTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsHighlyAvailable(input)
	return &out, nil
}

type LimitCPUForMigration string

const (
	LimitCPUForMigrationFalse LimitCPUForMigration = "false"
	LimitCPUForMigrationTrue  LimitCPUForMigration = "true"
)

func PossibleValuesForLimitCPUForMigration() []string {
	return []string{
		string(LimitCPUForMigrationFalse),
		string(LimitCPUForMigrationTrue),
	}
}

func (s *LimitCPUForMigration) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLimitCPUForMigration(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLimitCPUForMigration(input string) (*LimitCPUForMigration, error) {
	vals := map[string]LimitCPUForMigration{
		"false": LimitCPUForMigrationFalse,
		"true":  LimitCPUForMigrationTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LimitCPUForMigration(input)
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
