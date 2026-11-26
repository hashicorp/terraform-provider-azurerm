package virtualnetworkappliancecapabilities

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

type VirtualNetworkApplianceCapabilityIPVersion string

const (
	VirtualNetworkApplianceCapabilityIPVersionDualStack VirtualNetworkApplianceCapabilityIPVersion = "DualStack"
	VirtualNetworkApplianceCapabilityIPVersionIPvSix    VirtualNetworkApplianceCapabilityIPVersion = "IPv6"
)

func PossibleValuesForVirtualNetworkApplianceCapabilityIPVersion() []string {
	return []string{
		string(VirtualNetworkApplianceCapabilityIPVersionDualStack),
		string(VirtualNetworkApplianceCapabilityIPVersionIPvSix),
	}
}

func (s *VirtualNetworkApplianceCapabilityIPVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkApplianceCapabilityIPVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkApplianceCapabilityIPVersion(input string) (*VirtualNetworkApplianceCapabilityIPVersion, error) {
	vals := map[string]VirtualNetworkApplianceCapabilityIPVersion{
		"dualstack": VirtualNetworkApplianceCapabilityIPVersionDualStack,
		"ipv6":      VirtualNetworkApplianceCapabilityIPVersionIPvSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkApplianceCapabilityIPVersion(input)
	return &out, nil
}

type VirtualNetworkApplianceCapabilityKind string

const (
	VirtualNetworkApplianceCapabilityKindNATSixFour        VirtualNetworkApplianceCapabilityKind = "NAT64"
	VirtualNetworkApplianceCapabilityKindPLGateway         VirtualNetworkApplianceCapabilityKind = "PLGateway"
	VirtualNetworkApplianceCapabilityKindPLGatewayFastpath VirtualNetworkApplianceCapabilityKind = "PLGatewayFastpath"
	VirtualNetworkApplianceCapabilityKindPLIPForwarders    VirtualNetworkApplianceCapabilityKind = "PLIPForwarders"
)

func PossibleValuesForVirtualNetworkApplianceCapabilityKind() []string {
	return []string{
		string(VirtualNetworkApplianceCapabilityKindNATSixFour),
		string(VirtualNetworkApplianceCapabilityKindPLGateway),
		string(VirtualNetworkApplianceCapabilityKindPLGatewayFastpath),
		string(VirtualNetworkApplianceCapabilityKindPLIPForwarders),
	}
}

func (s *VirtualNetworkApplianceCapabilityKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkApplianceCapabilityKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkApplianceCapabilityKind(input string) (*VirtualNetworkApplianceCapabilityKind, error) {
	vals := map[string]VirtualNetworkApplianceCapabilityKind{
		"nat64":             VirtualNetworkApplianceCapabilityKindNATSixFour,
		"plgateway":         VirtualNetworkApplianceCapabilityKindPLGateway,
		"plgatewayfastpath": VirtualNetworkApplianceCapabilityKindPLGatewayFastpath,
		"plipforwarders":    VirtualNetworkApplianceCapabilityKindPLIPForwarders,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkApplianceCapabilityKind(input)
	return &out, nil
}
