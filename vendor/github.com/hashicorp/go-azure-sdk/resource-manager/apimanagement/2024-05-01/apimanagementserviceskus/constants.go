package apimanagementserviceskus

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceSkuCapacityScaleType string

const (
	ResourceSkuCapacityScaleTypeAutomatic ResourceSkuCapacityScaleType = "automatic"
	ResourceSkuCapacityScaleTypeManual    ResourceSkuCapacityScaleType = "manual"
	ResourceSkuCapacityScaleTypeNone      ResourceSkuCapacityScaleType = "none"
)

func PossibleValuesForResourceSkuCapacityScaleType() []string {
	return []string{
		string(ResourceSkuCapacityScaleTypeAutomatic),
		string(ResourceSkuCapacityScaleTypeManual),
		string(ResourceSkuCapacityScaleTypeNone),
	}
}

func (s *ResourceSkuCapacityScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceSkuCapacityScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceSkuCapacityScaleType(input string) (*ResourceSkuCapacityScaleType, error) {
	vals := map[string]ResourceSkuCapacityScaleType{
		"automatic": ResourceSkuCapacityScaleTypeAutomatic,
		"manual":    ResourceSkuCapacityScaleTypeManual,
		"none":      ResourceSkuCapacityScaleTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceSkuCapacityScaleType(input)
	return &out, nil
}

type SkuType string

const (
	SkuTypeBasic        SkuType = "Basic"
	SkuTypeBasicVTwo    SkuType = "BasicV2"
	SkuTypeConsumption  SkuType = "Consumption"
	SkuTypeDeveloper    SkuType = "Developer"
	SkuTypeIsolated     SkuType = "Isolated"
	SkuTypePremium      SkuType = "Premium"
	SkuTypeStandard     SkuType = "Standard"
	SkuTypeStandardVTwo SkuType = "StandardV2"
)

func PossibleValuesForSkuType() []string {
	return []string{
		string(SkuTypeBasic),
		string(SkuTypeBasicVTwo),
		string(SkuTypeConsumption),
		string(SkuTypeDeveloper),
		string(SkuTypeIsolated),
		string(SkuTypePremium),
		string(SkuTypeStandard),
		string(SkuTypeStandardVTwo),
	}
}

func (s *SkuType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuType(input string) (*SkuType, error) {
	vals := map[string]SkuType{
		"basic":       SkuTypeBasic,
		"basicv2":     SkuTypeBasicVTwo,
		"consumption": SkuTypeConsumption,
		"developer":   SkuTypeDeveloper,
		"isolated":    SkuTypeIsolated,
		"premium":     SkuTypePremium,
		"standard":    SkuTypeStandard,
		"standardv2":  SkuTypeStandardVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuType(input)
	return &out, nil
}
