package routefilters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Access string

const (
	AccessAllow Access = "Allow"
	AccessDeny  Access = "Deny"
)

func PossibleValuesForAccess() []string {
	return []string{
		string(AccessAllow),
		string(AccessDeny),
	}
}

func (s *Access) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccess(input string) (*Access, error) {
	vals := map[string]Access{
		"allow": AccessAllow,
		"deny":  AccessDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Access(input)
	return &out, nil
}

type CircuitConnectionStatus string

const (
	CircuitConnectionStatusConnected    CircuitConnectionStatus = "Connected"
	CircuitConnectionStatusConnecting   CircuitConnectionStatus = "Connecting"
	CircuitConnectionStatusDisconnected CircuitConnectionStatus = "Disconnected"
)

func PossibleValuesForCircuitConnectionStatus() []string {
	return []string{
		string(CircuitConnectionStatusConnected),
		string(CircuitConnectionStatusConnecting),
		string(CircuitConnectionStatusDisconnected),
	}
}

func (s *CircuitConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCircuitConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCircuitConnectionStatus(input string) (*CircuitConnectionStatus, error) {
	vals := map[string]CircuitConnectionStatus{
		"connected":    CircuitConnectionStatusConnected,
		"connecting":   CircuitConnectionStatusConnecting,
		"disconnected": CircuitConnectionStatusDisconnected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CircuitConnectionStatus(input)
	return &out, nil
}

type ExpressRouteCircuitPeeringAdvertisedPublicPrefixState string

const (
	ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateConfigured       ExpressRouteCircuitPeeringAdvertisedPublicPrefixState = "Configured"
	ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateConfiguring      ExpressRouteCircuitPeeringAdvertisedPublicPrefixState = "Configuring"
	ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateNotConfigured    ExpressRouteCircuitPeeringAdvertisedPublicPrefixState = "NotConfigured"
	ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateValidationNeeded ExpressRouteCircuitPeeringAdvertisedPublicPrefixState = "ValidationNeeded"
)

func PossibleValuesForExpressRouteCircuitPeeringAdvertisedPublicPrefixState() []string {
	return []string{
		string(ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateConfigured),
		string(ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateConfiguring),
		string(ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateNotConfigured),
		string(ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateValidationNeeded),
	}
}

func (s *ExpressRouteCircuitPeeringAdvertisedPublicPrefixState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRouteCircuitPeeringAdvertisedPublicPrefixState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRouteCircuitPeeringAdvertisedPublicPrefixState(input string) (*ExpressRouteCircuitPeeringAdvertisedPublicPrefixState, error) {
	vals := map[string]ExpressRouteCircuitPeeringAdvertisedPublicPrefixState{
		"configured":       ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateConfigured,
		"configuring":      ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateConfiguring,
		"notconfigured":    ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateNotConfigured,
		"validationneeded": ExpressRouteCircuitPeeringAdvertisedPublicPrefixStateValidationNeeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRouteCircuitPeeringAdvertisedPublicPrefixState(input)
	return &out, nil
}

type ExpressRouteCircuitPeeringState string

const (
	ExpressRouteCircuitPeeringStateDisabled ExpressRouteCircuitPeeringState = "Disabled"
	ExpressRouteCircuitPeeringStateEnabled  ExpressRouteCircuitPeeringState = "Enabled"
)

func PossibleValuesForExpressRouteCircuitPeeringState() []string {
	return []string{
		string(ExpressRouteCircuitPeeringStateDisabled),
		string(ExpressRouteCircuitPeeringStateEnabled),
	}
}

func (s *ExpressRouteCircuitPeeringState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRouteCircuitPeeringState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRouteCircuitPeeringState(input string) (*ExpressRouteCircuitPeeringState, error) {
	vals := map[string]ExpressRouteCircuitPeeringState{
		"disabled": ExpressRouteCircuitPeeringStateDisabled,
		"enabled":  ExpressRouteCircuitPeeringStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRouteCircuitPeeringState(input)
	return &out, nil
}

type ExpressRoutePeeringState string

const (
	ExpressRoutePeeringStateDisabled ExpressRoutePeeringState = "Disabled"
	ExpressRoutePeeringStateEnabled  ExpressRoutePeeringState = "Enabled"
)

func PossibleValuesForExpressRoutePeeringState() []string {
	return []string{
		string(ExpressRoutePeeringStateDisabled),
		string(ExpressRoutePeeringStateEnabled),
	}
}

func (s *ExpressRoutePeeringState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRoutePeeringState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRoutePeeringState(input string) (*ExpressRoutePeeringState, error) {
	vals := map[string]ExpressRoutePeeringState{
		"disabled": ExpressRoutePeeringStateDisabled,
		"enabled":  ExpressRoutePeeringStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRoutePeeringState(input)
	return &out, nil
}

type ExpressRoutePeeringType string

const (
	ExpressRoutePeeringTypeAzurePrivatePeering ExpressRoutePeeringType = "AzurePrivatePeering"
	ExpressRoutePeeringTypeAzurePublicPeering  ExpressRoutePeeringType = "AzurePublicPeering"
	ExpressRoutePeeringTypeMicrosoftPeering    ExpressRoutePeeringType = "MicrosoftPeering"
)

func PossibleValuesForExpressRoutePeeringType() []string {
	return []string{
		string(ExpressRoutePeeringTypeAzurePrivatePeering),
		string(ExpressRoutePeeringTypeAzurePublicPeering),
		string(ExpressRoutePeeringTypeMicrosoftPeering),
	}
}

func (s *ExpressRoutePeeringType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRoutePeeringType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRoutePeeringType(input string) (*ExpressRoutePeeringType, error) {
	vals := map[string]ExpressRoutePeeringType{
		"azureprivatepeering": ExpressRoutePeeringTypeAzurePrivatePeering,
		"azurepublicpeering":  ExpressRoutePeeringTypeAzurePublicPeering,
		"microsoftpeering":    ExpressRoutePeeringTypeMicrosoftPeering,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRoutePeeringType(input)
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

type RouteFilterRuleType string

const (
	RouteFilterRuleTypeCommunity RouteFilterRuleType = "Community"
)

func PossibleValuesForRouteFilterRuleType() []string {
	return []string{
		string(RouteFilterRuleTypeCommunity),
	}
}

func (s *RouteFilterRuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRouteFilterRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRouteFilterRuleType(input string) (*RouteFilterRuleType, error) {
	vals := map[string]RouteFilterRuleType{
		"community": RouteFilterRuleTypeCommunity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RouteFilterRuleType(input)
	return &out, nil
}
