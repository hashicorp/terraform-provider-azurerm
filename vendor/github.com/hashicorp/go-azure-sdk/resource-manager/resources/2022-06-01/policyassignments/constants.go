package policyassignments

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnforcementMode string

const (
	EnforcementModeDefault      EnforcementMode = "Default"
	EnforcementModeDoNotEnforce EnforcementMode = "DoNotEnforce"
)

func PossibleValuesForEnforcementMode() []string {
	return []string{
		string(EnforcementModeDefault),
		string(EnforcementModeDoNotEnforce),
	}
}

func parseEnforcementMode(input string) (*EnforcementMode, error) {
	vals := map[string]EnforcementMode{
		"default":      EnforcementModeDefault,
		"donotenforce": EnforcementModeDoNotEnforce,
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
	OverrideKindPolicyEffect OverrideKind = "policyEffect"
)

func PossibleValuesForOverrideKind() []string {
	return []string{
		string(OverrideKindPolicyEffect),
	}
}

func parseOverrideKind(input string) (*OverrideKind, error) {
	vals := map[string]OverrideKind{
		"policyeffect": OverrideKindPolicyEffect,
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
