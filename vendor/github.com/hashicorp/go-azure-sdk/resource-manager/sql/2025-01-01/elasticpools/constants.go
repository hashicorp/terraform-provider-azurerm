package elasticpools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlwaysEncryptedEnclaveType string

const (
	AlwaysEncryptedEnclaveTypeDefault AlwaysEncryptedEnclaveType = "Default"
	AlwaysEncryptedEnclaveTypeVBS     AlwaysEncryptedEnclaveType = "VBS"
)

func PossibleValuesForAlwaysEncryptedEnclaveType() []string {
	return []string{
		string(AlwaysEncryptedEnclaveTypeDefault),
		string(AlwaysEncryptedEnclaveTypeVBS),
	}
}

func (s *AlwaysEncryptedEnclaveType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlwaysEncryptedEnclaveType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlwaysEncryptedEnclaveType(input string) (*AlwaysEncryptedEnclaveType, error) {
	vals := map[string]AlwaysEncryptedEnclaveType{
		"default": AlwaysEncryptedEnclaveTypeDefault,
		"vbs":     AlwaysEncryptedEnclaveTypeVBS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlwaysEncryptedEnclaveType(input)
	return &out, nil
}

type AvailabilityZoneType string

const (
	AvailabilityZoneTypeNoPreference AvailabilityZoneType = "NoPreference"
	AvailabilityZoneTypeOne          AvailabilityZoneType = "1"
	AvailabilityZoneTypeThree        AvailabilityZoneType = "3"
	AvailabilityZoneTypeTwo          AvailabilityZoneType = "2"
)

func PossibleValuesForAvailabilityZoneType() []string {
	return []string{
		string(AvailabilityZoneTypeNoPreference),
		string(AvailabilityZoneTypeOne),
		string(AvailabilityZoneTypeThree),
		string(AvailabilityZoneTypeTwo),
	}
}

func (s *AvailabilityZoneType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAvailabilityZoneType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAvailabilityZoneType(input string) (*AvailabilityZoneType, error) {
	vals := map[string]AvailabilityZoneType{
		"nopreference": AvailabilityZoneTypeNoPreference,
		"1":            AvailabilityZoneTypeOne,
		"3":            AvailabilityZoneTypeThree,
		"2":            AvailabilityZoneTypeTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AvailabilityZoneType(input)
	return &out, nil
}

type ElasticPoolLicenseType string

const (
	ElasticPoolLicenseTypeBasePrice       ElasticPoolLicenseType = "BasePrice"
	ElasticPoolLicenseTypeLicenseIncluded ElasticPoolLicenseType = "LicenseIncluded"
)

func PossibleValuesForElasticPoolLicenseType() []string {
	return []string{
		string(ElasticPoolLicenseTypeBasePrice),
		string(ElasticPoolLicenseTypeLicenseIncluded),
	}
}

func (s *ElasticPoolLicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseElasticPoolLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseElasticPoolLicenseType(input string) (*ElasticPoolLicenseType, error) {
	vals := map[string]ElasticPoolLicenseType{
		"baseprice":       ElasticPoolLicenseTypeBasePrice,
		"licenseincluded": ElasticPoolLicenseTypeLicenseIncluded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ElasticPoolLicenseType(input)
	return &out, nil
}

type ElasticPoolState string

const (
	ElasticPoolStateCreating ElasticPoolState = "Creating"
	ElasticPoolStateDisabled ElasticPoolState = "Disabled"
	ElasticPoolStateReady    ElasticPoolState = "Ready"
)

func PossibleValuesForElasticPoolState() []string {
	return []string{
		string(ElasticPoolStateCreating),
		string(ElasticPoolStateDisabled),
		string(ElasticPoolStateReady),
	}
}

func (s *ElasticPoolState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseElasticPoolState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseElasticPoolState(input string) (*ElasticPoolState, error) {
	vals := map[string]ElasticPoolState{
		"creating": ElasticPoolStateCreating,
		"disabled": ElasticPoolStateDisabled,
		"ready":    ElasticPoolStateReady,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ElasticPoolState(input)
	return &out, nil
}
