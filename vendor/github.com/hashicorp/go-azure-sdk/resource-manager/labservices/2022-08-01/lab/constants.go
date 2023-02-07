package lab

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionType string

const (
	ConnectionTypeNone    ConnectionType = "None"
	ConnectionTypePrivate ConnectionType = "Private"
	ConnectionTypePublic  ConnectionType = "Public"
)

func PossibleValuesForConnectionType() []string {
	return []string{
		string(ConnectionTypeNone),
		string(ConnectionTypePrivate),
		string(ConnectionTypePublic),
	}
}

func parseConnectionType(input string) (*ConnectionType, error) {
	vals := map[string]ConnectionType{
		"none":    ConnectionTypeNone,
		"private": ConnectionTypePrivate,
		"public":  ConnectionTypePublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionType(input)
	return &out, nil
}

type CreateOption string

const (
	CreateOptionImage      CreateOption = "Image"
	CreateOptionTemplateVM CreateOption = "TemplateVM"
)

func PossibleValuesForCreateOption() []string {
	return []string{
		string(CreateOptionImage),
		string(CreateOptionTemplateVM),
	}
}

func parseCreateOption(input string) (*CreateOption, error) {
	vals := map[string]CreateOption{
		"image":      CreateOptionImage,
		"templatevm": CreateOptionTemplateVM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateOption(input)
	return &out, nil
}

type EnableState string

const (
	EnableStateDisabled EnableState = "Disabled"
	EnableStateEnabled  EnableState = "Enabled"
)

func PossibleValuesForEnableState() []string {
	return []string{
		string(EnableStateDisabled),
		string(EnableStateEnabled),
	}
}

func parseEnableState(input string) (*EnableState, error) {
	vals := map[string]EnableState{
		"disabled": EnableStateDisabled,
		"enabled":  EnableStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnableState(input)
	return &out, nil
}

type LabState string

const (
	LabStateDraft      LabState = "Draft"
	LabStatePublished  LabState = "Published"
	LabStatePublishing LabState = "Publishing"
	LabStateScaling    LabState = "Scaling"
	LabStateSyncing    LabState = "Syncing"
)

func PossibleValuesForLabState() []string {
	return []string{
		string(LabStateDraft),
		string(LabStatePublished),
		string(LabStatePublishing),
		string(LabStateScaling),
		string(LabStateSyncing),
	}
}

func parseLabState(input string) (*LabState, error) {
	vals := map[string]LabState{
		"draft":      LabStateDraft,
		"published":  LabStatePublished,
		"publishing": LabStatePublishing,
		"scaling":    LabStateScaling,
		"syncing":    LabStateSyncing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LabState(input)
	return &out, nil
}

type OsType string

const (
	OsTypeLinux   OsType = "Linux"
	OsTypeWindows OsType = "Windows"
)

func PossibleValuesForOsType() []string {
	return []string{
		string(OsTypeLinux),
		string(OsTypeWindows),
	}
}

func parseOsType(input string) (*OsType, error) {
	vals := map[string]OsType{
		"linux":   OsTypeLinux,
		"windows": OsTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OsType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateLocked    ProvisioningState = "Locked"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateLocked),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"locked":    ProvisioningStateLocked,
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

type ShutdownOnIdleMode string

const (
	ShutdownOnIdleModeLowUsage    ShutdownOnIdleMode = "LowUsage"
	ShutdownOnIdleModeNone        ShutdownOnIdleMode = "None"
	ShutdownOnIdleModeUserAbsence ShutdownOnIdleMode = "UserAbsence"
)

func PossibleValuesForShutdownOnIdleMode() []string {
	return []string{
		string(ShutdownOnIdleModeLowUsage),
		string(ShutdownOnIdleModeNone),
		string(ShutdownOnIdleModeUserAbsence),
	}
}

func parseShutdownOnIdleMode(input string) (*ShutdownOnIdleMode, error) {
	vals := map[string]ShutdownOnIdleMode{
		"lowusage":    ShutdownOnIdleModeLowUsage,
		"none":        ShutdownOnIdleModeNone,
		"userabsence": ShutdownOnIdleModeUserAbsence,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ShutdownOnIdleMode(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierFree),
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    SkuTierBasic,
		"free":     SkuTierFree,
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}
