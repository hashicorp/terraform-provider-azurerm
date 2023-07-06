package cognitiveservicesaccounts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AbusePenaltyAction string

const (
	AbusePenaltyActionBlock    AbusePenaltyAction = "Block"
	AbusePenaltyActionThrottle AbusePenaltyAction = "Throttle"
)

func PossibleValuesForAbusePenaltyAction() []string {
	return []string{
		string(AbusePenaltyActionBlock),
		string(AbusePenaltyActionThrottle),
	}
}

func (s *AbusePenaltyAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAbusePenaltyAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAbusePenaltyAction(input string) (*AbusePenaltyAction, error) {
	vals := map[string]AbusePenaltyAction{
		"block":    AbusePenaltyActionBlock,
		"throttle": AbusePenaltyActionThrottle,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AbusePenaltyAction(input)
	return &out, nil
}

type KeyName string

const (
	KeyNameKeyOne KeyName = "Key1"
	KeyNameKeyTwo KeyName = "Key2"
)

func PossibleValuesForKeyName() []string {
	return []string{
		string(KeyNameKeyOne),
		string(KeyNameKeyTwo),
	}
}

func (s *KeyName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyName(input string) (*KeyName, error) {
	vals := map[string]KeyName{
		"key1": KeyNameKeyOne,
		"key2": KeyNameKeyTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyName(input)
	return &out, nil
}

type KeySource string

const (
	KeySourceMicrosoftPointCognitiveServices KeySource = "Microsoft.CognitiveServices"
	KeySourceMicrosoftPointKeyVault          KeySource = "Microsoft.KeyVault"
)

func PossibleValuesForKeySource() []string {
	return []string{
		string(KeySourceMicrosoftPointCognitiveServices),
		string(KeySourceMicrosoftPointKeyVault),
	}
}

func (s *KeySource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeySource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"microsoft.cognitiveservices": KeySourceMicrosoftPointCognitiveServices,
		"microsoft.keyvault":          KeySourceMicrosoftPointKeyVault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeySource(input)
	return &out, nil
}

type ModelLifecycleStatus string

const (
	ModelLifecycleStatusGenerallyAvailable ModelLifecycleStatus = "GenerallyAvailable"
	ModelLifecycleStatusPreview            ModelLifecycleStatus = "Preview"
)

func PossibleValuesForModelLifecycleStatus() []string {
	return []string{
		string(ModelLifecycleStatusGenerallyAvailable),
		string(ModelLifecycleStatusPreview),
	}
}

func (s *ModelLifecycleStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseModelLifecycleStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseModelLifecycleStatus(input string) (*ModelLifecycleStatus, error) {
	vals := map[string]ModelLifecycleStatus{
		"generallyavailable": ModelLifecycleStatusGenerallyAvailable,
		"preview":            ModelLifecycleStatusPreview,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ModelLifecycleStatus(input)
	return &out, nil
}

type NetworkRuleAction string

const (
	NetworkRuleActionAllow NetworkRuleAction = "Allow"
	NetworkRuleActionDeny  NetworkRuleAction = "Deny"
)

func PossibleValuesForNetworkRuleAction() []string {
	return []string{
		string(NetworkRuleActionAllow),
		string(NetworkRuleActionDeny),
	}
}

func (s *NetworkRuleAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkRuleAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkRuleAction(input string) (*NetworkRuleAction, error) {
	vals := map[string]NetworkRuleAction{
		"allow": NetworkRuleActionAllow,
		"deny":  NetworkRuleActionDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkRuleAction(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
	}
}

func (s *PrivateEndpointConnectionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointConnectionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":  PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusPending  PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		string(PrivateEndpointServiceConnectionStatusApproved),
		string(PrivateEndpointServiceConnectionStatusPending),
		string(PrivateEndpointServiceConnectionStatusRejected),
	}
}

func (s *PrivateEndpointServiceConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointServiceConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointServiceConnectionStatus(input string) (*PrivateEndpointServiceConnectionStatus, error) {
	vals := map[string]PrivateEndpointServiceConnectionStatus{
		"approved": PrivateEndpointServiceConnectionStatusApproved,
		"pending":  PrivateEndpointServiceConnectionStatusPending,
		"rejected": PrivateEndpointServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointServiceConnectionStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateMoving       ProvisioningState = "Moving"
	ProvisioningStateResolvingDNS ProvisioningState = "ResolvingDNS"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateResolvingDNS),
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
		"accepted":     ProvisioningStateAccepted,
		"creating":     ProvisioningStateCreating,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"moving":       ProvisioningStateMoving,
		"resolvingdns": ProvisioningStateResolvingDNS,
		"succeeded":    ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type QuotaUsageStatus string

const (
	QuotaUsageStatusBlocked   QuotaUsageStatus = "Blocked"
	QuotaUsageStatusInOverage QuotaUsageStatus = "InOverage"
	QuotaUsageStatusIncluded  QuotaUsageStatus = "Included"
	QuotaUsageStatusUnknown   QuotaUsageStatus = "Unknown"
)

func PossibleValuesForQuotaUsageStatus() []string {
	return []string{
		string(QuotaUsageStatusBlocked),
		string(QuotaUsageStatusInOverage),
		string(QuotaUsageStatusIncluded),
		string(QuotaUsageStatusUnknown),
	}
}

func (s *QuotaUsageStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQuotaUsageStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQuotaUsageStatus(input string) (*QuotaUsageStatus, error) {
	vals := map[string]QuotaUsageStatus{
		"blocked":   QuotaUsageStatusBlocked,
		"inoverage": QuotaUsageStatusInOverage,
		"included":  QuotaUsageStatusIncluded,
		"unknown":   QuotaUsageStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QuotaUsageStatus(input)
	return &out, nil
}

type ResourceSkuRestrictionsReasonCode string

const (
	ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription ResourceSkuRestrictionsReasonCode = "NotAvailableForSubscription"
	ResourceSkuRestrictionsReasonCodeQuotaId                     ResourceSkuRestrictionsReasonCode = "QuotaId"
)

func PossibleValuesForResourceSkuRestrictionsReasonCode() []string {
	return []string{
		string(ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription),
		string(ResourceSkuRestrictionsReasonCodeQuotaId),
	}
}

func (s *ResourceSkuRestrictionsReasonCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceSkuRestrictionsReasonCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceSkuRestrictionsReasonCode(input string) (*ResourceSkuRestrictionsReasonCode, error) {
	vals := map[string]ResourceSkuRestrictionsReasonCode{
		"notavailableforsubscription": ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription,
		"quotaid":                     ResourceSkuRestrictionsReasonCodeQuotaId,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceSkuRestrictionsReasonCode(input)
	return &out, nil
}

type ResourceSkuRestrictionsType string

const (
	ResourceSkuRestrictionsTypeLocation ResourceSkuRestrictionsType = "Location"
	ResourceSkuRestrictionsTypeZone     ResourceSkuRestrictionsType = "Zone"
)

func PossibleValuesForResourceSkuRestrictionsType() []string {
	return []string{
		string(ResourceSkuRestrictionsTypeLocation),
		string(ResourceSkuRestrictionsTypeZone),
	}
}

func (s *ResourceSkuRestrictionsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceSkuRestrictionsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceSkuRestrictionsType(input string) (*ResourceSkuRestrictionsType, error) {
	vals := map[string]ResourceSkuRestrictionsType{
		"location": ResourceSkuRestrictionsTypeLocation,
		"zone":     ResourceSkuRestrictionsTypeZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceSkuRestrictionsType(input)
	return &out, nil
}

type RoutingMethods string

const (
	RoutingMethodsPerformance RoutingMethods = "Performance"
	RoutingMethodsPriority    RoutingMethods = "Priority"
	RoutingMethodsWeighted    RoutingMethods = "Weighted"
)

func PossibleValuesForRoutingMethods() []string {
	return []string{
		string(RoutingMethodsPerformance),
		string(RoutingMethodsPriority),
		string(RoutingMethodsWeighted),
	}
}

func (s *RoutingMethods) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoutingMethods(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoutingMethods(input string) (*RoutingMethods, error) {
	vals := map[string]RoutingMethods{
		"performance": RoutingMethodsPerformance,
		"priority":    RoutingMethodsPriority,
		"weighted":    RoutingMethodsWeighted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingMethods(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic      SkuTier = "Basic"
	SkuTierEnterprise SkuTier = "Enterprise"
	SkuTierFree       SkuTier = "Free"
	SkuTierPremium    SkuTier = "Premium"
	SkuTierStandard   SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierEnterprise),
		string(SkuTierFree),
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
		"basic":      SkuTierBasic,
		"enterprise": SkuTierEnterprise,
		"free":       SkuTierFree,
		"premium":    SkuTierPremium,
		"standard":   SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type UnitType string

const (
	UnitTypeBytes          UnitType = "Bytes"
	UnitTypeBytesPerSecond UnitType = "BytesPerSecond"
	UnitTypeCount          UnitType = "Count"
	UnitTypeCountPerSecond UnitType = "CountPerSecond"
	UnitTypeMilliseconds   UnitType = "Milliseconds"
	UnitTypePercent        UnitType = "Percent"
	UnitTypeSeconds        UnitType = "Seconds"
)

func PossibleValuesForUnitType() []string {
	return []string{
		string(UnitTypeBytes),
		string(UnitTypeBytesPerSecond),
		string(UnitTypeCount),
		string(UnitTypeCountPerSecond),
		string(UnitTypeMilliseconds),
		string(UnitTypePercent),
		string(UnitTypeSeconds),
	}
}

func (s *UnitType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUnitType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUnitType(input string) (*UnitType, error) {
	vals := map[string]UnitType{
		"bytes":          UnitTypeBytes,
		"bytespersecond": UnitTypeBytesPerSecond,
		"count":          UnitTypeCount,
		"countpersecond": UnitTypeCountPerSecond,
		"milliseconds":   UnitTypeMilliseconds,
		"percent":        UnitTypePercent,
		"seconds":        UnitTypeSeconds,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UnitType(input)
	return &out, nil
}
