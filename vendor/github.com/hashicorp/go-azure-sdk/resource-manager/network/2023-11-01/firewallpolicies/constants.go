package firewallpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoLearnPrivateRangesMode string

const (
	AutoLearnPrivateRangesModeDisabled AutoLearnPrivateRangesMode = "Disabled"
	AutoLearnPrivateRangesModeEnabled  AutoLearnPrivateRangesMode = "Enabled"
)

func PossibleValuesForAutoLearnPrivateRangesMode() []string {
	return []string{
		string(AutoLearnPrivateRangesModeDisabled),
		string(AutoLearnPrivateRangesModeEnabled),
	}
}

func (s *AutoLearnPrivateRangesMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoLearnPrivateRangesMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoLearnPrivateRangesMode(input string) (*AutoLearnPrivateRangesMode, error) {
	vals := map[string]AutoLearnPrivateRangesMode{
		"disabled": AutoLearnPrivateRangesModeDisabled,
		"enabled":  AutoLearnPrivateRangesModeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoLearnPrivateRangesMode(input)
	return &out, nil
}

type AzureFirewallThreatIntelMode string

const (
	AzureFirewallThreatIntelModeAlert AzureFirewallThreatIntelMode = "Alert"
	AzureFirewallThreatIntelModeDeny  AzureFirewallThreatIntelMode = "Deny"
	AzureFirewallThreatIntelModeOff   AzureFirewallThreatIntelMode = "Off"
)

func PossibleValuesForAzureFirewallThreatIntelMode() []string {
	return []string{
		string(AzureFirewallThreatIntelModeAlert),
		string(AzureFirewallThreatIntelModeDeny),
		string(AzureFirewallThreatIntelModeOff),
	}
}

func (s *AzureFirewallThreatIntelMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFirewallThreatIntelMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFirewallThreatIntelMode(input string) (*AzureFirewallThreatIntelMode, error) {
	vals := map[string]AzureFirewallThreatIntelMode{
		"alert": AzureFirewallThreatIntelModeAlert,
		"deny":  AzureFirewallThreatIntelModeDeny,
		"off":   AzureFirewallThreatIntelModeOff,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFirewallThreatIntelMode(input)
	return &out, nil
}

type FirewallPolicyFilterRuleCollectionActionType string

const (
	FirewallPolicyFilterRuleCollectionActionTypeAllow FirewallPolicyFilterRuleCollectionActionType = "Allow"
	FirewallPolicyFilterRuleCollectionActionTypeDeny  FirewallPolicyFilterRuleCollectionActionType = "Deny"
)

func PossibleValuesForFirewallPolicyFilterRuleCollectionActionType() []string {
	return []string{
		string(FirewallPolicyFilterRuleCollectionActionTypeAllow),
		string(FirewallPolicyFilterRuleCollectionActionTypeDeny),
	}
}

func (s *FirewallPolicyFilterRuleCollectionActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyFilterRuleCollectionActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyFilterRuleCollectionActionType(input string) (*FirewallPolicyFilterRuleCollectionActionType, error) {
	vals := map[string]FirewallPolicyFilterRuleCollectionActionType{
		"allow": FirewallPolicyFilterRuleCollectionActionTypeAllow,
		"deny":  FirewallPolicyFilterRuleCollectionActionTypeDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyFilterRuleCollectionActionType(input)
	return &out, nil
}

type FirewallPolicyIDPSQuerySortOrder string

const (
	FirewallPolicyIDPSQuerySortOrderAscending  FirewallPolicyIDPSQuerySortOrder = "Ascending"
	FirewallPolicyIDPSQuerySortOrderDescending FirewallPolicyIDPSQuerySortOrder = "Descending"
)

func PossibleValuesForFirewallPolicyIDPSQuerySortOrder() []string {
	return []string{
		string(FirewallPolicyIDPSQuerySortOrderAscending),
		string(FirewallPolicyIDPSQuerySortOrderDescending),
	}
}

func (s *FirewallPolicyIDPSQuerySortOrder) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyIDPSQuerySortOrder(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyIDPSQuerySortOrder(input string) (*FirewallPolicyIDPSQuerySortOrder, error) {
	vals := map[string]FirewallPolicyIDPSQuerySortOrder{
		"ascending":  FirewallPolicyIDPSQuerySortOrderAscending,
		"descending": FirewallPolicyIDPSQuerySortOrderDescending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyIDPSQuerySortOrder(input)
	return &out, nil
}

type FirewallPolicyIDPSSignatureDirection int64

const (
	FirewallPolicyIDPSSignatureDirectionFour  FirewallPolicyIDPSSignatureDirection = 4
	FirewallPolicyIDPSSignatureDirectionOne   FirewallPolicyIDPSSignatureDirection = 1
	FirewallPolicyIDPSSignatureDirectionThree FirewallPolicyIDPSSignatureDirection = 3
	FirewallPolicyIDPSSignatureDirectionTwo   FirewallPolicyIDPSSignatureDirection = 2
	FirewallPolicyIDPSSignatureDirectionZero  FirewallPolicyIDPSSignatureDirection = 0
)

func PossibleValuesForFirewallPolicyIDPSSignatureDirection() []int64 {
	return []int64{
		int64(FirewallPolicyIDPSSignatureDirectionFour),
		int64(FirewallPolicyIDPSSignatureDirectionOne),
		int64(FirewallPolicyIDPSSignatureDirectionThree),
		int64(FirewallPolicyIDPSSignatureDirectionTwo),
		int64(FirewallPolicyIDPSSignatureDirectionZero),
	}
}

type FirewallPolicyIDPSSignatureMode int64

const (
	FirewallPolicyIDPSSignatureModeOne  FirewallPolicyIDPSSignatureMode = 1
	FirewallPolicyIDPSSignatureModeTwo  FirewallPolicyIDPSSignatureMode = 2
	FirewallPolicyIDPSSignatureModeZero FirewallPolicyIDPSSignatureMode = 0
)

func PossibleValuesForFirewallPolicyIDPSSignatureMode() []int64 {
	return []int64{
		int64(FirewallPolicyIDPSSignatureModeOne),
		int64(FirewallPolicyIDPSSignatureModeTwo),
		int64(FirewallPolicyIDPSSignatureModeZero),
	}
}

type FirewallPolicyIDPSSignatureSeverity int64

const (
	FirewallPolicyIDPSSignatureSeverityOne   FirewallPolicyIDPSSignatureSeverity = 1
	FirewallPolicyIDPSSignatureSeverityThree FirewallPolicyIDPSSignatureSeverity = 3
	FirewallPolicyIDPSSignatureSeverityTwo   FirewallPolicyIDPSSignatureSeverity = 2
)

func PossibleValuesForFirewallPolicyIDPSSignatureSeverity() []int64 {
	return []int64{
		int64(FirewallPolicyIDPSSignatureSeverityOne),
		int64(FirewallPolicyIDPSSignatureSeverityThree),
		int64(FirewallPolicyIDPSSignatureSeverityTwo),
	}
}

type FirewallPolicyIntrusionDetectionProfileType string

const (
	FirewallPolicyIntrusionDetectionProfileTypeAdvanced FirewallPolicyIntrusionDetectionProfileType = "Advanced"
	FirewallPolicyIntrusionDetectionProfileTypeBasic    FirewallPolicyIntrusionDetectionProfileType = "Basic"
	FirewallPolicyIntrusionDetectionProfileTypeExtended FirewallPolicyIntrusionDetectionProfileType = "Extended"
	FirewallPolicyIntrusionDetectionProfileTypeStandard FirewallPolicyIntrusionDetectionProfileType = "Standard"
)

func PossibleValuesForFirewallPolicyIntrusionDetectionProfileType() []string {
	return []string{
		string(FirewallPolicyIntrusionDetectionProfileTypeAdvanced),
		string(FirewallPolicyIntrusionDetectionProfileTypeBasic),
		string(FirewallPolicyIntrusionDetectionProfileTypeExtended),
		string(FirewallPolicyIntrusionDetectionProfileTypeStandard),
	}
}

func (s *FirewallPolicyIntrusionDetectionProfileType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyIntrusionDetectionProfileType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyIntrusionDetectionProfileType(input string) (*FirewallPolicyIntrusionDetectionProfileType, error) {
	vals := map[string]FirewallPolicyIntrusionDetectionProfileType{
		"advanced": FirewallPolicyIntrusionDetectionProfileTypeAdvanced,
		"basic":    FirewallPolicyIntrusionDetectionProfileTypeBasic,
		"extended": FirewallPolicyIntrusionDetectionProfileTypeExtended,
		"standard": FirewallPolicyIntrusionDetectionProfileTypeStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyIntrusionDetectionProfileType(input)
	return &out, nil
}

type FirewallPolicyIntrusionDetectionProtocol string

const (
	FirewallPolicyIntrusionDetectionProtocolANY  FirewallPolicyIntrusionDetectionProtocol = "ANY"
	FirewallPolicyIntrusionDetectionProtocolICMP FirewallPolicyIntrusionDetectionProtocol = "ICMP"
	FirewallPolicyIntrusionDetectionProtocolTCP  FirewallPolicyIntrusionDetectionProtocol = "TCP"
	FirewallPolicyIntrusionDetectionProtocolUDP  FirewallPolicyIntrusionDetectionProtocol = "UDP"
)

func PossibleValuesForFirewallPolicyIntrusionDetectionProtocol() []string {
	return []string{
		string(FirewallPolicyIntrusionDetectionProtocolANY),
		string(FirewallPolicyIntrusionDetectionProtocolICMP),
		string(FirewallPolicyIntrusionDetectionProtocolTCP),
		string(FirewallPolicyIntrusionDetectionProtocolUDP),
	}
}

func (s *FirewallPolicyIntrusionDetectionProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyIntrusionDetectionProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyIntrusionDetectionProtocol(input string) (*FirewallPolicyIntrusionDetectionProtocol, error) {
	vals := map[string]FirewallPolicyIntrusionDetectionProtocol{
		"any":  FirewallPolicyIntrusionDetectionProtocolANY,
		"icmp": FirewallPolicyIntrusionDetectionProtocolICMP,
		"tcp":  FirewallPolicyIntrusionDetectionProtocolTCP,
		"udp":  FirewallPolicyIntrusionDetectionProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyIntrusionDetectionProtocol(input)
	return &out, nil
}

type FirewallPolicyIntrusionDetectionStateType string

const (
	FirewallPolicyIntrusionDetectionStateTypeAlert FirewallPolicyIntrusionDetectionStateType = "Alert"
	FirewallPolicyIntrusionDetectionStateTypeDeny  FirewallPolicyIntrusionDetectionStateType = "Deny"
	FirewallPolicyIntrusionDetectionStateTypeOff   FirewallPolicyIntrusionDetectionStateType = "Off"
)

func PossibleValuesForFirewallPolicyIntrusionDetectionStateType() []string {
	return []string{
		string(FirewallPolicyIntrusionDetectionStateTypeAlert),
		string(FirewallPolicyIntrusionDetectionStateTypeDeny),
		string(FirewallPolicyIntrusionDetectionStateTypeOff),
	}
}

func (s *FirewallPolicyIntrusionDetectionStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyIntrusionDetectionStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyIntrusionDetectionStateType(input string) (*FirewallPolicyIntrusionDetectionStateType, error) {
	vals := map[string]FirewallPolicyIntrusionDetectionStateType{
		"alert": FirewallPolicyIntrusionDetectionStateTypeAlert,
		"deny":  FirewallPolicyIntrusionDetectionStateTypeDeny,
		"off":   FirewallPolicyIntrusionDetectionStateTypeOff,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyIntrusionDetectionStateType(input)
	return &out, nil
}

type FirewallPolicyNatRuleCollectionActionType string

const (
	FirewallPolicyNatRuleCollectionActionTypeDNAT FirewallPolicyNatRuleCollectionActionType = "DNAT"
)

func PossibleValuesForFirewallPolicyNatRuleCollectionActionType() []string {
	return []string{
		string(FirewallPolicyNatRuleCollectionActionTypeDNAT),
	}
}

func (s *FirewallPolicyNatRuleCollectionActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyNatRuleCollectionActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyNatRuleCollectionActionType(input string) (*FirewallPolicyNatRuleCollectionActionType, error) {
	vals := map[string]FirewallPolicyNatRuleCollectionActionType{
		"dnat": FirewallPolicyNatRuleCollectionActionTypeDNAT,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyNatRuleCollectionActionType(input)
	return &out, nil
}

type FirewallPolicyRuleApplicationProtocolType string

const (
	FirewallPolicyRuleApplicationProtocolTypeHTTP  FirewallPolicyRuleApplicationProtocolType = "Http"
	FirewallPolicyRuleApplicationProtocolTypeHTTPS FirewallPolicyRuleApplicationProtocolType = "Https"
)

func PossibleValuesForFirewallPolicyRuleApplicationProtocolType() []string {
	return []string{
		string(FirewallPolicyRuleApplicationProtocolTypeHTTP),
		string(FirewallPolicyRuleApplicationProtocolTypeHTTPS),
	}
}

func (s *FirewallPolicyRuleApplicationProtocolType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyRuleApplicationProtocolType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyRuleApplicationProtocolType(input string) (*FirewallPolicyRuleApplicationProtocolType, error) {
	vals := map[string]FirewallPolicyRuleApplicationProtocolType{
		"http":  FirewallPolicyRuleApplicationProtocolTypeHTTP,
		"https": FirewallPolicyRuleApplicationProtocolTypeHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyRuleApplicationProtocolType(input)
	return &out, nil
}

type FirewallPolicyRuleCollectionType string

const (
	FirewallPolicyRuleCollectionTypeFirewallPolicyFilterRuleCollection FirewallPolicyRuleCollectionType = "FirewallPolicyFilterRuleCollection"
	FirewallPolicyRuleCollectionTypeFirewallPolicyNatRuleCollection    FirewallPolicyRuleCollectionType = "FirewallPolicyNatRuleCollection"
)

func PossibleValuesForFirewallPolicyRuleCollectionType() []string {
	return []string{
		string(FirewallPolicyRuleCollectionTypeFirewallPolicyFilterRuleCollection),
		string(FirewallPolicyRuleCollectionTypeFirewallPolicyNatRuleCollection),
	}
}

func (s *FirewallPolicyRuleCollectionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyRuleCollectionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyRuleCollectionType(input string) (*FirewallPolicyRuleCollectionType, error) {
	vals := map[string]FirewallPolicyRuleCollectionType{
		"firewallpolicyfilterrulecollection": FirewallPolicyRuleCollectionTypeFirewallPolicyFilterRuleCollection,
		"firewallpolicynatrulecollection":    FirewallPolicyRuleCollectionTypeFirewallPolicyNatRuleCollection,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyRuleCollectionType(input)
	return &out, nil
}

type FirewallPolicyRuleNetworkProtocol string

const (
	FirewallPolicyRuleNetworkProtocolAny  FirewallPolicyRuleNetworkProtocol = "Any"
	FirewallPolicyRuleNetworkProtocolICMP FirewallPolicyRuleNetworkProtocol = "ICMP"
	FirewallPolicyRuleNetworkProtocolTCP  FirewallPolicyRuleNetworkProtocol = "TCP"
	FirewallPolicyRuleNetworkProtocolUDP  FirewallPolicyRuleNetworkProtocol = "UDP"
)

func PossibleValuesForFirewallPolicyRuleNetworkProtocol() []string {
	return []string{
		string(FirewallPolicyRuleNetworkProtocolAny),
		string(FirewallPolicyRuleNetworkProtocolICMP),
		string(FirewallPolicyRuleNetworkProtocolTCP),
		string(FirewallPolicyRuleNetworkProtocolUDP),
	}
}

func (s *FirewallPolicyRuleNetworkProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyRuleNetworkProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyRuleNetworkProtocol(input string) (*FirewallPolicyRuleNetworkProtocol, error) {
	vals := map[string]FirewallPolicyRuleNetworkProtocol{
		"any":  FirewallPolicyRuleNetworkProtocolAny,
		"icmp": FirewallPolicyRuleNetworkProtocolICMP,
		"tcp":  FirewallPolicyRuleNetworkProtocolTCP,
		"udp":  FirewallPolicyRuleNetworkProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyRuleNetworkProtocol(input)
	return &out, nil
}

type FirewallPolicyRuleType string

const (
	FirewallPolicyRuleTypeApplicationRule FirewallPolicyRuleType = "ApplicationRule"
	FirewallPolicyRuleTypeNatRule         FirewallPolicyRuleType = "NatRule"
	FirewallPolicyRuleTypeNetworkRule     FirewallPolicyRuleType = "NetworkRule"
)

func PossibleValuesForFirewallPolicyRuleType() []string {
	return []string{
		string(FirewallPolicyRuleTypeApplicationRule),
		string(FirewallPolicyRuleTypeNatRule),
		string(FirewallPolicyRuleTypeNetworkRule),
	}
}

func (s *FirewallPolicyRuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicyRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicyRuleType(input string) (*FirewallPolicyRuleType, error) {
	vals := map[string]FirewallPolicyRuleType{
		"applicationrule": FirewallPolicyRuleTypeApplicationRule,
		"natrule":         FirewallPolicyRuleTypeNatRule,
		"networkrule":     FirewallPolicyRuleTypeNetworkRule,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicyRuleType(input)
	return &out, nil
}

type FirewallPolicySkuTier string

const (
	FirewallPolicySkuTierBasic    FirewallPolicySkuTier = "Basic"
	FirewallPolicySkuTierPremium  FirewallPolicySkuTier = "Premium"
	FirewallPolicySkuTierStandard FirewallPolicySkuTier = "Standard"
)

func PossibleValuesForFirewallPolicySkuTier() []string {
	return []string{
		string(FirewallPolicySkuTierBasic),
		string(FirewallPolicySkuTierPremium),
		string(FirewallPolicySkuTierStandard),
	}
}

func (s *FirewallPolicySkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFirewallPolicySkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFirewallPolicySkuTier(input string) (*FirewallPolicySkuTier, error) {
	vals := map[string]FirewallPolicySkuTier{
		"basic":    FirewallPolicySkuTierBasic,
		"premium":  FirewallPolicySkuTierPremium,
		"standard": FirewallPolicySkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FirewallPolicySkuTier(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
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
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
