package networkmanagereffectivesecurityadminrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddressPrefixType string

const (
	AddressPrefixTypeIPPrefix     AddressPrefixType = "IPPrefix"
	AddressPrefixTypeNetworkGroup AddressPrefixType = "NetworkGroup"
	AddressPrefixTypeServiceTag   AddressPrefixType = "ServiceTag"
)

func PossibleValuesForAddressPrefixType() []string {
	return []string{
		string(AddressPrefixTypeIPPrefix),
		string(AddressPrefixTypeNetworkGroup),
		string(AddressPrefixTypeServiceTag),
	}
}

func (s *AddressPrefixType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAddressPrefixType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAddressPrefixType(input string) (*AddressPrefixType, error) {
	vals := map[string]AddressPrefixType{
		"ipprefix":     AddressPrefixTypeIPPrefix,
		"networkgroup": AddressPrefixTypeNetworkGroup,
		"servicetag":   AddressPrefixTypeServiceTag,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AddressPrefixType(input)
	return &out, nil
}

type EffectiveAdminRuleKind string

const (
	EffectiveAdminRuleKindCustom  EffectiveAdminRuleKind = "Custom"
	EffectiveAdminRuleKindDefault EffectiveAdminRuleKind = "Default"
)

func PossibleValuesForEffectiveAdminRuleKind() []string {
	return []string{
		string(EffectiveAdminRuleKindCustom),
		string(EffectiveAdminRuleKindDefault),
	}
}

func (s *EffectiveAdminRuleKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEffectiveAdminRuleKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEffectiveAdminRuleKind(input string) (*EffectiveAdminRuleKind, error) {
	vals := map[string]EffectiveAdminRuleKind{
		"custom":  EffectiveAdminRuleKindCustom,
		"default": EffectiveAdminRuleKindDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EffectiveAdminRuleKind(input)
	return &out, nil
}

type GroupMemberType string

const (
	GroupMemberTypeSubnet         GroupMemberType = "Subnet"
	GroupMemberTypeVirtualNetwork GroupMemberType = "VirtualNetwork"
)

func PossibleValuesForGroupMemberType() []string {
	return []string{
		string(GroupMemberTypeSubnet),
		string(GroupMemberTypeVirtualNetwork),
	}
}

func (s *GroupMemberType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGroupMemberType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGroupMemberType(input string) (*GroupMemberType, error) {
	vals := map[string]GroupMemberType{
		"subnet":         GroupMemberTypeSubnet,
		"virtualnetwork": GroupMemberTypeVirtualNetwork,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GroupMemberType(input)
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

type SecurityConfigurationRuleAccess string

const (
	SecurityConfigurationRuleAccessAllow       SecurityConfigurationRuleAccess = "Allow"
	SecurityConfigurationRuleAccessAlwaysAllow SecurityConfigurationRuleAccess = "AlwaysAllow"
	SecurityConfigurationRuleAccessDeny        SecurityConfigurationRuleAccess = "Deny"
)

func PossibleValuesForSecurityConfigurationRuleAccess() []string {
	return []string{
		string(SecurityConfigurationRuleAccessAllow),
		string(SecurityConfigurationRuleAccessAlwaysAllow),
		string(SecurityConfigurationRuleAccessDeny),
	}
}

func (s *SecurityConfigurationRuleAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityConfigurationRuleAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityConfigurationRuleAccess(input string) (*SecurityConfigurationRuleAccess, error) {
	vals := map[string]SecurityConfigurationRuleAccess{
		"allow":       SecurityConfigurationRuleAccessAllow,
		"alwaysallow": SecurityConfigurationRuleAccessAlwaysAllow,
		"deny":        SecurityConfigurationRuleAccessDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityConfigurationRuleAccess(input)
	return &out, nil
}

type SecurityConfigurationRuleDirection string

const (
	SecurityConfigurationRuleDirectionInbound  SecurityConfigurationRuleDirection = "Inbound"
	SecurityConfigurationRuleDirectionOutbound SecurityConfigurationRuleDirection = "Outbound"
)

func PossibleValuesForSecurityConfigurationRuleDirection() []string {
	return []string{
		string(SecurityConfigurationRuleDirectionInbound),
		string(SecurityConfigurationRuleDirectionOutbound),
	}
}

func (s *SecurityConfigurationRuleDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityConfigurationRuleDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityConfigurationRuleDirection(input string) (*SecurityConfigurationRuleDirection, error) {
	vals := map[string]SecurityConfigurationRuleDirection{
		"inbound":  SecurityConfigurationRuleDirectionInbound,
		"outbound": SecurityConfigurationRuleDirectionOutbound,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityConfigurationRuleDirection(input)
	return &out, nil
}

type SecurityConfigurationRuleProtocol string

const (
	SecurityConfigurationRuleProtocolAh   SecurityConfigurationRuleProtocol = "Ah"
	SecurityConfigurationRuleProtocolAny  SecurityConfigurationRuleProtocol = "Any"
	SecurityConfigurationRuleProtocolEsp  SecurityConfigurationRuleProtocol = "Esp"
	SecurityConfigurationRuleProtocolIcmp SecurityConfigurationRuleProtocol = "Icmp"
	SecurityConfigurationRuleProtocolTcp  SecurityConfigurationRuleProtocol = "Tcp"
	SecurityConfigurationRuleProtocolUdp  SecurityConfigurationRuleProtocol = "Udp"
)

func PossibleValuesForSecurityConfigurationRuleProtocol() []string {
	return []string{
		string(SecurityConfigurationRuleProtocolAh),
		string(SecurityConfigurationRuleProtocolAny),
		string(SecurityConfigurationRuleProtocolEsp),
		string(SecurityConfigurationRuleProtocolIcmp),
		string(SecurityConfigurationRuleProtocolTcp),
		string(SecurityConfigurationRuleProtocolUdp),
	}
}

func (s *SecurityConfigurationRuleProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityConfigurationRuleProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityConfigurationRuleProtocol(input string) (*SecurityConfigurationRuleProtocol, error) {
	vals := map[string]SecurityConfigurationRuleProtocol{
		"ah":   SecurityConfigurationRuleProtocolAh,
		"any":  SecurityConfigurationRuleProtocolAny,
		"esp":  SecurityConfigurationRuleProtocolEsp,
		"icmp": SecurityConfigurationRuleProtocolIcmp,
		"tcp":  SecurityConfigurationRuleProtocolTcp,
		"udp":  SecurityConfigurationRuleProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityConfigurationRuleProtocol(input)
	return &out, nil
}
