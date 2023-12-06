package azurefirewalls

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFirewallApplicationRuleProtocolType string

const (
	AzureFirewallApplicationRuleProtocolTypeHTTP  AzureFirewallApplicationRuleProtocolType = "Http"
	AzureFirewallApplicationRuleProtocolTypeHTTPS AzureFirewallApplicationRuleProtocolType = "Https"
	AzureFirewallApplicationRuleProtocolTypeMssql AzureFirewallApplicationRuleProtocolType = "Mssql"
)

func PossibleValuesForAzureFirewallApplicationRuleProtocolType() []string {
	return []string{
		string(AzureFirewallApplicationRuleProtocolTypeHTTP),
		string(AzureFirewallApplicationRuleProtocolTypeHTTPS),
		string(AzureFirewallApplicationRuleProtocolTypeMssql),
	}
}

func (s *AzureFirewallApplicationRuleProtocolType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFirewallApplicationRuleProtocolType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFirewallApplicationRuleProtocolType(input string) (*AzureFirewallApplicationRuleProtocolType, error) {
	vals := map[string]AzureFirewallApplicationRuleProtocolType{
		"http":  AzureFirewallApplicationRuleProtocolTypeHTTP,
		"https": AzureFirewallApplicationRuleProtocolTypeHTTPS,
		"mssql": AzureFirewallApplicationRuleProtocolTypeMssql,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFirewallApplicationRuleProtocolType(input)
	return &out, nil
}

type AzureFirewallNatRCActionType string

const (
	AzureFirewallNatRCActionTypeDnat AzureFirewallNatRCActionType = "Dnat"
	AzureFirewallNatRCActionTypeSnat AzureFirewallNatRCActionType = "Snat"
)

func PossibleValuesForAzureFirewallNatRCActionType() []string {
	return []string{
		string(AzureFirewallNatRCActionTypeDnat),
		string(AzureFirewallNatRCActionTypeSnat),
	}
}

func (s *AzureFirewallNatRCActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFirewallNatRCActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFirewallNatRCActionType(input string) (*AzureFirewallNatRCActionType, error) {
	vals := map[string]AzureFirewallNatRCActionType{
		"dnat": AzureFirewallNatRCActionTypeDnat,
		"snat": AzureFirewallNatRCActionTypeSnat,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFirewallNatRCActionType(input)
	return &out, nil
}

type AzureFirewallNetworkRuleProtocol string

const (
	AzureFirewallNetworkRuleProtocolAny  AzureFirewallNetworkRuleProtocol = "Any"
	AzureFirewallNetworkRuleProtocolICMP AzureFirewallNetworkRuleProtocol = "ICMP"
	AzureFirewallNetworkRuleProtocolTCP  AzureFirewallNetworkRuleProtocol = "TCP"
	AzureFirewallNetworkRuleProtocolUDP  AzureFirewallNetworkRuleProtocol = "UDP"
)

func PossibleValuesForAzureFirewallNetworkRuleProtocol() []string {
	return []string{
		string(AzureFirewallNetworkRuleProtocolAny),
		string(AzureFirewallNetworkRuleProtocolICMP),
		string(AzureFirewallNetworkRuleProtocolTCP),
		string(AzureFirewallNetworkRuleProtocolUDP),
	}
}

func (s *AzureFirewallNetworkRuleProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFirewallNetworkRuleProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFirewallNetworkRuleProtocol(input string) (*AzureFirewallNetworkRuleProtocol, error) {
	vals := map[string]AzureFirewallNetworkRuleProtocol{
		"any":  AzureFirewallNetworkRuleProtocolAny,
		"icmp": AzureFirewallNetworkRuleProtocolICMP,
		"tcp":  AzureFirewallNetworkRuleProtocolTCP,
		"udp":  AzureFirewallNetworkRuleProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFirewallNetworkRuleProtocol(input)
	return &out, nil
}

type AzureFirewallPacketCaptureFlagsType string

const (
	AzureFirewallPacketCaptureFlagsTypeAck  AzureFirewallPacketCaptureFlagsType = "ack"
	AzureFirewallPacketCaptureFlagsTypeFin  AzureFirewallPacketCaptureFlagsType = "fin"
	AzureFirewallPacketCaptureFlagsTypePush AzureFirewallPacketCaptureFlagsType = "push"
	AzureFirewallPacketCaptureFlagsTypeRst  AzureFirewallPacketCaptureFlagsType = "rst"
	AzureFirewallPacketCaptureFlagsTypeSyn  AzureFirewallPacketCaptureFlagsType = "syn"
	AzureFirewallPacketCaptureFlagsTypeUrg  AzureFirewallPacketCaptureFlagsType = "urg"
)

func PossibleValuesForAzureFirewallPacketCaptureFlagsType() []string {
	return []string{
		string(AzureFirewallPacketCaptureFlagsTypeAck),
		string(AzureFirewallPacketCaptureFlagsTypeFin),
		string(AzureFirewallPacketCaptureFlagsTypePush),
		string(AzureFirewallPacketCaptureFlagsTypeRst),
		string(AzureFirewallPacketCaptureFlagsTypeSyn),
		string(AzureFirewallPacketCaptureFlagsTypeUrg),
	}
}

func (s *AzureFirewallPacketCaptureFlagsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFirewallPacketCaptureFlagsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFirewallPacketCaptureFlagsType(input string) (*AzureFirewallPacketCaptureFlagsType, error) {
	vals := map[string]AzureFirewallPacketCaptureFlagsType{
		"ack":  AzureFirewallPacketCaptureFlagsTypeAck,
		"fin":  AzureFirewallPacketCaptureFlagsTypeFin,
		"push": AzureFirewallPacketCaptureFlagsTypePush,
		"rst":  AzureFirewallPacketCaptureFlagsTypeRst,
		"syn":  AzureFirewallPacketCaptureFlagsTypeSyn,
		"urg":  AzureFirewallPacketCaptureFlagsTypeUrg,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFirewallPacketCaptureFlagsType(input)
	return &out, nil
}

type AzureFirewallRCActionType string

const (
	AzureFirewallRCActionTypeAllow AzureFirewallRCActionType = "Allow"
	AzureFirewallRCActionTypeDeny  AzureFirewallRCActionType = "Deny"
)

func PossibleValuesForAzureFirewallRCActionType() []string {
	return []string{
		string(AzureFirewallRCActionTypeAllow),
		string(AzureFirewallRCActionTypeDeny),
	}
}

func (s *AzureFirewallRCActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFirewallRCActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFirewallRCActionType(input string) (*AzureFirewallRCActionType, error) {
	vals := map[string]AzureFirewallRCActionType{
		"allow": AzureFirewallRCActionTypeAllow,
		"deny":  AzureFirewallRCActionTypeDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFirewallRCActionType(input)
	return &out, nil
}

type AzureFirewallSkuName string

const (
	AzureFirewallSkuNameAZFWHub  AzureFirewallSkuName = "AZFW_Hub"
	AzureFirewallSkuNameAZFWVNet AzureFirewallSkuName = "AZFW_VNet"
)

func PossibleValuesForAzureFirewallSkuName() []string {
	return []string{
		string(AzureFirewallSkuNameAZFWHub),
		string(AzureFirewallSkuNameAZFWVNet),
	}
}

func (s *AzureFirewallSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFirewallSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFirewallSkuName(input string) (*AzureFirewallSkuName, error) {
	vals := map[string]AzureFirewallSkuName{
		"azfw_hub":  AzureFirewallSkuNameAZFWHub,
		"azfw_vnet": AzureFirewallSkuNameAZFWVNet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFirewallSkuName(input)
	return &out, nil
}

type AzureFirewallSkuTier string

const (
	AzureFirewallSkuTierBasic    AzureFirewallSkuTier = "Basic"
	AzureFirewallSkuTierPremium  AzureFirewallSkuTier = "Premium"
	AzureFirewallSkuTierStandard AzureFirewallSkuTier = "Standard"
)

func PossibleValuesForAzureFirewallSkuTier() []string {
	return []string{
		string(AzureFirewallSkuTierBasic),
		string(AzureFirewallSkuTierPremium),
		string(AzureFirewallSkuTierStandard),
	}
}

func (s *AzureFirewallSkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureFirewallSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureFirewallSkuTier(input string) (*AzureFirewallSkuTier, error) {
	vals := map[string]AzureFirewallSkuTier{
		"basic":    AzureFirewallSkuTierBasic,
		"premium":  AzureFirewallSkuTierPremium,
		"standard": AzureFirewallSkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureFirewallSkuTier(input)
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
