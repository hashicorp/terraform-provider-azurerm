package factories

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FactoryIdentityType string

const (
	FactoryIdentityTypeSystemAssigned             FactoryIdentityType = "SystemAssigned"
	FactoryIdentityTypeSystemAssignedUserAssigned FactoryIdentityType = "SystemAssigned,UserAssigned"
	FactoryIdentityTypeUserAssigned               FactoryIdentityType = "UserAssigned"
)

func PossibleValuesForFactoryIdentityType() []string {
	return []string{
		string(FactoryIdentityTypeSystemAssigned),
		string(FactoryIdentityTypeSystemAssignedUserAssigned),
		string(FactoryIdentityTypeUserAssigned),
	}
}

func (s *FactoryIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFactoryIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFactoryIdentityType(input string) (*FactoryIdentityType, error) {
	vals := map[string]FactoryIdentityType{
		"systemassigned":              FactoryIdentityTypeSystemAssigned,
		"systemassigned,userassigned": FactoryIdentityTypeSystemAssignedUserAssigned,
		"userassigned":                FactoryIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FactoryIdentityType(input)
	return &out, nil
}

type GlobalParameterType string

const (
	GlobalParameterTypeArray  GlobalParameterType = "Array"
	GlobalParameterTypeBool   GlobalParameterType = "Bool"
	GlobalParameterTypeFloat  GlobalParameterType = "Float"
	GlobalParameterTypeInt    GlobalParameterType = "Int"
	GlobalParameterTypeObject GlobalParameterType = "Object"
	GlobalParameterTypeString GlobalParameterType = "String"
)

func PossibleValuesForGlobalParameterType() []string {
	return []string{
		string(GlobalParameterTypeArray),
		string(GlobalParameterTypeBool),
		string(GlobalParameterTypeFloat),
		string(GlobalParameterTypeInt),
		string(GlobalParameterTypeObject),
		string(GlobalParameterTypeString),
	}
}

func (s *GlobalParameterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGlobalParameterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGlobalParameterType(input string) (*GlobalParameterType, error) {
	vals := map[string]GlobalParameterType{
		"array":  GlobalParameterTypeArray,
		"bool":   GlobalParameterTypeBool,
		"float":  GlobalParameterTypeFloat,
		"int":    GlobalParameterTypeInt,
		"object": GlobalParameterTypeObject,
		"string": GlobalParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GlobalParameterType(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}
