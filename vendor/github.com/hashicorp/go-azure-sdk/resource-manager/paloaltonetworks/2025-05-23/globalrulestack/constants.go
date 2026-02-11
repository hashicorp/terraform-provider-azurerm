package globalrulestack

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvSecurityObjectTypeEnum string

const (
	AdvSecurityObjectTypeEnumFeeds     AdvSecurityObjectTypeEnum = "feeds"
	AdvSecurityObjectTypeEnumUrlCustom AdvSecurityObjectTypeEnum = "urlCustom"
)

func PossibleValuesForAdvSecurityObjectTypeEnum() []string {
	return []string{
		string(AdvSecurityObjectTypeEnumFeeds),
		string(AdvSecurityObjectTypeEnumUrlCustom),
	}
}

func (s *AdvSecurityObjectTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdvSecurityObjectTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdvSecurityObjectTypeEnum(input string) (*AdvSecurityObjectTypeEnum, error) {
	vals := map[string]AdvSecurityObjectTypeEnum{
		"feeds":     AdvSecurityObjectTypeEnumFeeds,
		"urlcustom": AdvSecurityObjectTypeEnumUrlCustom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdvSecurityObjectTypeEnum(input)
	return &out, nil
}

type DefaultMode string

const (
	DefaultModeFIREWALL DefaultMode = "FIREWALL"
	DefaultModeIPS      DefaultMode = "IPS"
	DefaultModeNONE     DefaultMode = "NONE"
)

func PossibleValuesForDefaultMode() []string {
	return []string{
		string(DefaultModeFIREWALL),
		string(DefaultModeIPS),
		string(DefaultModeNONE),
	}
}

func (s *DefaultMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDefaultMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDefaultMode(input string) (*DefaultMode, error) {
	vals := map[string]DefaultMode{
		"firewall": DefaultModeFIREWALL,
		"ips":      DefaultModeIPS,
		"none":     DefaultModeNONE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DefaultMode(input)
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

type ScopeType string

const (
	ScopeTypeGLOBAL ScopeType = "GLOBAL"
	ScopeTypeLOCAL  ScopeType = "LOCAL"
)

func PossibleValuesForScopeType() []string {
	return []string{
		string(ScopeTypeGLOBAL),
		string(ScopeTypeLOCAL),
	}
}

func (s *ScopeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScopeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScopeType(input string) (*ScopeType, error) {
	vals := map[string]ScopeType{
		"global": ScopeTypeGLOBAL,
		"local":  ScopeTypeLOCAL,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScopeType(input)
	return &out, nil
}

type SecurityServicesTypeEnum string

const (
	SecurityServicesTypeEnumAntiSpyware      SecurityServicesTypeEnum = "antiSpyware"
	SecurityServicesTypeEnumAntiVirus        SecurityServicesTypeEnum = "antiVirus"
	SecurityServicesTypeEnumDnsSubscription  SecurityServicesTypeEnum = "dnsSubscription"
	SecurityServicesTypeEnumFileBlocking     SecurityServicesTypeEnum = "fileBlocking"
	SecurityServicesTypeEnumIPsVulnerability SecurityServicesTypeEnum = "ipsVulnerability"
	SecurityServicesTypeEnumUrlFiltering     SecurityServicesTypeEnum = "urlFiltering"
)

func PossibleValuesForSecurityServicesTypeEnum() []string {
	return []string{
		string(SecurityServicesTypeEnumAntiSpyware),
		string(SecurityServicesTypeEnumAntiVirus),
		string(SecurityServicesTypeEnumDnsSubscription),
		string(SecurityServicesTypeEnumFileBlocking),
		string(SecurityServicesTypeEnumIPsVulnerability),
		string(SecurityServicesTypeEnumUrlFiltering),
	}
}

func (s *SecurityServicesTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityServicesTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityServicesTypeEnum(input string) (*SecurityServicesTypeEnum, error) {
	vals := map[string]SecurityServicesTypeEnum{
		"antispyware":      SecurityServicesTypeEnumAntiSpyware,
		"antivirus":        SecurityServicesTypeEnumAntiVirus,
		"dnssubscription":  SecurityServicesTypeEnumDnsSubscription,
		"fileblocking":     SecurityServicesTypeEnumFileBlocking,
		"ipsvulnerability": SecurityServicesTypeEnumIPsVulnerability,
		"urlfiltering":     SecurityServicesTypeEnumUrlFiltering,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityServicesTypeEnum(input)
	return &out, nil
}
