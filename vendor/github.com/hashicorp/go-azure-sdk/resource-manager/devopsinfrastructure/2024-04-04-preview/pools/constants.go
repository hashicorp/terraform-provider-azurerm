package pools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDevOpsPermissionType string

const (
	AzureDevOpsPermissionTypeCreatorOnly      AzureDevOpsPermissionType = "CreatorOnly"
	AzureDevOpsPermissionTypeInherit          AzureDevOpsPermissionType = "Inherit"
	AzureDevOpsPermissionTypeSpecificAccounts AzureDevOpsPermissionType = "SpecificAccounts"
)

func PossibleValuesForAzureDevOpsPermissionType() []string {
	return []string{
		string(AzureDevOpsPermissionTypeCreatorOnly),
		string(AzureDevOpsPermissionTypeInherit),
		string(AzureDevOpsPermissionTypeSpecificAccounts),
	}
}

func (s *AzureDevOpsPermissionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureDevOpsPermissionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureDevOpsPermissionType(input string) (*AzureDevOpsPermissionType, error) {
	vals := map[string]AzureDevOpsPermissionType{
		"creatoronly":      AzureDevOpsPermissionTypeCreatorOnly,
		"inherit":          AzureDevOpsPermissionTypeInherit,
		"specificaccounts": AzureDevOpsPermissionTypeSpecificAccounts,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureDevOpsPermissionType(input)
	return &out, nil
}

type CachingType string

const (
	CachingTypeNone      CachingType = "None"
	CachingTypeReadOnly  CachingType = "ReadOnly"
	CachingTypeReadWrite CachingType = "ReadWrite"
)

func PossibleValuesForCachingType() []string {
	return []string{
		string(CachingTypeNone),
		string(CachingTypeReadOnly),
		string(CachingTypeReadWrite),
	}
}

func (s *CachingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCachingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCachingType(input string) (*CachingType, error) {
	vals := map[string]CachingType{
		"none":      CachingTypeNone,
		"readonly":  CachingTypeReadOnly,
		"readwrite": CachingTypeReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CachingType(input)
	return &out, nil
}

type LogonType string

const (
	LogonTypeInteractive LogonType = "Interactive"
	LogonTypeService     LogonType = "Service"
)

func PossibleValuesForLogonType() []string {
	return []string{
		string(LogonTypeInteractive),
		string(LogonTypeService),
	}
}

func (s *LogonType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLogonType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLogonType(input string) (*LogonType, error) {
	vals := map[string]LogonType{
		"interactive": LogonTypeInteractive,
		"service":     LogonTypeService,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LogonType(input)
	return &out, nil
}

type OsDiskStorageAccountType string

const (
	OsDiskStorageAccountTypePremium     OsDiskStorageAccountType = "Premium"
	OsDiskStorageAccountTypeStandard    OsDiskStorageAccountType = "Standard"
	OsDiskStorageAccountTypeStandardSSD OsDiskStorageAccountType = "StandardSSD"
)

func PossibleValuesForOsDiskStorageAccountType() []string {
	return []string{
		string(OsDiskStorageAccountTypePremium),
		string(OsDiskStorageAccountTypeStandard),
		string(OsDiskStorageAccountTypeStandardSSD),
	}
}

func (s *OsDiskStorageAccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOsDiskStorageAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOsDiskStorageAccountType(input string) (*OsDiskStorageAccountType, error) {
	vals := map[string]OsDiskStorageAccountType{
		"premium":     OsDiskStorageAccountTypePremium,
		"standard":    OsDiskStorageAccountTypeStandard,
		"standardssd": OsDiskStorageAccountTypeStandardSSD,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OsDiskStorageAccountType(input)
	return &out, nil
}

type PredictionPreference string

const (
	PredictionPreferenceBalanced          PredictionPreference = "Balanced"
	PredictionPreferenceBestPerformance   PredictionPreference = "BestPerformance"
	PredictionPreferenceMoreCostEffective PredictionPreference = "MoreCostEffective"
	PredictionPreferenceMorePerformance   PredictionPreference = "MorePerformance"
	PredictionPreferenceMostCostEffective PredictionPreference = "MostCostEffective"
)

func PossibleValuesForPredictionPreference() []string {
	return []string{
		string(PredictionPreferenceBalanced),
		string(PredictionPreferenceBestPerformance),
		string(PredictionPreferenceMoreCostEffective),
		string(PredictionPreferenceMorePerformance),
		string(PredictionPreferenceMostCostEffective),
	}
}

func (s *PredictionPreference) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePredictionPreference(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePredictionPreference(input string) (*PredictionPreference, error) {
	vals := map[string]PredictionPreference{
		"balanced":          PredictionPreferenceBalanced,
		"bestperformance":   PredictionPreferenceBestPerformance,
		"morecosteffective": PredictionPreferenceMoreCostEffective,
		"moreperformance":   PredictionPreferenceMorePerformance,
		"mostcosteffective": PredictionPreferenceMostCostEffective,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PredictionPreference(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ResourcePredictionsProfileType string

const (
	ResourcePredictionsProfileTypeAutomatic ResourcePredictionsProfileType = "Automatic"
	ResourcePredictionsProfileTypeManual    ResourcePredictionsProfileType = "Manual"
)

func PossibleValuesForResourcePredictionsProfileType() []string {
	return []string{
		string(ResourcePredictionsProfileTypeAutomatic),
		string(ResourcePredictionsProfileTypeManual),
	}
}

func (s *ResourcePredictionsProfileType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourcePredictionsProfileType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourcePredictionsProfileType(input string) (*ResourcePredictionsProfileType, error) {
	vals := map[string]ResourcePredictionsProfileType{
		"automatic": ResourcePredictionsProfileTypeAutomatic,
		"manual":    ResourcePredictionsProfileTypeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourcePredictionsProfileType(input)
	return &out, nil
}

type StorageAccountType string

const (
	StorageAccountTypePremiumLRS     StorageAccountType = "Premium_LRS"
	StorageAccountTypePremiumZRS     StorageAccountType = "Premium_ZRS"
	StorageAccountTypeStandardLRS    StorageAccountType = "Standard_LRS"
	StorageAccountTypeStandardSSDLRS StorageAccountType = "StandardSSD_LRS"
	StorageAccountTypeStandardSSDZRS StorageAccountType = "StandardSSD_ZRS"
)

func PossibleValuesForStorageAccountType() []string {
	return []string{
		string(StorageAccountTypePremiumLRS),
		string(StorageAccountTypePremiumZRS),
		string(StorageAccountTypeStandardLRS),
		string(StorageAccountTypeStandardSSDLRS),
		string(StorageAccountTypeStandardSSDZRS),
	}
}

func (s *StorageAccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageAccountType(input string) (*StorageAccountType, error) {
	vals := map[string]StorageAccountType{
		"premium_lrs":     StorageAccountTypePremiumLRS,
		"premium_zrs":     StorageAccountTypePremiumZRS,
		"standard_lrs":    StorageAccountTypeStandardLRS,
		"standardssd_lrs": StorageAccountTypeStandardSSDLRS,
		"standardssd_zrs": StorageAccountTypeStandardSSDZRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAccountType(input)
	return &out, nil
}
