package customizationtasks

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

type CustomizationTaskInputType string

const (
	CustomizationTaskInputTypeBoolean CustomizationTaskInputType = "boolean"
	CustomizationTaskInputTypeNumber  CustomizationTaskInputType = "number"
	CustomizationTaskInputTypeString  CustomizationTaskInputType = "string"
)

func PossibleValuesForCustomizationTaskInputType() []string {
	return []string{
		string(CustomizationTaskInputTypeBoolean),
		string(CustomizationTaskInputTypeNumber),
		string(CustomizationTaskInputTypeString),
	}
}

func (s *CustomizationTaskInputType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomizationTaskInputType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomizationTaskInputType(input string) (*CustomizationTaskInputType, error) {
	vals := map[string]CustomizationTaskInputType{
		"boolean": CustomizationTaskInputTypeBoolean,
		"number":  CustomizationTaskInputTypeNumber,
		"string":  CustomizationTaskInputTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomizationTaskInputType(input)
	return &out, nil
}
