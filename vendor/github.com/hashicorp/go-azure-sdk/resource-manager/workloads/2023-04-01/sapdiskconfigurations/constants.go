package sapdiskconfigurations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskSkuName string

const (
	DiskSkuNamePremiumLRS     DiskSkuName = "Premium_LRS"
	DiskSkuNamePremiumVTwoLRS DiskSkuName = "PremiumV2_LRS"
	DiskSkuNamePremiumZRS     DiskSkuName = "Premium_ZRS"
	DiskSkuNameStandardLRS    DiskSkuName = "Standard_LRS"
	DiskSkuNameStandardSSDLRS DiskSkuName = "StandardSSD_LRS"
	DiskSkuNameStandardSSDZRS DiskSkuName = "StandardSSD_ZRS"
	DiskSkuNameUltraSSDLRS    DiskSkuName = "UltraSSD_LRS"
)

func PossibleValuesForDiskSkuName() []string {
	return []string{
		string(DiskSkuNamePremiumLRS),
		string(DiskSkuNamePremiumVTwoLRS),
		string(DiskSkuNamePremiumZRS),
		string(DiskSkuNameStandardLRS),
		string(DiskSkuNameStandardSSDLRS),
		string(DiskSkuNameStandardSSDZRS),
		string(DiskSkuNameUltraSSDLRS),
	}
}

func (s *DiskSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskSkuName(input string) (*DiskSkuName, error) {
	vals := map[string]DiskSkuName{
		"premium_lrs":     DiskSkuNamePremiumLRS,
		"premiumv2_lrs":   DiskSkuNamePremiumVTwoLRS,
		"premium_zrs":     DiskSkuNamePremiumZRS,
		"standard_lrs":    DiskSkuNameStandardLRS,
		"standardssd_lrs": DiskSkuNameStandardSSDLRS,
		"standardssd_zrs": DiskSkuNameStandardSSDZRS,
		"ultrassd_lrs":    DiskSkuNameUltraSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskSkuName(input)
	return &out, nil
}

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

func (s *SAPDatabaseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPDatabaseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *SAPDeploymentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPDeploymentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *SAPEnvironmentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPEnvironmentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *SAPProductType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPProductType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
