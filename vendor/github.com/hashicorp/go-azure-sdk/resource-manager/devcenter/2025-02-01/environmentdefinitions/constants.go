package environmentdefinitions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CatalogResourceValidationStatus string

const (
	CatalogResourceValidationStatusFailed    CatalogResourceValidationStatus = "Failed"
	CatalogResourceValidationStatusPending   CatalogResourceValidationStatus = "Pending"
	CatalogResourceValidationStatusSucceeded CatalogResourceValidationStatus = "Succeeded"
	CatalogResourceValidationStatusUnknown   CatalogResourceValidationStatus = "Unknown"
)

func PossibleValuesForCatalogResourceValidationStatus() []string {
	return []string{
		string(CatalogResourceValidationStatusFailed),
		string(CatalogResourceValidationStatusPending),
		string(CatalogResourceValidationStatusSucceeded),
		string(CatalogResourceValidationStatusUnknown),
	}
}

func (s *CatalogResourceValidationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCatalogResourceValidationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCatalogResourceValidationStatus(input string) (*CatalogResourceValidationStatus, error) {
	vals := map[string]CatalogResourceValidationStatus{
		"failed":    CatalogResourceValidationStatusFailed,
		"pending":   CatalogResourceValidationStatusPending,
		"succeeded": CatalogResourceValidationStatusSucceeded,
		"unknown":   CatalogResourceValidationStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CatalogResourceValidationStatus(input)
	return &out, nil
}

type ParameterType string

const (
	ParameterTypeArray   ParameterType = "array"
	ParameterTypeBoolean ParameterType = "boolean"
	ParameterTypeInteger ParameterType = "integer"
	ParameterTypeNumber  ParameterType = "number"
	ParameterTypeObject  ParameterType = "object"
	ParameterTypeString  ParameterType = "string"
)

func PossibleValuesForParameterType() []string {
	return []string{
		string(ParameterTypeArray),
		string(ParameterTypeBoolean),
		string(ParameterTypeInteger),
		string(ParameterTypeNumber),
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
		"array":   ParameterTypeArray,
		"boolean": ParameterTypeBoolean,
		"integer": ParameterTypeInteger,
		"number":  ParameterTypeNumber,
		"object":  ParameterTypeObject,
		"string":  ParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ParameterType(input)
	return &out, nil
}
