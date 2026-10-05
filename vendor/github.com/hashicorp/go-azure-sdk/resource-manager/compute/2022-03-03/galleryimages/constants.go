package galleryimages

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Architecture string

const (
	ArchitectureArmSixFour Architecture = "Arm64"
	ArchitectureXSixFour   Architecture = "x64"
)

func PossibleValuesForArchitecture() []string {
	return []string{
		string(ArchitectureArmSixFour),
		string(ArchitectureXSixFour),
	}
}

func (s *Architecture) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseArchitecture(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseArchitecture(input string) (*Architecture, error) {
	vals := map[string]Architecture{
		"arm64": ArchitectureArmSixFour,
		"x64":   ArchitectureXSixFour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Architecture(input)
	return &out, nil
}

type GalleryProvisioningState string

const (
	GalleryProvisioningStateCreating  GalleryProvisioningState = "Creating"
	GalleryProvisioningStateDeleting  GalleryProvisioningState = "Deleting"
	GalleryProvisioningStateFailed    GalleryProvisioningState = "Failed"
	GalleryProvisioningStateMigrating GalleryProvisioningState = "Migrating"
	GalleryProvisioningStateSucceeded GalleryProvisioningState = "Succeeded"
	GalleryProvisioningStateUpdating  GalleryProvisioningState = "Updating"
)

func PossibleValuesForGalleryProvisioningState() []string {
	return []string{
		string(GalleryProvisioningStateCreating),
		string(GalleryProvisioningStateDeleting),
		string(GalleryProvisioningStateFailed),
		string(GalleryProvisioningStateMigrating),
		string(GalleryProvisioningStateSucceeded),
		string(GalleryProvisioningStateUpdating),
	}
}

func (s *GalleryProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGalleryProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGalleryProvisioningState(input string) (*GalleryProvisioningState, error) {
	vals := map[string]GalleryProvisioningState{
		"creating":  GalleryProvisioningStateCreating,
		"deleting":  GalleryProvisioningStateDeleting,
		"failed":    GalleryProvisioningStateFailed,
		"migrating": GalleryProvisioningStateMigrating,
		"succeeded": GalleryProvisioningStateSucceeded,
		"updating":  GalleryProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GalleryProvisioningState(input)
	return &out, nil
}

type HyperVGeneration string

const (
	HyperVGenerationVOne HyperVGeneration = "V1"
	HyperVGenerationVTwo HyperVGeneration = "V2"
)

func PossibleValuesForHyperVGeneration() []string {
	return []string{
		string(HyperVGenerationVOne),
		string(HyperVGenerationVTwo),
	}
}

func (s *HyperVGeneration) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHyperVGeneration(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHyperVGeneration(input string) (*HyperVGeneration, error) {
	vals := map[string]HyperVGeneration{
		"v1": HyperVGenerationVOne,
		"v2": HyperVGenerationVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HyperVGeneration(input)
	return &out, nil
}

type OperatingSystemStateTypes string

const (
	OperatingSystemStateTypesGeneralized OperatingSystemStateTypes = "Generalized"
	OperatingSystemStateTypesSpecialized OperatingSystemStateTypes = "Specialized"
)

func PossibleValuesForOperatingSystemStateTypes() []string {
	return []string{
		string(OperatingSystemStateTypesGeneralized),
		string(OperatingSystemStateTypesSpecialized),
	}
}

func (s *OperatingSystemStateTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperatingSystemStateTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperatingSystemStateTypes(input string) (*OperatingSystemStateTypes, error) {
	vals := map[string]OperatingSystemStateTypes{
		"generalized": OperatingSystemStateTypesGeneralized,
		"specialized": OperatingSystemStateTypesSpecialized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemStateTypes(input)
	return &out, nil
}

type OperatingSystemTypes string

const (
	OperatingSystemTypesLinux   OperatingSystemTypes = "Linux"
	OperatingSystemTypesWindows OperatingSystemTypes = "Windows"
)

func PossibleValuesForOperatingSystemTypes() []string {
	return []string{
		string(OperatingSystemTypesLinux),
		string(OperatingSystemTypesWindows),
	}
}

func (s *OperatingSystemTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperatingSystemTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperatingSystemTypes(input string) (*OperatingSystemTypes, error) {
	vals := map[string]OperatingSystemTypes{
		"linux":   OperatingSystemTypesLinux,
		"windows": OperatingSystemTypesWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemTypes(input)
	return &out, nil
}
