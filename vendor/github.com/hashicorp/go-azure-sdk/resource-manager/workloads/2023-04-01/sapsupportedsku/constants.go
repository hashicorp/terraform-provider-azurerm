package sapsupportedsku

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPDatabaseType string

const (
	SAPDatabaseTypeDBTwo SAPDatabaseType = "DB2"
	SAPDatabaseTypeHANA  SAPDatabaseType = "HANA"
)

func PossibleValuesForSAPDatabaseType() []string {
	return []string{
		string(SAPDatabaseTypeDBTwo),
		string(SAPDatabaseTypeHANA),
	}
}

func parseSAPDatabaseType(input string) (*SAPDatabaseType, error) {
	vals := map[string]SAPDatabaseType{
		"db2":  SAPDatabaseTypeDBTwo,
		"hana": SAPDatabaseTypeHANA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPDatabaseType(input)
	return &out, nil
}

type SAPDeploymentType string

const (
	SAPDeploymentTypeSingleServer SAPDeploymentType = "SingleServer"
	SAPDeploymentTypeThreeTier    SAPDeploymentType = "ThreeTier"
)

func PossibleValuesForSAPDeploymentType() []string {
	return []string{
		string(SAPDeploymentTypeSingleServer),
		string(SAPDeploymentTypeThreeTier),
	}
}

func parseSAPDeploymentType(input string) (*SAPDeploymentType, error) {
	vals := map[string]SAPDeploymentType{
		"singleserver": SAPDeploymentTypeSingleServer,
		"threetier":    SAPDeploymentTypeThreeTier,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPDeploymentType(input)
	return &out, nil
}

type SAPEnvironmentType string

const (
	SAPEnvironmentTypeNonProd SAPEnvironmentType = "NonProd"
	SAPEnvironmentTypeProd    SAPEnvironmentType = "Prod"
)

func PossibleValuesForSAPEnvironmentType() []string {
	return []string{
		string(SAPEnvironmentTypeNonProd),
		string(SAPEnvironmentTypeProd),
	}
}

func parseSAPEnvironmentType(input string) (*SAPEnvironmentType, error) {
	vals := map[string]SAPEnvironmentType{
		"nonprod": SAPEnvironmentTypeNonProd,
		"prod":    SAPEnvironmentTypeProd,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPEnvironmentType(input)
	return &out, nil
}

type SAPHighAvailabilityType string

const (
	SAPHighAvailabilityTypeAvailabilitySet  SAPHighAvailabilityType = "AvailabilitySet"
	SAPHighAvailabilityTypeAvailabilityZone SAPHighAvailabilityType = "AvailabilityZone"
)

func PossibleValuesForSAPHighAvailabilityType() []string {
	return []string{
		string(SAPHighAvailabilityTypeAvailabilitySet),
		string(SAPHighAvailabilityTypeAvailabilityZone),
	}
}

func parseSAPHighAvailabilityType(input string) (*SAPHighAvailabilityType, error) {
	vals := map[string]SAPHighAvailabilityType{
		"availabilityset":  SAPHighAvailabilityTypeAvailabilitySet,
		"availabilityzone": SAPHighAvailabilityTypeAvailabilityZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPHighAvailabilityType(input)
	return &out, nil
}

type SAPProductType string

const (
	SAPProductTypeECC       SAPProductType = "ECC"
	SAPProductTypeOther     SAPProductType = "Other"
	SAPProductTypeSFourHANA SAPProductType = "S4HANA"
)

func PossibleValuesForSAPProductType() []string {
	return []string{
		string(SAPProductTypeECC),
		string(SAPProductTypeOther),
		string(SAPProductTypeSFourHANA),
	}
}

func parseSAPProductType(input string) (*SAPProductType, error) {
	vals := map[string]SAPProductType{
		"ecc":    SAPProductTypeECC,
		"other":  SAPProductTypeOther,
		"s4hana": SAPProductTypeSFourHANA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPProductType(input)
	return &out, nil
}
