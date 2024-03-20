package backupresourcestorageconfigsnoncrr

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedupState string

const (
	DedupStateDisabled DedupState = "Disabled"
	DedupStateEnabled  DedupState = "Enabled"
	DedupStateInvalid  DedupState = "Invalid"
)

func PossibleValuesForDedupState() []string {
	return []string{
		string(DedupStateDisabled),
		string(DedupStateEnabled),
		string(DedupStateInvalid),
	}
}

func parseDedupState(input string) (*DedupState, error) {
	vals := map[string]DedupState{
		"disabled": DedupStateDisabled,
		"enabled":  DedupStateEnabled,
		"invalid":  DedupStateInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DedupState(input)
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

type XcoolState string

const (
	XcoolStateDisabled XcoolState = "Disabled"
	XcoolStateEnabled  XcoolState = "Enabled"
	XcoolStateInvalid  XcoolState = "Invalid"
)

func PossibleValuesForXcoolState() []string {
	return []string{
		string(XcoolStateDisabled),
		string(XcoolStateEnabled),
		string(XcoolStateInvalid),
	}
}

func parseXcoolState(input string) (*XcoolState, error) {
	vals := map[string]XcoolState{
		"disabled": XcoolStateDisabled,
		"enabled":  XcoolStateEnabled,
		"invalid":  XcoolStateInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := XcoolState(input)
	return &out, nil
}
