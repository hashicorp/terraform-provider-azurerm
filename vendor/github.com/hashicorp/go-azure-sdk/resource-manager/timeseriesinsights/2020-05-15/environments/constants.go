package environments

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentKind string

const (
	EnvironmentKindGenOne EnvironmentKind = "Gen1"
	EnvironmentKindGenTwo EnvironmentKind = "Gen2"
)

func PossibleValuesForEnvironmentKind() []string {
	return []string{
		string(EnvironmentKindGenOne),
		string(EnvironmentKindGenTwo),
	}
}

func parseEnvironmentKind(input string) (*EnvironmentKind, error) {
	vals := map[string]EnvironmentKind{
		"gen1": EnvironmentKindGenOne,
		"gen2": EnvironmentKindGenTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnvironmentKind(input)
	return &out, nil
}

type IngressState string

const (
	IngressStateDisabled IngressState = "Disabled"
	IngressStatePaused   IngressState = "Paused"
	IngressStateReady    IngressState = "Ready"
	IngressStateRunning  IngressState = "Running"
	IngressStateUnknown  IngressState = "Unknown"
)

func PossibleValuesForIngressState() []string {
	return []string{
		string(IngressStateDisabled),
		string(IngressStatePaused),
		string(IngressStateReady),
		string(IngressStateRunning),
		string(IngressStateUnknown),
	}
}

func parseIngressState(input string) (*IngressState, error) {
	vals := map[string]IngressState{
		"disabled": IngressStateDisabled,
		"paused":   IngressStatePaused,
		"ready":    IngressStateReady,
		"running":  IngressStateRunning,
		"unknown":  IngressStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IngressState(input)
	return &out, nil
}

type Kind string

const (
	KindGenOne Kind = "Gen1"
	KindGenTwo Kind = "Gen2"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindGenOne),
		string(KindGenTwo),
	}
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"gen1": KindGenOne,
		"gen2": KindGenTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type PropertyType string

const (
	PropertyTypeString PropertyType = "String"
)

func PossibleValuesForPropertyType() []string {
	return []string{
		string(PropertyTypeString),
	}
}

func parsePropertyType(input string) (*PropertyType, error) {
	vals := map[string]PropertyType{
		"string": PropertyTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PropertyType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
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

type SkuName string

const (
	SkuNameLOne SkuName = "L1"
	SkuNamePOne SkuName = "P1"
	SkuNameSOne SkuName = "S1"
	SkuNameSTwo SkuName = "S2"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameLOne),
		string(SkuNamePOne),
		string(SkuNameSOne),
		string(SkuNameSTwo),
	}
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"l1": SkuNameLOne,
		"p1": SkuNamePOne,
		"s1": SkuNameSOne,
		"s2": SkuNameSTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type StorageLimitExceededBehavior string

const (
	StorageLimitExceededBehaviorPauseIngress StorageLimitExceededBehavior = "PauseIngress"
	StorageLimitExceededBehaviorPurgeOldData StorageLimitExceededBehavior = "PurgeOldData"
)

func PossibleValuesForStorageLimitExceededBehavior() []string {
	return []string{
		string(StorageLimitExceededBehaviorPauseIngress),
		string(StorageLimitExceededBehaviorPurgeOldData),
	}
}

func parseStorageLimitExceededBehavior(input string) (*StorageLimitExceededBehavior, error) {
	vals := map[string]StorageLimitExceededBehavior{
		"pauseingress": StorageLimitExceededBehaviorPauseIngress,
		"purgeolddata": StorageLimitExceededBehaviorPurgeOldData,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageLimitExceededBehavior(input)
	return &out, nil
}

type WarmStoragePropertiesState string

const (
	WarmStoragePropertiesStateError   WarmStoragePropertiesState = "Error"
	WarmStoragePropertiesStateOk      WarmStoragePropertiesState = "Ok"
	WarmStoragePropertiesStateUnknown WarmStoragePropertiesState = "Unknown"
)

func PossibleValuesForWarmStoragePropertiesState() []string {
	return []string{
		string(WarmStoragePropertiesStateError),
		string(WarmStoragePropertiesStateOk),
		string(WarmStoragePropertiesStateUnknown),
	}
}

func parseWarmStoragePropertiesState(input string) (*WarmStoragePropertiesState, error) {
	vals := map[string]WarmStoragePropertiesState{
		"error":   WarmStoragePropertiesStateError,
		"ok":      WarmStoragePropertiesStateOk,
		"unknown": WarmStoragePropertiesStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WarmStoragePropertiesState(input)
	return &out, nil
}
