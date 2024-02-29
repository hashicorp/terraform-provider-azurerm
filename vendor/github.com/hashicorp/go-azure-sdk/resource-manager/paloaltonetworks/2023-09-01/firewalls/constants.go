package firewalls

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingCycle string

const (
	BillingCycleMONTHLY BillingCycle = "MONTHLY"
	BillingCycleWEEKLY  BillingCycle = "WEEKLY"
)

func PossibleValuesForBillingCycle() []string {
	return []string{
		string(BillingCycleMONTHLY),
		string(BillingCycleWEEKLY),
	}
}

func (s *BillingCycle) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBillingCycle(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBillingCycle(input string) (*BillingCycle, error) {
	vals := map[string]BillingCycle{
		"monthly": BillingCycleMONTHLY,
		"weekly":  BillingCycleWEEKLY,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BillingCycle(input)
	return &out, nil
}

type BooleanEnum string

const (
	BooleanEnumFALSE BooleanEnum = "FALSE"
	BooleanEnumTRUE  BooleanEnum = "TRUE"
)

func PossibleValuesForBooleanEnum() []string {
	return []string{
		string(BooleanEnumFALSE),
		string(BooleanEnumTRUE),
	}
}

func (s *BooleanEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBooleanEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBooleanEnum(input string) (*BooleanEnum, error) {
	vals := map[string]BooleanEnum{
		"false": BooleanEnumFALSE,
		"true":  BooleanEnumTRUE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BooleanEnum(input)
	return &out, nil
}

type DNSProxy string

const (
	DNSProxyDISABLED DNSProxy = "DISABLED"
	DNSProxyENABLED  DNSProxy = "ENABLED"
)

func PossibleValuesForDNSProxy() []string {
	return []string{
		string(DNSProxyDISABLED),
		string(DNSProxyENABLED),
	}
}

func (s *DNSProxy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDNSProxy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDNSProxy(input string) (*DNSProxy, error) {
	vals := map[string]DNSProxy{
		"disabled": DNSProxyDISABLED,
		"enabled":  DNSProxyENABLED,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DNSProxy(input)
	return &out, nil
}

type EgressNat string

const (
	EgressNatDISABLED EgressNat = "DISABLED"
	EgressNatENABLED  EgressNat = "ENABLED"
)

func PossibleValuesForEgressNat() []string {
	return []string{
		string(EgressNatDISABLED),
		string(EgressNatENABLED),
	}
}

func (s *EgressNat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEgressNat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEgressNat(input string) (*EgressNat, error) {
	vals := map[string]EgressNat{
		"disabled": EgressNatDISABLED,
		"enabled":  EgressNatENABLED,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EgressNat(input)
	return &out, nil
}

type EnabledDNSType string

const (
	EnabledDNSTypeAZURE  EnabledDNSType = "AZURE"
	EnabledDNSTypeCUSTOM EnabledDNSType = "CUSTOM"
)

func PossibleValuesForEnabledDNSType() []string {
	return []string{
		string(EnabledDNSTypeAZURE),
		string(EnabledDNSTypeCUSTOM),
	}
}

func (s *EnabledDNSType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnabledDNSType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnabledDNSType(input string) (*EnabledDNSType, error) {
	vals := map[string]EnabledDNSType{
		"azure":  EnabledDNSTypeAZURE,
		"custom": EnabledDNSTypeCUSTOM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnabledDNSType(input)
	return &out, nil
}

type LogOption string

const (
	LogOptionINDIVIDUALDESTINATION LogOption = "INDIVIDUAL_DESTINATION"
	LogOptionSAMEDESTINATION       LogOption = "SAME_DESTINATION"
)

func PossibleValuesForLogOption() []string {
	return []string{
		string(LogOptionINDIVIDUALDESTINATION),
		string(LogOptionSAMEDESTINATION),
	}
}

func (s *LogOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLogOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLogOption(input string) (*LogOption, error) {
	vals := map[string]LogOption{
		"individual_destination": LogOptionINDIVIDUALDESTINATION,
		"same_destination":       LogOptionSAMEDESTINATION,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LogOption(input)
	return &out, nil
}

type LogType string

const (
	LogTypeAUDIT      LogType = "AUDIT"
	LogTypeDECRYPTION LogType = "DECRYPTION"
	LogTypeDLP        LogType = "DLP"
	LogTypeTHREAT     LogType = "THREAT"
	LogTypeTRAFFIC    LogType = "TRAFFIC"
	LogTypeWILDFIRE   LogType = "WILDFIRE"
)

func PossibleValuesForLogType() []string {
	return []string{
		string(LogTypeAUDIT),
		string(LogTypeDECRYPTION),
		string(LogTypeDLP),
		string(LogTypeTHREAT),
		string(LogTypeTRAFFIC),
		string(LogTypeWILDFIRE),
	}
}

func (s *LogType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLogType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLogType(input string) (*LogType, error) {
	vals := map[string]LogType{
		"audit":      LogTypeAUDIT,
		"decryption": LogTypeDECRYPTION,
		"dlp":        LogTypeDLP,
		"threat":     LogTypeTHREAT,
		"traffic":    LogTypeTRAFFIC,
		"wildfire":   LogTypeWILDFIRE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LogType(input)
	return &out, nil
}

type MarketplaceSubscriptionStatus string

const (
	MarketplaceSubscriptionStatusFulfillmentRequested    MarketplaceSubscriptionStatus = "FulfillmentRequested"
	MarketplaceSubscriptionStatusNotStarted              MarketplaceSubscriptionStatus = "NotStarted"
	MarketplaceSubscriptionStatusPendingFulfillmentStart MarketplaceSubscriptionStatus = "PendingFulfillmentStart"
	MarketplaceSubscriptionStatusSubscribed              MarketplaceSubscriptionStatus = "Subscribed"
	MarketplaceSubscriptionStatusSuspended               MarketplaceSubscriptionStatus = "Suspended"
	MarketplaceSubscriptionStatusUnsubscribed            MarketplaceSubscriptionStatus = "Unsubscribed"
)

func PossibleValuesForMarketplaceSubscriptionStatus() []string {
	return []string{
		string(MarketplaceSubscriptionStatusFulfillmentRequested),
		string(MarketplaceSubscriptionStatusNotStarted),
		string(MarketplaceSubscriptionStatusPendingFulfillmentStart),
		string(MarketplaceSubscriptionStatusSubscribed),
		string(MarketplaceSubscriptionStatusSuspended),
		string(MarketplaceSubscriptionStatusUnsubscribed),
	}
}

func (s *MarketplaceSubscriptionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMarketplaceSubscriptionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMarketplaceSubscriptionStatus(input string) (*MarketplaceSubscriptionStatus, error) {
	vals := map[string]MarketplaceSubscriptionStatus{
		"fulfillmentrequested":    MarketplaceSubscriptionStatusFulfillmentRequested,
		"notstarted":              MarketplaceSubscriptionStatusNotStarted,
		"pendingfulfillmentstart": MarketplaceSubscriptionStatusPendingFulfillmentStart,
		"subscribed":              MarketplaceSubscriptionStatusSubscribed,
		"suspended":               MarketplaceSubscriptionStatusSuspended,
		"unsubscribed":            MarketplaceSubscriptionStatusUnsubscribed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MarketplaceSubscriptionStatus(input)
	return &out, nil
}

type NetworkType string

const (
	NetworkTypeVNET NetworkType = "VNET"
	NetworkTypeVWAN NetworkType = "VWAN"
)

func PossibleValuesForNetworkType() []string {
	return []string{
		string(NetworkTypeVNET),
		string(NetworkTypeVWAN),
	}
}

func (s *NetworkType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkType(input string) (*NetworkType, error) {
	vals := map[string]NetworkType{
		"vnet": NetworkTypeVNET,
		"vwan": NetworkTypeVWAN,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkType(input)
	return &out, nil
}

type ProtocolType string

const (
	ProtocolTypeTCP ProtocolType = "TCP"
	ProtocolTypeUDP ProtocolType = "UDP"
)

func PossibleValuesForProtocolType() []string {
	return []string{
		string(ProtocolTypeTCP),
		string(ProtocolTypeUDP),
	}
}

func (s *ProtocolType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProtocolType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProtocolType(input string) (*ProtocolType, error) {
	vals := map[string]ProtocolType{
		"tcp": ProtocolTypeTCP,
		"udp": ProtocolTypeUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProtocolType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleted      ProvisioningState = "Deleted"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateNotSpecified),
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
		"creating":     ProvisioningStateCreating,
		"deleted":      ProvisioningStateDeleted,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"notspecified": ProvisioningStateNotSpecified,
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

type UsageType string

const (
	UsageTypeCOMMITTED UsageType = "COMMITTED"
	UsageTypePAYG      UsageType = "PAYG"
)

func PossibleValuesForUsageType() []string {
	return []string{
		string(UsageTypeCOMMITTED),
		string(UsageTypePAYG),
	}
}

func (s *UsageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUsageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUsageType(input string) (*UsageType, error) {
	vals := map[string]UsageType{
		"committed": UsageTypeCOMMITTED,
		"payg":      UsageTypePAYG,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsageType(input)
	return &out, nil
}
