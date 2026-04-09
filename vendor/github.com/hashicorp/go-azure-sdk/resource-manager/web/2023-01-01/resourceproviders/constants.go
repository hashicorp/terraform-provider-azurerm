package resourceproviders

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServicePlanRestrictions string

const (
	AppServicePlanRestrictionsBasic    AppServicePlanRestrictions = "Basic"
	AppServicePlanRestrictionsFree     AppServicePlanRestrictions = "Free"
	AppServicePlanRestrictionsNone     AppServicePlanRestrictions = "None"
	AppServicePlanRestrictionsPremium  AppServicePlanRestrictions = "Premium"
	AppServicePlanRestrictionsShared   AppServicePlanRestrictions = "Shared"
	AppServicePlanRestrictionsStandard AppServicePlanRestrictions = "Standard"
)

func PossibleValuesForAppServicePlanRestrictions() []string {
	return []string{
		string(AppServicePlanRestrictionsBasic),
		string(AppServicePlanRestrictionsFree),
		string(AppServicePlanRestrictionsNone),
		string(AppServicePlanRestrictionsPremium),
		string(AppServicePlanRestrictionsShared),
		string(AppServicePlanRestrictionsStandard),
	}
}

func (s *AppServicePlanRestrictions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAppServicePlanRestrictions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAppServicePlanRestrictions(input string) (*AppServicePlanRestrictions, error) {
	vals := map[string]AppServicePlanRestrictions{
		"basic":    AppServicePlanRestrictionsBasic,
		"free":     AppServicePlanRestrictionsFree,
		"none":     AppServicePlanRestrictionsNone,
		"premium":  AppServicePlanRestrictionsPremium,
		"shared":   AppServicePlanRestrictionsShared,
		"standard": AppServicePlanRestrictionsStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AppServicePlanRestrictions(input)
	return &out, nil
}

type CheckNameResourceTypes string

const (
	CheckNameResourceTypesHostingEnvironment                   CheckNameResourceTypes = "HostingEnvironment"
	CheckNameResourceTypesMicrosoftPointWebHostingEnvironments CheckNameResourceTypes = "Microsoft.Web/hostingEnvironments"
	CheckNameResourceTypesMicrosoftPointWebPublishingUsers     CheckNameResourceTypes = "Microsoft.Web/publishingUsers"
	CheckNameResourceTypesMicrosoftPointWebSites               CheckNameResourceTypes = "Microsoft.Web/sites"
	CheckNameResourceTypesMicrosoftPointWebSitesSlots          CheckNameResourceTypes = "Microsoft.Web/sites/slots"
	CheckNameResourceTypesPublishingUser                       CheckNameResourceTypes = "PublishingUser"
	CheckNameResourceTypesSite                                 CheckNameResourceTypes = "Site"
	CheckNameResourceTypesSlot                                 CheckNameResourceTypes = "Slot"
)

func PossibleValuesForCheckNameResourceTypes() []string {
	return []string{
		string(CheckNameResourceTypesHostingEnvironment),
		string(CheckNameResourceTypesMicrosoftPointWebHostingEnvironments),
		string(CheckNameResourceTypesMicrosoftPointWebPublishingUsers),
		string(CheckNameResourceTypesMicrosoftPointWebSites),
		string(CheckNameResourceTypesMicrosoftPointWebSitesSlots),
		string(CheckNameResourceTypesPublishingUser),
		string(CheckNameResourceTypesSite),
		string(CheckNameResourceTypesSlot),
	}
}

func (s *CheckNameResourceTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCheckNameResourceTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCheckNameResourceTypes(input string) (*CheckNameResourceTypes, error) {
	vals := map[string]CheckNameResourceTypes{
		"hostingenvironment":                CheckNameResourceTypesHostingEnvironment,
		"microsoft.web/hostingenvironments": CheckNameResourceTypesMicrosoftPointWebHostingEnvironments,
		"microsoft.web/publishingusers":     CheckNameResourceTypesMicrosoftPointWebPublishingUsers,
		"microsoft.web/sites":               CheckNameResourceTypesMicrosoftPointWebSites,
		"microsoft.web/sites/slots":         CheckNameResourceTypesMicrosoftPointWebSitesSlots,
		"publishinguser":                    CheckNameResourceTypesPublishingUser,
		"site":                              CheckNameResourceTypesSite,
		"slot":                              CheckNameResourceTypesSlot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CheckNameResourceTypes(input)
	return &out, nil
}

type CustomDnsSuffixProvisioningState string

const (
	CustomDnsSuffixProvisioningStateDegraded   CustomDnsSuffixProvisioningState = "Degraded"
	CustomDnsSuffixProvisioningStateFailed     CustomDnsSuffixProvisioningState = "Failed"
	CustomDnsSuffixProvisioningStateInProgress CustomDnsSuffixProvisioningState = "InProgress"
	CustomDnsSuffixProvisioningStateSucceeded  CustomDnsSuffixProvisioningState = "Succeeded"
)

func PossibleValuesForCustomDnsSuffixProvisioningState() []string {
	return []string{
		string(CustomDnsSuffixProvisioningStateDegraded),
		string(CustomDnsSuffixProvisioningStateFailed),
		string(CustomDnsSuffixProvisioningStateInProgress),
		string(CustomDnsSuffixProvisioningStateSucceeded),
	}
}

func (s *CustomDnsSuffixProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomDnsSuffixProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomDnsSuffixProvisioningState(input string) (*CustomDnsSuffixProvisioningState, error) {
	vals := map[string]CustomDnsSuffixProvisioningState{
		"degraded":   CustomDnsSuffixProvisioningStateDegraded,
		"failed":     CustomDnsSuffixProvisioningStateFailed,
		"inprogress": CustomDnsSuffixProvisioningStateInProgress,
		"succeeded":  CustomDnsSuffixProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomDnsSuffixProvisioningState(input)
	return &out, nil
}

type HostingEnvironmentStatus string

const (
	HostingEnvironmentStatusDeleting  HostingEnvironmentStatus = "Deleting"
	HostingEnvironmentStatusPreparing HostingEnvironmentStatus = "Preparing"
	HostingEnvironmentStatusReady     HostingEnvironmentStatus = "Ready"
	HostingEnvironmentStatusScaling   HostingEnvironmentStatus = "Scaling"
)

func PossibleValuesForHostingEnvironmentStatus() []string {
	return []string{
		string(HostingEnvironmentStatusDeleting),
		string(HostingEnvironmentStatusPreparing),
		string(HostingEnvironmentStatusReady),
		string(HostingEnvironmentStatusScaling),
	}
}

func (s *HostingEnvironmentStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHostingEnvironmentStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHostingEnvironmentStatus(input string) (*HostingEnvironmentStatus, error) {
	vals := map[string]HostingEnvironmentStatus{
		"deleting":  HostingEnvironmentStatusDeleting,
		"preparing": HostingEnvironmentStatusPreparing,
		"ready":     HostingEnvironmentStatusReady,
		"scaling":   HostingEnvironmentStatusScaling,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostingEnvironmentStatus(input)
	return &out, nil
}

type InAvailabilityReasonType string

const (
	InAvailabilityReasonTypeAlreadyExists InAvailabilityReasonType = "AlreadyExists"
	InAvailabilityReasonTypeInvalid       InAvailabilityReasonType = "Invalid"
)

func PossibleValuesForInAvailabilityReasonType() []string {
	return []string{
		string(InAvailabilityReasonTypeAlreadyExists),
		string(InAvailabilityReasonTypeInvalid),
	}
}

func (s *InAvailabilityReasonType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInAvailabilityReasonType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInAvailabilityReasonType(input string) (*InAvailabilityReasonType, error) {
	vals := map[string]InAvailabilityReasonType{
		"alreadyexists": InAvailabilityReasonTypeAlreadyExists,
		"invalid":       InAvailabilityReasonTypeInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InAvailabilityReasonType(input)
	return &out, nil
}

type LoadBalancingMode string

const (
	LoadBalancingModeNone          LoadBalancingMode = "None"
	LoadBalancingModePublishing    LoadBalancingMode = "Publishing"
	LoadBalancingModeWeb           LoadBalancingMode = "Web"
	LoadBalancingModeWebPublishing LoadBalancingMode = "Web, Publishing"
)

func PossibleValuesForLoadBalancingMode() []string {
	return []string{
		string(LoadBalancingModeNone),
		string(LoadBalancingModePublishing),
		string(LoadBalancingModeWeb),
		string(LoadBalancingModeWebPublishing),
	}
}

func (s *LoadBalancingMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoadBalancingMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoadBalancingMode(input string) (*LoadBalancingMode, error) {
	vals := map[string]LoadBalancingMode{
		"none":            LoadBalancingModeNone,
		"publishing":      LoadBalancingModePublishing,
		"web":             LoadBalancingModeWeb,
		"web, publishing": LoadBalancingModeWebPublishing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancingMode(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled   ProvisioningState = "Canceled"
	ProvisioningStateDeleting   ProvisioningState = "Deleting"
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateInProgress ProvisioningState = "InProgress"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
		string(ProvisioningStateSucceeded),
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
		"canceled":   ProvisioningStateCanceled,
		"deleting":   ProvisioningStateDeleting,
		"failed":     ProvisioningStateFailed,
		"inprogress": ProvisioningStateInProgress,
		"succeeded":  ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameBasic            SkuName = "Basic"
	SkuNameDynamic          SkuName = "Dynamic"
	SkuNameElasticIsolated  SkuName = "ElasticIsolated"
	SkuNameElasticPremium   SkuName = "ElasticPremium"
	SkuNameFree             SkuName = "Free"
	SkuNameIsolated         SkuName = "Isolated"
	SkuNameIsolatedVTwo     SkuName = "IsolatedV2"
	SkuNamePremium          SkuName = "Premium"
	SkuNamePremiumContainer SkuName = "PremiumContainer"
	SkuNamePremiumVThree    SkuName = "PremiumV3"
	SkuNamePremiumVTwo      SkuName = "PremiumV2"
	SkuNameShared           SkuName = "Shared"
	SkuNameStandard         SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameBasic),
		string(SkuNameDynamic),
		string(SkuNameElasticIsolated),
		string(SkuNameElasticPremium),
		string(SkuNameFree),
		string(SkuNameIsolated),
		string(SkuNameIsolatedVTwo),
		string(SkuNamePremium),
		string(SkuNamePremiumContainer),
		string(SkuNamePremiumVThree),
		string(SkuNamePremiumVTwo),
		string(SkuNameShared),
		string(SkuNameStandard),
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
		"basic":            SkuNameBasic,
		"dynamic":          SkuNameDynamic,
		"elasticisolated":  SkuNameElasticIsolated,
		"elasticpremium":   SkuNameElasticPremium,
		"free":             SkuNameFree,
		"isolated":         SkuNameIsolated,
		"isolatedv2":       SkuNameIsolatedVTwo,
		"premium":          SkuNamePremium,
		"premiumcontainer": SkuNamePremiumContainer,
		"premiumv3":        SkuNamePremiumVThree,
		"premiumv2":        SkuNamePremiumVTwo,
		"shared":           SkuNameShared,
		"standard":         SkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type UpgradeAvailability string

const (
	UpgradeAvailabilityNone  UpgradeAvailability = "None"
	UpgradeAvailabilityReady UpgradeAvailability = "Ready"
)

func PossibleValuesForUpgradeAvailability() []string {
	return []string{
		string(UpgradeAvailabilityNone),
		string(UpgradeAvailabilityReady),
	}
}

func (s *UpgradeAvailability) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpgradeAvailability(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpgradeAvailability(input string) (*UpgradeAvailability, error) {
	vals := map[string]UpgradeAvailability{
		"none":  UpgradeAvailabilityNone,
		"ready": UpgradeAvailabilityReady,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpgradeAvailability(input)
	return &out, nil
}

type UpgradePreference string

const (
	UpgradePreferenceEarly  UpgradePreference = "Early"
	UpgradePreferenceLate   UpgradePreference = "Late"
	UpgradePreferenceManual UpgradePreference = "Manual"
	UpgradePreferenceNone   UpgradePreference = "None"
)

func PossibleValuesForUpgradePreference() []string {
	return []string{
		string(UpgradePreferenceEarly),
		string(UpgradePreferenceLate),
		string(UpgradePreferenceManual),
		string(UpgradePreferenceNone),
	}
}

func (s *UpgradePreference) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpgradePreference(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpgradePreference(input string) (*UpgradePreference, error) {
	vals := map[string]UpgradePreference{
		"early":  UpgradePreferenceEarly,
		"late":   UpgradePreferenceLate,
		"manual": UpgradePreferenceManual,
		"none":   UpgradePreferenceNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpgradePreference(input)
	return &out, nil
}

type ValidateResourceTypes string

const (
	ValidateResourceTypesMicrosoftPointWebHostingEnvironments ValidateResourceTypes = "Microsoft.Web/hostingEnvironments"
	ValidateResourceTypesServerFarm                           ValidateResourceTypes = "ServerFarm"
	ValidateResourceTypesSite                                 ValidateResourceTypes = "Site"
)

func PossibleValuesForValidateResourceTypes() []string {
	return []string{
		string(ValidateResourceTypesMicrosoftPointWebHostingEnvironments),
		string(ValidateResourceTypesServerFarm),
		string(ValidateResourceTypesSite),
	}
}

func (s *ValidateResourceTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseValidateResourceTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseValidateResourceTypes(input string) (*ValidateResourceTypes, error) {
	vals := map[string]ValidateResourceTypes{
		"microsoft.web/hostingenvironments": ValidateResourceTypesMicrosoftPointWebHostingEnvironments,
		"serverfarm":                        ValidateResourceTypesServerFarm,
		"site":                              ValidateResourceTypesSite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ValidateResourceTypes(input)
	return &out, nil
}
