package routingrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type RoutingRuleDestinationType string

const (
	RoutingRuleDestinationTypeAddressPrefix RoutingRuleDestinationType = "AddressPrefix"
	RoutingRuleDestinationTypeServiceTag    RoutingRuleDestinationType = "ServiceTag"
)

func PossibleValuesForRoutingRuleDestinationType() []string {
	return []string{
		string(RoutingRuleDestinationTypeAddressPrefix),
		string(RoutingRuleDestinationTypeServiceTag),
	}
}

func (s *RoutingRuleDestinationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoutingRuleDestinationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoutingRuleDestinationType(input string) (*RoutingRuleDestinationType, error) {
	vals := map[string]RoutingRuleDestinationType{
		"addressprefix": RoutingRuleDestinationTypeAddressPrefix,
		"servicetag":    RoutingRuleDestinationTypeServiceTag,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingRuleDestinationType(input)
	return &out, nil
}

type RoutingRuleNextHopType string

const (
	RoutingRuleNextHopTypeInternet              RoutingRuleNextHopType = "Internet"
	RoutingRuleNextHopTypeNoNextHop             RoutingRuleNextHopType = "NoNextHop"
	RoutingRuleNextHopTypeVirtualAppliance      RoutingRuleNextHopType = "VirtualAppliance"
	RoutingRuleNextHopTypeVirtualNetworkGateway RoutingRuleNextHopType = "VirtualNetworkGateway"
	RoutingRuleNextHopTypeVnetLocal             RoutingRuleNextHopType = "VnetLocal"
)

func PossibleValuesForRoutingRuleNextHopType() []string {
	return []string{
		string(RoutingRuleNextHopTypeInternet),
		string(RoutingRuleNextHopTypeNoNextHop),
		string(RoutingRuleNextHopTypeVirtualAppliance),
		string(RoutingRuleNextHopTypeVirtualNetworkGateway),
		string(RoutingRuleNextHopTypeVnetLocal),
	}
}

func (s *RoutingRuleNextHopType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoutingRuleNextHopType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoutingRuleNextHopType(input string) (*RoutingRuleNextHopType, error) {
	vals := map[string]RoutingRuleNextHopType{
		"internet":              RoutingRuleNextHopTypeInternet,
		"nonexthop":             RoutingRuleNextHopTypeNoNextHop,
		"virtualappliance":      RoutingRuleNextHopTypeVirtualAppliance,
		"virtualnetworkgateway": RoutingRuleNextHopTypeVirtualNetworkGateway,
		"vnetlocal":             RoutingRuleNextHopTypeVnetLocal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingRuleNextHopType(input)
	return &out, nil
}
