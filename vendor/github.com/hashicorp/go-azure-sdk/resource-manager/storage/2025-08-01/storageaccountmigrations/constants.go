package storageaccountmigrations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationStatus string

const (
	MigrationStatusComplete               MigrationStatus = "Complete"
	MigrationStatusFailed                 MigrationStatus = "Failed"
	MigrationStatusInProgress             MigrationStatus = "InProgress"
	MigrationStatusInvalid                MigrationStatus = "Invalid"
	MigrationStatusSubmittedForConversion MigrationStatus = "SubmittedForConversion"
)

func PossibleValuesForMigrationStatus() []string {
	return []string{
		string(MigrationStatusComplete),
		string(MigrationStatusFailed),
		string(MigrationStatusInProgress),
		string(MigrationStatusInvalid),
		string(MigrationStatusSubmittedForConversion),
	}
}

func (s *MigrationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMigrationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMigrationStatus(input string) (*MigrationStatus, error) {
	vals := map[string]MigrationStatus{
		"complete":               MigrationStatusComplete,
		"failed":                 MigrationStatusFailed,
		"inprogress":             MigrationStatusInProgress,
		"invalid":                MigrationStatusInvalid,
		"submittedforconversion": MigrationStatusSubmittedForConversion,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MigrationStatus(input)
	return &out, nil
}

type SkuName string

const (
	SkuNamePremiumLRS       SkuName = "Premium_LRS"
	SkuNamePremiumVTwoLRS   SkuName = "PremiumV2_LRS"
	SkuNamePremiumVTwoZRS   SkuName = "PremiumV2_ZRS"
	SkuNamePremiumZRS       SkuName = "Premium_ZRS"
	SkuNameStandardGRS      SkuName = "Standard_GRS"
	SkuNameStandardGZRS     SkuName = "Standard_GZRS"
	SkuNameStandardLRS      SkuName = "Standard_LRS"
	SkuNameStandardRAGRS    SkuName = "Standard_RAGRS"
	SkuNameStandardRAGZRS   SkuName = "Standard_RAGZRS"
	SkuNameStandardVTwoGRS  SkuName = "StandardV2_GRS"
	SkuNameStandardVTwoGZRS SkuName = "StandardV2_GZRS"
	SkuNameStandardVTwoLRS  SkuName = "StandardV2_LRS"
	SkuNameStandardVTwoZRS  SkuName = "StandardV2_ZRS"
	SkuNameStandardZRS      SkuName = "Standard_ZRS"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNamePremiumLRS),
		string(SkuNamePremiumVTwoLRS),
		string(SkuNamePremiumVTwoZRS),
		string(SkuNamePremiumZRS),
		string(SkuNameStandardGRS),
		string(SkuNameStandardGZRS),
		string(SkuNameStandardLRS),
		string(SkuNameStandardRAGRS),
		string(SkuNameStandardRAGZRS),
		string(SkuNameStandardVTwoGRS),
		string(SkuNameStandardVTwoGZRS),
		string(SkuNameStandardVTwoLRS),
		string(SkuNameStandardVTwoZRS),
		string(SkuNameStandardZRS),
	}
}

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"premium_lrs":     SkuNamePremiumLRS,
		"premiumv2_lrs":   SkuNamePremiumVTwoLRS,
		"premiumv2_zrs":   SkuNamePremiumVTwoZRS,
		"premium_zrs":     SkuNamePremiumZRS,
		"standard_grs":    SkuNameStandardGRS,
		"standard_gzrs":   SkuNameStandardGZRS,
		"standard_lrs":    SkuNameStandardLRS,
		"standard_ragrs":  SkuNameStandardRAGRS,
		"standard_ragzrs": SkuNameStandardRAGZRS,
		"standardv2_grs":  SkuNameStandardVTwoGRS,
		"standardv2_gzrs": SkuNameStandardVTwoGZRS,
		"standardv2_lrs":  SkuNameStandardVTwoLRS,
		"standardv2_zrs":  SkuNameStandardVTwoZRS,
		"standard_zrs":    SkuNameStandardZRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}
