package policyassignments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssignmentType string

const (
	AssignmentTypeCustom       AssignmentType = "Custom"
	AssignmentTypeNotSpecified AssignmentType = "NotSpecified"
	AssignmentTypeSystem       AssignmentType = "System"
	AssignmentTypeSystemHidden AssignmentType = "SystemHidden"
)

func PossibleValuesForAssignmentType() []string {
	return []string{
		string(AssignmentTypeCustom),
		string(AssignmentTypeNotSpecified),
		string(AssignmentTypeSystem),
		string(AssignmentTypeSystemHidden),
	}
}

func (s *AssignmentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAssignmentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAssignmentType(input string) (*AssignmentType, error) {
	vals := map[string]AssignmentType{
		"custom":       AssignmentTypeCustom,
		"notspecified": AssignmentTypeNotSpecified,
		"system":       AssignmentTypeSystem,
		"systemhidden": AssignmentTypeSystemHidden,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssignmentType(input)
	return &out, nil
}

type EnforcementMode string

const (
	EnforcementModeDefault      EnforcementMode = "Default"
	EnforcementModeDoNotEnforce EnforcementMode = "DoNotEnforce"
	EnforcementModeEnroll       EnforcementMode = "Enroll"
)

func PossibleValuesForEnforcementMode() []string {
	return []string{
		string(EnforcementModeDefault),
		string(EnforcementModeDoNotEnforce),
		string(EnforcementModeEnroll),
	}
}

func (s *EnforcementMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnforcementMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnforcementMode(input string) (*EnforcementMode, error) {
	vals := map[string]EnforcementMode{
		"default":      EnforcementModeDefault,
		"donotenforce": EnforcementModeDoNotEnforce,
		"enroll":       EnforcementModeEnroll,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnforcementMode(input)
	return &out, nil
}

type OverrideKind string

const (
	OverrideKindDefinitionVersion OverrideKind = "definitionVersion"
	OverrideKindPolicyEffect      OverrideKind = "policyEffect"
)

func PossibleValuesForOverrideKind() []string {
	return []string{
		string(OverrideKindDefinitionVersion),
		string(OverrideKindPolicyEffect),
	}
}

func (s *OverrideKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOverrideKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOverrideKind(input string) (*OverrideKind, error) {
	vals := map[string]OverrideKind{
		"definitionversion": OverrideKindDefinitionVersion,
		"policyeffect":      OverrideKindPolicyEffect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OverrideKind(input)
	return &out, nil
}

type SelectorKind string

const (
	SelectorKindPolicyDefinitionReferenceId SelectorKind = "policyDefinitionReferenceId"
	SelectorKindResourceLocation            SelectorKind = "resourceLocation"
	SelectorKindResourceType                SelectorKind = "resourceType"
	SelectorKindResourceWithoutLocation     SelectorKind = "resourceWithoutLocation"
)

func PossibleValuesForSelectorKind() []string {
	return []string{
		string(SelectorKindPolicyDefinitionReferenceId),
		string(SelectorKindResourceLocation),
		string(SelectorKindResourceType),
		string(SelectorKindResourceWithoutLocation),
	}
}

func (s *SelectorKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSelectorKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSelectorKind(input string) (*SelectorKind, error) {
	vals := map[string]SelectorKind{
		"policydefinitionreferenceid": SelectorKindPolicyDefinitionReferenceId,
		"resourcelocation":            SelectorKindResourceLocation,
		"resourcetype":                SelectorKindResourceType,
		"resourcewithoutlocation":     SelectorKindResourceWithoutLocation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SelectorKind(input)
	return &out, nil
}
