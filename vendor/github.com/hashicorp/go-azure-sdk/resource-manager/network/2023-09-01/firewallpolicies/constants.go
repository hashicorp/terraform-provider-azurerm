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
