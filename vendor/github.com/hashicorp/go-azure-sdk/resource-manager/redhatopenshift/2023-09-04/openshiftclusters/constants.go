package openshiftclusters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionAtHost string

const (
	EncryptionAtHostDisabled EncryptionAtHost = "Disabled"
	EncryptionAtHostEnabled  EncryptionAtHost = "Enabled"
)

func PossibleValuesForEncryptionAtHost() []string {
	return []string{
		string(EncryptionAtHostDisabled),
		string(EncryptionAtHostEnabled),
	}
}

func (s *EncryptionAtHost) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionAtHost(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionAtHost(input string) (*EncryptionAtHost, error) {
	vals := map[string]EncryptionAtHost{
		"disabled": EncryptionAtHostDisabled,
		"enabled":  EncryptionAtHostEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionAtHost(input)
	return &out, nil
}

type FipsValidatedModules string

const (
	FipsValidatedModulesDisabled FipsValidatedModules = "Disabled"
	FipsValidatedModulesEnabled  FipsValidatedModules = "Enabled"
)

func PossibleValuesForFipsValidatedModules() []string {
	return []string{
		string(FipsValidatedModulesDisabled),
		string(FipsValidatedModulesEnabled),
	}
}

func (s *FipsValidatedModules) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFipsValidatedModules(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFipsValidatedModules(input string) (*FipsValidatedModules, error) {
	vals := map[string]FipsValidatedModules{
		"disabled": FipsValidatedModulesDisabled,
		"enabled":  FipsValidatedModulesEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FipsValidatedModules(input)
	return &out, nil
}

type OutboundType string

const (
	OutboundTypeLoadbalancer       OutboundType = "Loadbalancer"
	OutboundTypeUserDefinedRouting OutboundType = "UserDefinedRouting"
)

func PossibleValuesForOutboundType() []string {
	return []string{
		string(OutboundTypeLoadbalancer),
		string(OutboundTypeUserDefinedRouting),
	}
}

func (s *OutboundType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOutboundType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOutboundType(input string) (*OutboundType, error) {
	vals := map[string]OutboundType{
		"loadbalancer":       OutboundTypeLoadbalancer,
		"userdefinedrouting": OutboundTypeUserDefinedRouting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OutboundType(input)
	return &out, nil
}

type PreconfiguredNSG string

const (
	PreconfiguredNSGDisabled PreconfiguredNSG = "Disabled"
	PreconfiguredNSGEnabled  PreconfiguredNSG = "Enabled"
)

func PossibleValuesForPreconfiguredNSG() []string {
	return []string{
		string(PreconfiguredNSGDisabled),
		string(PreconfiguredNSGEnabled),
	}
}

func (s *PreconfiguredNSG) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePreconfiguredNSG(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePreconfiguredNSG(input string) (*PreconfiguredNSG, error) {
	vals := map[string]PreconfiguredNSG{
		"disabled": PreconfiguredNSGDisabled,
		"enabled":  PreconfiguredNSGEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreconfiguredNSG(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAdminUpdating ProvisioningState = "AdminUpdating"
	ProvisioningStateCreating      ProvisioningState = "Creating"
	ProvisioningStateDeleting      ProvisioningState = "Deleting"
	ProvisioningStateFailed        ProvisioningState = "Failed"
	ProvisioningStateSucceeded     ProvisioningState = "Succeeded"
	ProvisioningStateUpdating      ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAdminUpdating),
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
		"adminupdating": ProvisioningStateAdminUpdating,
		"creating":      ProvisioningStateCreating,
		"deleting":      ProvisioningStateDeleting,
		"failed":        ProvisioningStateFailed,
		"succeeded":     ProvisioningStateSucceeded,
		"updating":      ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type Visibility string

const (
	VisibilityPrivate Visibility = "Private"
	VisibilityPublic  Visibility = "Public"
)

func PossibleValuesForVisibility() []string {
	return []string{
		string(VisibilityPrivate),
		string(VisibilityPublic),
	}
}

func (s *Visibility) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVisibility(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVisibility(input string) (*Visibility, error) {
	vals := map[string]Visibility{
		"private": VisibilityPrivate,
		"public":  VisibilityPublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Visibility(input)
	return &out, nil
}
