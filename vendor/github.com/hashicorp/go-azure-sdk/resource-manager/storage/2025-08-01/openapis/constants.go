package openapis

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Kind string

const (
	KindBlobStorage      Kind = "BlobStorage"
	KindBlockBlobStorage Kind = "BlockBlobStorage"
	KindFileStorage      Kind = "FileStorage"
	KindStorage          Kind = "Storage"
	KindStorageVTwo      Kind = "StorageV2"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindBlobStorage),
		string(KindBlockBlobStorage),
		string(KindFileStorage),
		string(KindStorage),
		string(KindStorageVTwo),
	}
}

func (s *Kind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"blobstorage":      KindBlobStorage,
		"blockblobstorage": KindBlockBlobStorage,
		"filestorage":      KindFileStorage,
		"storage":          KindStorage,
		"storagev2":        KindStorageVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type ReasonCode string

const (
	ReasonCodeNotAvailableForSubscription ReasonCode = "NotAvailableForSubscription"
	ReasonCodeQuotaId                     ReasonCode = "QuotaId"
)

func PossibleValuesForReasonCode() []string {
	return []string{
		string(ReasonCodeNotAvailableForSubscription),
		string(ReasonCodeQuotaId),
	}
}

func (s *ReasonCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReasonCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReasonCode(input string) (*ReasonCode, error) {
	vals := map[string]ReasonCode{
		"notavailableforsubscription": ReasonCodeNotAvailableForSubscription,
		"quotaid":                     ReasonCodeQuotaId,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReasonCode(input)
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

type SkuTier string

const (
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type UsageUnit string

const (
	UsageUnitBytes           UsageUnit = "Bytes"
	UsageUnitBytesPerSecond  UsageUnit = "BytesPerSecond"
	UsageUnitCount           UsageUnit = "Count"
	UsageUnitCountsPerSecond UsageUnit = "CountsPerSecond"
	UsageUnitPercent         UsageUnit = "Percent"
	UsageUnitSeconds         UsageUnit = "Seconds"
)

func PossibleValuesForUsageUnit() []string {
	return []string{
		string(UsageUnitBytes),
		string(UsageUnitBytesPerSecond),
		string(UsageUnitCount),
		string(UsageUnitCountsPerSecond),
		string(UsageUnitPercent),
		string(UsageUnitSeconds),
	}
}

func (s *UsageUnit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUsageUnit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUsageUnit(input string) (*UsageUnit, error) {
	vals := map[string]UsageUnit{
		"bytes":           UsageUnitBytes,
		"bytespersecond":  UsageUnitBytesPerSecond,
		"count":           UsageUnitCount,
		"countspersecond": UsageUnitCountsPerSecond,
		"percent":         UsageUnitPercent,
		"seconds":         UsageUnitSeconds,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsageUnit(input)
	return &out, nil
}
