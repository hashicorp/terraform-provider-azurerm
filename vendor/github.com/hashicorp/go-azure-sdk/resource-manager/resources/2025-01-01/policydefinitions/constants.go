package policydefinitions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParameterType string

const (
	ParameterTypeArray    ParameterType = "Array"
	ParameterTypeBoolean  ParameterType = "Boolean"
	ParameterTypeDateTime ParameterType = "DateTime"
	ParameterTypeFloat    ParameterType = "Float"
	ParameterTypeInteger  ParameterType = "Integer"
	ParameterTypeObject   ParameterType = "Object"
	ParameterTypeString   ParameterType = "String"
)

func PossibleValuesForParameterType() []string {
	return []string{
		string(ParameterTypeArray),
		string(ParameterTypeBoolean),
		string(ParameterTypeDateTime),
		string(ParameterTypeFloat),
		string(ParameterTypeInteger),
		string(ParameterTypeObject),
		string(ParameterTypeString),
	}
}

func (s *ParameterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseParameterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseParameterType(input string) (*ParameterType, error) {
	vals := map[string]ParameterType{
		"array":    ParameterTypeArray,
		"boolean":  ParameterTypeBoolean,
		"datetime": ParameterTypeDateTime,
		"float":    ParameterTypeFloat,
		"integer":  ParameterTypeInteger,
		"object":   ParameterTypeObject,
		"string":   ParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ParameterType(input)
	return &out, nil
}

type PolicyType string

const (
	PolicyTypeBuiltIn      PolicyType = "BuiltIn"
	PolicyTypeCustom       PolicyType = "Custom"
	PolicyTypeNotSpecified PolicyType = "NotSpecified"
	PolicyTypeStatic       PolicyType = "Static"
)

func PossibleValuesForPolicyType() []string {
	return []string{
		string(PolicyTypeBuiltIn),
		string(PolicyTypeCustom),
		string(PolicyTypeNotSpecified),
		string(PolicyTypeStatic),
	}
}

func (s *PolicyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyType(input string) (*PolicyType, error) {
	vals := map[string]PolicyType{
		"builtin":      PolicyTypeBuiltIn,
		"custom":       PolicyTypeCustom,
		"notspecified": PolicyTypeNotSpecified,
		"static":       PolicyTypeStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyType(input)
	return &out, nil
}
