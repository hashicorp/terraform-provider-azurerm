package profiles

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CanMigrateDefaultSku string

const (
	CanMigrateDefaultSkuPremiumAzureFrontDoor  CanMigrateDefaultSku = "Premium_AzureFrontDoor"
	CanMigrateDefaultSkuStandardAzureFrontDoor CanMigrateDefaultSku = "Standard_AzureFrontDoor"
)

func PossibleValuesForCanMigrateDefaultSku() []string {
	return []string{
		string(CanMigrateDefaultSkuPremiumAzureFrontDoor),
		string(CanMigrateDefaultSkuStandardAzureFrontDoor),
	}
}

func (s *CanMigrateDefaultSku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCanMigrateDefaultSku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCanMigrateDefaultSku(input string) (*CanMigrateDefaultSku, error) {
	vals := map[string]CanMigrateDefaultSku{
		"premium_azurefrontdoor":  CanMigrateDefaultSkuPremiumAzureFrontDoor,
		"standard_azurefrontdoor": CanMigrateDefaultSkuStandardAzureFrontDoor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CanMigrateDefaultSku(input)
	return &out, nil
}

type OptimizationType string

const (
	OptimizationTypeDynamicSiteAcceleration     OptimizationType = "DynamicSiteAcceleration"
	OptimizationTypeGeneralMediaStreaming       OptimizationType = "GeneralMediaStreaming"
	OptimizationTypeGeneralWebDelivery          OptimizationType = "GeneralWebDelivery"
	OptimizationTypeLargeFileDownload           OptimizationType = "LargeFileDownload"
	OptimizationTypeVideoOnDemandMediaStreaming OptimizationType = "VideoOnDemandMediaStreaming"
)

func PossibleValuesForOptimizationType() []string {
	return []string{
		string(OptimizationTypeDynamicSiteAcceleration),
		string(OptimizationTypeGeneralMediaStreaming),
		string(OptimizationTypeGeneralWebDelivery),
		string(OptimizationTypeLargeFileDownload),
		string(OptimizationTypeVideoOnDemandMediaStreaming),
	}
}

func (s *OptimizationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOptimizationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOptimizationType(input string) (*OptimizationType, error) {
	vals := map[string]OptimizationType{
		"dynamicsiteacceleration":     OptimizationTypeDynamicSiteAcceleration,
		"generalmediastreaming":       OptimizationTypeGeneralMediaStreaming,
		"generalwebdelivery":          OptimizationTypeGeneralWebDelivery,
		"largefiledownload":           OptimizationTypeLargeFileDownload,
		"videoondemandmediastreaming": OptimizationTypeVideoOnDemandMediaStreaming,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OptimizationType(input)
	return &out, nil
}

type ProfileProvisioningState string

const (
	ProfileProvisioningStateCreating  ProfileProvisioningState = "Creating"
	ProfileProvisioningStateDeleting  ProfileProvisioningState = "Deleting"
	ProfileProvisioningStateFailed    ProfileProvisioningState = "Failed"
	ProfileProvisioningStateSucceeded ProfileProvisioningState = "Succeeded"
	ProfileProvisioningStateUpdating  ProfileProvisioningState = "Updating"
)

func PossibleValuesForProfileProvisioningState() []string {
	return []string{
		string(ProfileProvisioningStateCreating),
		string(ProfileProvisioningStateDeleting),
		string(ProfileProvisioningStateFailed),
		string(ProfileProvisioningStateSucceeded),
		string(ProfileProvisioningStateUpdating),
	}
}

func (s *ProfileProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProfileProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProfileProvisioningState(input string) (*ProfileProvisioningState, error) {
	vals := map[string]ProfileProvisioningState{
		"creating":  ProfileProvisioningStateCreating,
		"deleting":  ProfileProvisioningStateDeleting,
		"failed":    ProfileProvisioningStateFailed,
		"succeeded": ProfileProvisioningStateSucceeded,
		"updating":  ProfileProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProfileProvisioningState(input)
	return &out, nil
}

type ProfileResourceState string

const (
	ProfileResourceStateAbortingMigration      ProfileResourceState = "AbortingMigration"
	ProfileResourceStateActive                 ProfileResourceState = "Active"
	ProfileResourceStateCommittingMigration    ProfileResourceState = "CommittingMigration"
	ProfileResourceStateCreating               ProfileResourceState = "Creating"
	ProfileResourceStateDeleting               ProfileResourceState = "Deleting"
	ProfileResourceStateDisabled               ProfileResourceState = "Disabled"
	ProfileResourceStateMigrated               ProfileResourceState = "Migrated"
	ProfileResourceStateMigrating              ProfileResourceState = "Migrating"
	ProfileResourceStatePendingMigrationCommit ProfileResourceState = "PendingMigrationCommit"
)

func PossibleValuesForProfileResourceState() []string {
	return []string{
		string(ProfileResourceStateAbortingMigration),
		string(ProfileResourceStateActive),
		string(ProfileResourceStateCommittingMigration),
		string(ProfileResourceStateCreating),
		string(ProfileResourceStateDeleting),
		string(ProfileResourceStateDisabled),
		string(ProfileResourceStateMigrated),
		string(ProfileResourceStateMigrating),
		string(ProfileResourceStatePendingMigrationCommit),
	}
}

func (s *ProfileResourceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProfileResourceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProfileResourceState(input string) (*ProfileResourceState, error) {
	vals := map[string]ProfileResourceState{
		"abortingmigration":      ProfileResourceStateAbortingMigration,
		"active":                 ProfileResourceStateActive,
		"committingmigration":    ProfileResourceStateCommittingMigration,
		"creating":               ProfileResourceStateCreating,
		"deleting":               ProfileResourceStateDeleting,
		"disabled":               ProfileResourceStateDisabled,
		"migrated":               ProfileResourceStateMigrated,
		"migrating":              ProfileResourceStateMigrating,
		"pendingmigrationcommit": ProfileResourceStatePendingMigrationCommit,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProfileResourceState(input)
	return &out, nil
}

type ProfileScrubbingState string

const (
	ProfileScrubbingStateDisabled ProfileScrubbingState = "Disabled"
	ProfileScrubbingStateEnabled  ProfileScrubbingState = "Enabled"
)

func PossibleValuesForProfileScrubbingState() []string {
	return []string{
		string(ProfileScrubbingStateDisabled),
		string(ProfileScrubbingStateEnabled),
	}
}

func (s *ProfileScrubbingState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProfileScrubbingState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProfileScrubbingState(input string) (*ProfileScrubbingState, error) {
	vals := map[string]ProfileScrubbingState{
		"disabled": ProfileScrubbingStateDisabled,
		"enabled":  ProfileScrubbingStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProfileScrubbingState(input)
	return &out, nil
}

type ResourceUsageUnit string

const (
	ResourceUsageUnitCount ResourceUsageUnit = "count"
)

func PossibleValuesForResourceUsageUnit() []string {
	return []string{
		string(ResourceUsageUnitCount),
	}
}

func (s *ResourceUsageUnit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceUsageUnit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceUsageUnit(input string) (*ResourceUsageUnit, error) {
	vals := map[string]ResourceUsageUnit{
		"count": ResourceUsageUnitCount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceUsageUnit(input)
	return &out, nil
}

type ScrubbingRuleEntryMatchOperator string

const (
	ScrubbingRuleEntryMatchOperatorEqualsAny ScrubbingRuleEntryMatchOperator = "EqualsAny"
)

func PossibleValuesForScrubbingRuleEntryMatchOperator() []string {
	return []string{
		string(ScrubbingRuleEntryMatchOperatorEqualsAny),
	}
}

func (s *ScrubbingRuleEntryMatchOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScrubbingRuleEntryMatchOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScrubbingRuleEntryMatchOperator(input string) (*ScrubbingRuleEntryMatchOperator, error) {
	vals := map[string]ScrubbingRuleEntryMatchOperator{
		"equalsany": ScrubbingRuleEntryMatchOperatorEqualsAny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScrubbingRuleEntryMatchOperator(input)
	return &out, nil
}

type ScrubbingRuleEntryMatchVariable string

const (
	ScrubbingRuleEntryMatchVariableQueryStringArgNames ScrubbingRuleEntryMatchVariable = "QueryStringArgNames"
	ScrubbingRuleEntryMatchVariableRequestIPAddress    ScrubbingRuleEntryMatchVariable = "RequestIPAddress"
	ScrubbingRuleEntryMatchVariableRequestUri          ScrubbingRuleEntryMatchVariable = "RequestUri"
)

func PossibleValuesForScrubbingRuleEntryMatchVariable() []string {
	return []string{
		string(ScrubbingRuleEntryMatchVariableQueryStringArgNames),
		string(ScrubbingRuleEntryMatchVariableRequestIPAddress),
		string(ScrubbingRuleEntryMatchVariableRequestUri),
	}
}

func (s *ScrubbingRuleEntryMatchVariable) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScrubbingRuleEntryMatchVariable(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScrubbingRuleEntryMatchVariable(input string) (*ScrubbingRuleEntryMatchVariable, error) {
	vals := map[string]ScrubbingRuleEntryMatchVariable{
		"querystringargnames": ScrubbingRuleEntryMatchVariableQueryStringArgNames,
		"requestipaddress":    ScrubbingRuleEntryMatchVariableRequestIPAddress,
		"requesturi":          ScrubbingRuleEntryMatchVariableRequestUri,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScrubbingRuleEntryMatchVariable(input)
	return &out, nil
}

type ScrubbingRuleEntryState string

const (
	ScrubbingRuleEntryStateDisabled ScrubbingRuleEntryState = "Disabled"
	ScrubbingRuleEntryStateEnabled  ScrubbingRuleEntryState = "Enabled"
)

func PossibleValuesForScrubbingRuleEntryState() []string {
	return []string{
		string(ScrubbingRuleEntryStateDisabled),
		string(ScrubbingRuleEntryStateEnabled),
	}
}

func (s *ScrubbingRuleEntryState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScrubbingRuleEntryState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScrubbingRuleEntryState(input string) (*ScrubbingRuleEntryState, error) {
	vals := map[string]ScrubbingRuleEntryState{
		"disabled": ScrubbingRuleEntryStateDisabled,
		"enabled":  ScrubbingRuleEntryStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScrubbingRuleEntryState(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameCustomVerizon                             SkuName = "Custom_Verizon"
	SkuNamePremiumAzureFrontDoor                     SkuName = "Premium_AzureFrontDoor"
	SkuNamePremiumVerizon                            SkuName = "Premium_Verizon"
	SkuNameStandardAkamai                            SkuName = "Standard_Akamai"
	SkuNameStandardAvgBandWidthChinaCdn              SkuName = "Standard_AvgBandWidth_ChinaCdn"
	SkuNameStandardAzureFrontDoor                    SkuName = "Standard_AzureFrontDoor"
	SkuNameStandardChinaCdn                          SkuName = "Standard_ChinaCdn"
	SkuNameStandardMicrosoft                         SkuName = "Standard_Microsoft"
	SkuNameStandardNineFiveFiveBandWidthChinaCdn     SkuName = "Standard_955BandWidth_ChinaCdn"
	SkuNameStandardPlusAvgBandWidthChinaCdn          SkuName = "StandardPlus_AvgBandWidth_ChinaCdn"
	SkuNameStandardPlusChinaCdn                      SkuName = "StandardPlus_ChinaCdn"
	SkuNameStandardPlusNineFiveFiveBandWidthChinaCdn SkuName = "StandardPlus_955BandWidth_ChinaCdn"
	SkuNameStandardVerizon                           SkuName = "Standard_Verizon"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameCustomVerizon),
		string(SkuNamePremiumAzureFrontDoor),
		string(SkuNamePremiumVerizon),
		string(SkuNameStandardAkamai),
		string(SkuNameStandardAvgBandWidthChinaCdn),
		string(SkuNameStandardAzureFrontDoor),
		string(SkuNameStandardChinaCdn),
		string(SkuNameStandardMicrosoft),
		string(SkuNameStandardNineFiveFiveBandWidthChinaCdn),
		string(SkuNameStandardPlusAvgBandWidthChinaCdn),
		string(SkuNameStandardPlusChinaCdn),
		string(SkuNameStandardPlusNineFiveFiveBandWidthChinaCdn),
		string(SkuNameStandardVerizon),
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
		"custom_verizon":                     SkuNameCustomVerizon,
		"premium_azurefrontdoor":             SkuNamePremiumAzureFrontDoor,
		"premium_verizon":                    SkuNamePremiumVerizon,
		"standard_akamai":                    SkuNameStandardAkamai,
		"standard_avgbandwidth_chinacdn":     SkuNameStandardAvgBandWidthChinaCdn,
		"standard_azurefrontdoor":            SkuNameStandardAzureFrontDoor,
		"standard_chinacdn":                  SkuNameStandardChinaCdn,
		"standard_microsoft":                 SkuNameStandardMicrosoft,
		"standard_955bandwidth_chinacdn":     SkuNameStandardNineFiveFiveBandWidthChinaCdn,
		"standardplus_avgbandwidth_chinacdn": SkuNameStandardPlusAvgBandWidthChinaCdn,
		"standardplus_chinacdn":              SkuNameStandardPlusChinaCdn,
		"standardplus_955bandwidth_chinacdn": SkuNameStandardPlusNineFiveFiveBandWidthChinaCdn,
		"standard_verizon":                   SkuNameStandardVerizon,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}
