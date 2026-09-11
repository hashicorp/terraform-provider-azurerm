package networkgateways

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
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
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
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

type VirtualNetworkGatewayMigrationType string

const (
	VirtualNetworkGatewayMigrationTypeUpgradeDeploymentToStandardIP VirtualNetworkGatewayMigrationType = "UpgradeDeploymentToStandardIP"
)

func PossibleValuesForVirtualNetworkGatewayMigrationType() []string {
	return []string{
		string(VirtualNetworkGatewayMigrationTypeUpgradeDeploymentToStandardIP),
	}
}

func (s *VirtualNetworkGatewayMigrationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkGatewayMigrationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkGatewayMigrationType(input string) (*VirtualNetworkGatewayMigrationType, error) {
	vals := map[string]VirtualNetworkGatewayMigrationType{
		"upgradedeploymenttostandardip": VirtualNetworkGatewayMigrationTypeUpgradeDeploymentToStandardIP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkGatewayMigrationType(input)
	return &out, nil
}

type VpnNatRuleMode string

const (
	VpnNatRuleModeEgressSnat  VpnNatRuleMode = "EgressSnat"
	VpnNatRuleModeIngressSnat VpnNatRuleMode = "IngressSnat"
)

func PossibleValuesForVpnNatRuleMode() []string {
	return []string{
		string(VpnNatRuleModeEgressSnat),
		string(VpnNatRuleModeIngressSnat),
	}
}

func (s *VpnNatRuleMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnNatRuleMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnNatRuleMode(input string) (*VpnNatRuleMode, error) {
	vals := map[string]VpnNatRuleMode{
		"egresssnat":  VpnNatRuleModeEgressSnat,
		"ingresssnat": VpnNatRuleModeIngressSnat,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnNatRuleMode(input)
	return &out, nil
}

type VpnNatRuleType string

const (
	VpnNatRuleTypeDynamic VpnNatRuleType = "Dynamic"
	VpnNatRuleTypeStatic  VpnNatRuleType = "Static"
)

func PossibleValuesForVpnNatRuleType() []string {
	return []string{
		string(VpnNatRuleTypeDynamic),
		string(VpnNatRuleTypeStatic),
	}
}

func (s *VpnNatRuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnNatRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnNatRuleType(input string) (*VpnNatRuleType, error) {
	vals := map[string]VpnNatRuleType{
		"dynamic": VpnNatRuleTypeDynamic,
		"static":  VpnNatRuleTypeStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnNatRuleType(input)
	return &out, nil
}
