package p2svpngateways

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationMethod string

const (
	AuthenticationMethodEAPMSCHAPvTwo AuthenticationMethod = "EAPMSCHAPv2"
	AuthenticationMethodEAPTLS        AuthenticationMethod = "EAPTLS"
)

func PossibleValuesForAuthenticationMethod() []string {
	return []string{
		string(AuthenticationMethodEAPMSCHAPvTwo),
		string(AuthenticationMethodEAPTLS),
	}
}

func (s *AuthenticationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthenticationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthenticationMethod(input string) (*AuthenticationMethod, error) {
	vals := map[string]AuthenticationMethod{
		"eapmschapv2": AuthenticationMethodEAPMSCHAPvTwo,
		"eaptls":      AuthenticationMethodEAPTLS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationMethod(input)
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

type VnetLocalRouteOverrideCriteria string

const (
	VnetLocalRouteOverrideCriteriaContains VnetLocalRouteOverrideCriteria = "Contains"
	VnetLocalRouteOverrideCriteriaEqual    VnetLocalRouteOverrideCriteria = "Equal"
)

func PossibleValuesForVnetLocalRouteOverrideCriteria() []string {
	return []string{
		string(VnetLocalRouteOverrideCriteriaContains),
		string(VnetLocalRouteOverrideCriteriaEqual),
	}
}

func (s *VnetLocalRouteOverrideCriteria) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVnetLocalRouteOverrideCriteria(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVnetLocalRouteOverrideCriteria(input string) (*VnetLocalRouteOverrideCriteria, error) {
	vals := map[string]VnetLocalRouteOverrideCriteria{
		"contains": VnetLocalRouteOverrideCriteriaContains,
		"equal":    VnetLocalRouteOverrideCriteriaEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VnetLocalRouteOverrideCriteria(input)
	return &out, nil
}

type VpnPolicyMemberAttributeType string

const (
	VpnPolicyMemberAttributeTypeAADGroupId         VpnPolicyMemberAttributeType = "AADGroupId"
	VpnPolicyMemberAttributeTypeCertificateGroupId VpnPolicyMemberAttributeType = "CertificateGroupId"
	VpnPolicyMemberAttributeTypeRadiusAzureGroupId VpnPolicyMemberAttributeType = "RadiusAzureGroupId"
)

func PossibleValuesForVpnPolicyMemberAttributeType() []string {
	return []string{
		string(VpnPolicyMemberAttributeTypeAADGroupId),
		string(VpnPolicyMemberAttributeTypeCertificateGroupId),
		string(VpnPolicyMemberAttributeTypeRadiusAzureGroupId),
	}
}

func (s *VpnPolicyMemberAttributeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnPolicyMemberAttributeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnPolicyMemberAttributeType(input string) (*VpnPolicyMemberAttributeType, error) {
	vals := map[string]VpnPolicyMemberAttributeType{
		"aadgroupid":         VpnPolicyMemberAttributeTypeAADGroupId,
		"certificategroupid": VpnPolicyMemberAttributeTypeCertificateGroupId,
		"radiusazuregroupid": VpnPolicyMemberAttributeTypeRadiusAzureGroupId,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnPolicyMemberAttributeType(input)
	return &out, nil
}
