package paloaltonetworkscloudngfws

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnableStatus string

const (
	EnableStatusDisabled EnableStatus = "Disabled"
	EnableStatusEnabled  EnableStatus = "Enabled"
)

func PossibleValuesForEnableStatus() []string {
	return []string{
		string(EnableStatusDisabled),
		string(EnableStatusEnabled),
	}
}

func (s *EnableStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnableStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnableStatus(input string) (*EnableStatus, error) {
	vals := map[string]EnableStatus{
		"disabled": EnableStatusDisabled,
		"enabled":  EnableStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnableStatus(input)
	return &out, nil
}

type ProductSerialStatusValues string

const (
	ProductSerialStatusValuesAllocated  ProductSerialStatusValues = "Allocated"
	ProductSerialStatusValuesInProgress ProductSerialStatusValues = "InProgress"
)

func PossibleValuesForProductSerialStatusValues() []string {
	return []string{
		string(ProductSerialStatusValuesAllocated),
		string(ProductSerialStatusValuesInProgress),
	}
}

func (s *ProductSerialStatusValues) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProductSerialStatusValues(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProductSerialStatusValues(input string) (*ProductSerialStatusValues, error) {
	vals := map[string]ProductSerialStatusValues{
		"allocated":  ProductSerialStatusValuesAllocated,
		"inprogress": ProductSerialStatusValuesInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProductSerialStatusValues(input)
	return &out, nil
}

type RegistrationStatus string

const (
	RegistrationStatusNotRegistered RegistrationStatus = "Not Registered"
	RegistrationStatusRegistered    RegistrationStatus = "Registered"
)

func PossibleValuesForRegistrationStatus() []string {
	return []string{
		string(RegistrationStatusNotRegistered),
		string(RegistrationStatusRegistered),
	}
}

func (s *RegistrationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRegistrationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRegistrationStatus(input string) (*RegistrationStatus, error) {
	vals := map[string]RegistrationStatus{
		"not registered": RegistrationStatusNotRegistered,
		"registered":     RegistrationStatusRegistered,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegistrationStatus(input)
	return &out, nil
}
