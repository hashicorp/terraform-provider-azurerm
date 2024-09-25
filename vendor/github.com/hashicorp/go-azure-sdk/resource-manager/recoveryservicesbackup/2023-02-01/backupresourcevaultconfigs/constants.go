package backupresourcevaultconfigs

import (
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnhancedSecurityState string

const (
	EnhancedSecurityStateDisabled EnhancedSecurityState = "Disabled"
	EnhancedSecurityStateEnabled  EnhancedSecurityState = "Enabled"
	EnhancedSecurityStateInvalid  EnhancedSecurityState = "Invalid"
)

func PossibleValuesForEnhancedSecurityState() []string {
	return []string{
		string(EnhancedSecurityStateDisabled),
		string(EnhancedSecurityStateEnabled),
		string(EnhancedSecurityStateInvalid),
	}
}

func parseEnhancedSecurityState(input string) (*EnhancedSecurityState, error) {
	vals := map[string]EnhancedSecurityState{
		"disabled": EnhancedSecurityStateDisabled,
		"enabled":  EnhancedSecurityStateEnabled,
		"invalid":  EnhancedSecurityStateInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnhancedSecurityState(input)
	return &out, nil
}

type SoftDeleteFeatureState string

const (
	SoftDeleteFeatureStateAlwaysON SoftDeleteFeatureState = "AlwaysON"
	SoftDeleteFeatureStateDisabled SoftDeleteFeatureState = "Disabled"
	SoftDeleteFeatureStateEnabled  SoftDeleteFeatureState = "Enabled"
	SoftDeleteFeatureStateInvalid  SoftDeleteFeatureState = "Invalid"
)

func PossibleValuesForSoftDeleteFeatureState() []string {
	return []string{
		string(SoftDeleteFeatureStateAlwaysON),
		string(SoftDeleteFeatureStateDisabled),
		string(SoftDeleteFeatureStateEnabled),
		string(SoftDeleteFeatureStateInvalid),
	}
}

func parseSoftDeleteFeatureState(input string) (*SoftDeleteFeatureState, error) {
	vals := map[string]SoftDeleteFeatureState{
		"alwayson": SoftDeleteFeatureStateAlwaysON,
		"disabled": SoftDeleteFeatureStateDisabled,
		"enabled":  SoftDeleteFeatureStateEnabled,
		"invalid":  SoftDeleteFeatureStateInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SoftDeleteFeatureState(input)
	return &out, nil
}

type StorageType string

const (
	StorageTypeGeoRedundant               StorageType = "GeoRedundant"
	StorageTypeInvalid                    StorageType = "Invalid"
	StorageTypeLocallyRedundant           StorageType = "LocallyRedundant"
	StorageTypeReadAccessGeoZoneRedundant StorageType = "ReadAccessGeoZoneRedundant"
	StorageTypeZoneRedundant              StorageType = "ZoneRedundant"
)

func PossibleValuesForStorageType() []string {
	return []string{
		string(StorageTypeGeoRedundant),
		string(StorageTypeInvalid),
		string(StorageTypeLocallyRedundant),
		string(StorageTypeReadAccessGeoZoneRedundant),
		string(StorageTypeZoneRedundant),
	}
}

func parseStorageType(input string) (*StorageType, error) {
	vals := map[string]StorageType{
		"georedundant":               StorageTypeGeoRedundant,
		"invalid":                    StorageTypeInvalid,
		"locallyredundant":           StorageTypeLocallyRedundant,
		"readaccessgeozoneredundant": StorageTypeReadAccessGeoZoneRedundant,
		"zoneredundant":              StorageTypeZoneRedundant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageType(input)
	return &out, nil
}

type StorageTypeState string

const (
	StorageTypeStateInvalid  StorageTypeState = "Invalid"
	StorageTypeStateLocked   StorageTypeState = "Locked"
	StorageTypeStateUnlocked StorageTypeState = "Unlocked"
)

func PossibleValuesForStorageTypeState() []string {
	return []string{
		string(StorageTypeStateInvalid),
		string(StorageTypeStateLocked),
		string(StorageTypeStateUnlocked),
	}
}

func parseStorageTypeState(input string) (*StorageTypeState, error) {
	vals := map[string]StorageTypeState{
		"invalid":  StorageTypeStateInvalid,
		"locked":   StorageTypeStateLocked,
		"unlocked": StorageTypeStateUnlocked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageTypeState(input)
	return &out, nil
}
