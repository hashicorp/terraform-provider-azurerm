package publishedblueprint

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlueprintTargetScope string

const (
	BlueprintTargetScopeManagementGroup BlueprintTargetScope = "managementGroup"
	BlueprintTargetScopeSubscription    BlueprintTargetScope = "subscription"
)

func PossibleValuesForBlueprintTargetScope() []string {
	return []string{
		string(BlueprintTargetScopeManagementGroup),
		string(BlueprintTargetScopeSubscription),
	}
}

func (s *BlueprintTargetScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlueprintTargetScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlueprintTargetScope(input string) (*BlueprintTargetScope, error) {
	vals := map[string]BlueprintTargetScope{
		"managementgroup": BlueprintTargetScopeManagementGroup,
		"subscription":    BlueprintTargetScopeSubscription,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlueprintTargetScope(input)
	return &out, nil
}

type TemplateParameterType string

const (
	TemplateParameterTypeArray        TemplateParameterType = "array"
	TemplateParameterTypeBool         TemplateParameterType = "bool"
	TemplateParameterTypeInt          TemplateParameterType = "int"
	TemplateParameterTypeObject       TemplateParameterType = "object"
	TemplateParameterTypeSecureObject TemplateParameterType = "secureObject"
	TemplateParameterTypeSecureString TemplateParameterType = "secureString"
	TemplateParameterTypeString       TemplateParameterType = "string"
)

func PossibleValuesForTemplateParameterType() []string {
	return []string{
		string(TemplateParameterTypeArray),
		string(TemplateParameterTypeBool),
		string(TemplateParameterTypeInt),
		string(TemplateParameterTypeObject),
		string(TemplateParameterTypeSecureObject),
		string(TemplateParameterTypeSecureString),
		string(TemplateParameterTypeString),
	}
}

func (s *TemplateParameterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTemplateParameterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTemplateParameterType(input string) (*TemplateParameterType, error) {
	vals := map[string]TemplateParameterType{
		"array":        TemplateParameterTypeArray,
		"bool":         TemplateParameterTypeBool,
		"int":          TemplateParameterTypeInt,
		"object":       TemplateParameterTypeObject,
		"secureobject": TemplateParameterTypeSecureObject,
		"securestring": TemplateParameterTypeSecureString,
		"string":       TemplateParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TemplateParameterType(input)
	return &out, nil
}
