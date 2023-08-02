package snapshots

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSSKU string

const (
	OSSKUCBLMariner            OSSKU = "CBLMariner"
	OSSKUMariner               OSSKU = "Mariner"
	OSSKUUbuntu                OSSKU = "Ubuntu"
	OSSKUWindowsTwoZeroOneNine OSSKU = "Windows2019"
	OSSKUWindowsTwoZeroTwoTwo  OSSKU = "Windows2022"
)

func PossibleValuesForOSSKU() []string {
	return []string{
		string(OSSKUCBLMariner),
		string(OSSKUMariner),
		string(OSSKUUbuntu),
		string(OSSKUWindowsTwoZeroOneNine),
		string(OSSKUWindowsTwoZeroTwoTwo),
	}
}

func parseOSSKU(input string) (*OSSKU, error) {
	vals := map[string]OSSKU{
		"cblmariner":  OSSKUCBLMariner,
		"mariner":     OSSKUMariner,
		"ubuntu":      OSSKUUbuntu,
		"windows2019": OSSKUWindowsTwoZeroOneNine,
		"windows2022": OSSKUWindowsTwoZeroTwoTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OSSKU(input)
	return &out, nil
}

type OSType string

const (
	OSTypeLinux   OSType = "Linux"
	OSTypeWindows OSType = "Windows"
)

func PossibleValuesForOSType() []string {
	return []string{
		string(OSTypeLinux),
		string(OSTypeWindows),
	}
}

func parseOSType(input string) (*OSType, error) {
	vals := map[string]OSType{
		"linux":   OSTypeLinux,
		"windows": OSTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OSType(input)
	return &out, nil
}

type SnapshotType string

const (
	SnapshotTypeManagedCluster SnapshotType = "ManagedCluster"
	SnapshotTypeNodePool       SnapshotType = "NodePool"
)

func PossibleValuesForSnapshotType() []string {
	return []string{
		string(SnapshotTypeManagedCluster),
		string(SnapshotTypeNodePool),
	}
}

func parseSnapshotType(input string) (*SnapshotType, error) {
	vals := map[string]SnapshotType{
		"managedcluster": SnapshotTypeManagedCluster,
		"nodepool":       SnapshotTypeNodePool,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SnapshotType(input)
	return &out, nil
}
