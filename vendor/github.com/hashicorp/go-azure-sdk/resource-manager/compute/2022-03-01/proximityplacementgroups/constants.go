package proximityplacementgroups

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProximityPlacementGroupType string

const (
	ProximityPlacementGroupTypeStandard ProximityPlacementGroupType = "Standard"
	ProximityPlacementGroupTypeUltra    ProximityPlacementGroupType = "Ultra"
)

func PossibleValuesForProximityPlacementGroupType() []string {
	return []string{
		string(ProximityPlacementGroupTypeStandard),
		string(ProximityPlacementGroupTypeUltra),
	}
}

func parseProximityPlacementGroupType(input string) (*ProximityPlacementGroupType, error) {
	vals := map[string]ProximityPlacementGroupType{
		"standard": ProximityPlacementGroupTypeStandard,
		"ultra":    ProximityPlacementGroupTypeUltra,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProximityPlacementGroupType(input)
	return &out, nil
}

type StatusLevelTypes string

const (
	StatusLevelTypesError   StatusLevelTypes = "Error"
	StatusLevelTypesInfo    StatusLevelTypes = "Info"
	StatusLevelTypesWarning StatusLevelTypes = "Warning"
)

func PossibleValuesForStatusLevelTypes() []string {
	return []string{
		string(StatusLevelTypesError),
		string(StatusLevelTypesInfo),
		string(StatusLevelTypesWarning),
	}
}

func parseStatusLevelTypes(input string) (*StatusLevelTypes, error) {
	vals := map[string]StatusLevelTypes{
		"error":   StatusLevelTypesError,
		"info":    StatusLevelTypesInfo,
		"warning": StatusLevelTypesWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusLevelTypes(input)
	return &out, nil
}
