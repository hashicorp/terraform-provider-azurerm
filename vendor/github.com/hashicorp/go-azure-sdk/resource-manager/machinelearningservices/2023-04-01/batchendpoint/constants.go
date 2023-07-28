package batchendpoint

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointAuthMode string

const (
	EndpointAuthModeAADToken EndpointAuthMode = "AADToken"
	EndpointAuthModeAMLToken EndpointAuthMode = "AMLToken"
	EndpointAuthModeKey      EndpointAuthMode = "Key"
)

func PossibleValuesForEndpointAuthMode() []string {
	return []string{
		string(EndpointAuthModeAADToken),
		string(EndpointAuthModeAMLToken),
		string(EndpointAuthModeKey),
	}
}

func (s *EndpointAuthMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointAuthMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEndpointAuthMode(input string) (*EndpointAuthMode, error) {
	vals := map[string]EndpointAuthMode{
		"aadtoken": EndpointAuthModeAADToken,
		"amltoken": EndpointAuthModeAMLToken,
		"key":      EndpointAuthModeKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointAuthMode(input)
	return &out, nil
}

type EndpointProvisioningState string

const (
	EndpointProvisioningStateCanceled  EndpointProvisioningState = "Canceled"
	EndpointProvisioningStateCreating  EndpointProvisioningState = "Creating"
	EndpointProvisioningStateDeleting  EndpointProvisioningState = "Deleting"
	EndpointProvisioningStateFailed    EndpointProvisioningState = "Failed"
	EndpointProvisioningStateSucceeded EndpointProvisioningState = "Succeeded"
	EndpointProvisioningStateUpdating  EndpointProvisioningState = "Updating"
)

func PossibleValuesForEndpointProvisioningState() []string {
	return []string{
		string(EndpointProvisioningStateCanceled),
		string(EndpointProvisioningStateCreating),
		string(EndpointProvisioningStateDeleting),
		string(EndpointProvisioningStateFailed),
		string(EndpointProvisioningStateSucceeded),
		string(EndpointProvisioningStateUpdating),
	}
}

func (s *EndpointProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEndpointProvisioningState(input string) (*EndpointProvisioningState, error) {
	vals := map[string]EndpointProvisioningState{
		"canceled":  EndpointProvisioningStateCanceled,
		"creating":  EndpointProvisioningStateCreating,
		"deleting":  EndpointProvisioningStateDeleting,
		"failed":    EndpointProvisioningStateFailed,
		"succeeded": EndpointProvisioningStateSucceeded,
		"updating":  EndpointProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointProvisioningState(input)
	return &out, nil
}

type ManagedServiceIdentityType string

const (
	ManagedServiceIdentityTypeNone                       ManagedServiceIdentityType = "None"
	ManagedServiceIdentityTypeSystemAssigned             ManagedServiceIdentityType = "SystemAssigned"
	ManagedServiceIdentityTypeSystemAssignedUserAssigned ManagedServiceIdentityType = "SystemAssigned,UserAssigned"
	ManagedServiceIdentityTypeUserAssigned               ManagedServiceIdentityType = "UserAssigned"
)

func PossibleValuesForManagedServiceIdentityType() []string {
	return []string{
		string(ManagedServiceIdentityTypeNone),
		string(ManagedServiceIdentityTypeSystemAssigned),
		string(ManagedServiceIdentityTypeSystemAssignedUserAssigned),
		string(ManagedServiceIdentityTypeUserAssigned),
	}
}

func (s *ManagedServiceIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedServiceIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedServiceIdentityType(input string) (*ManagedServiceIdentityType, error) {
	vals := map[string]ManagedServiceIdentityType{
		"none":                        ManagedServiceIdentityTypeNone,
		"systemassigned":              ManagedServiceIdentityTypeSystemAssigned,
		"systemassigned,userassigned": ManagedServiceIdentityTypeSystemAssignedUserAssigned,
		"userassigned":                ManagedServiceIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedServiceIdentityType(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierFree),
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    SkuTierBasic,
		"free":     SkuTierFree,
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}
