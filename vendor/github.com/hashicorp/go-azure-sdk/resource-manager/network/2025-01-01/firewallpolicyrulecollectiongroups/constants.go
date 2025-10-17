package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
