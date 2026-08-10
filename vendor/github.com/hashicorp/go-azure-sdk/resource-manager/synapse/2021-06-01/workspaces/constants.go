package workspaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActualState string

const (
	ActualStateDisabled  ActualState = "Disabled"
	ActualStateDisabling ActualState = "Disabling"
	ActualStateEnabled   ActualState = "Enabled"
	ActualStateEnabling  ActualState = "Enabling"
	ActualStateUnknown   ActualState = "Unknown"
)

func PossibleValuesForActualState() []string {
	return []string{
		string(ActualStateDisabled),
		string(ActualStateDisabling),
		string(ActualStateEnabled),
		string(ActualStateEnabling),
		string(ActualStateUnknown),
	}
}

func (s *ActualState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActualState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActualState(input string) (*ActualState, error) {
	vals := map[string]ActualState{
		"disabled":  ActualStateDisabled,
		"disabling": ActualStateDisabling,
		"enabled":   ActualStateEnabled,
		"enabling":  ActualStateEnabling,
		"unknown":   ActualStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActualState(input)
	return &out, nil
}

type DesiredState string

const (
	DesiredStateDisabled DesiredState = "Disabled"
	DesiredStateEnabled  DesiredState = "Enabled"
)

func PossibleValuesForDesiredState() []string {
	return []string{
		string(DesiredStateDisabled),
		string(DesiredStateEnabled),
	}
}

func (s *DesiredState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDesiredState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDesiredState(input string) (*DesiredState, error) {
	vals := map[string]DesiredState{
		"disabled": DesiredStateDisabled,
		"enabled":  DesiredStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DesiredState(input)
	return &out, nil
}

type ResourceIdentityType string

const (
	ResourceIdentityTypeNone                       ResourceIdentityType = "None"
	ResourceIdentityTypeSystemAssigned             ResourceIdentityType = "SystemAssigned"
	ResourceIdentityTypeSystemAssignedUserAssigned ResourceIdentityType = "SystemAssigned,UserAssigned"
)

func PossibleValuesForResourceIdentityType() []string {
	return []string{
		string(ResourceIdentityTypeNone),
		string(ResourceIdentityTypeSystemAssigned),
		string(ResourceIdentityTypeSystemAssignedUserAssigned),
	}
}

func (s *ResourceIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceIdentityType(input string) (*ResourceIdentityType, error) {
	vals := map[string]ResourceIdentityType{
		"none":                        ResourceIdentityTypeNone,
		"systemassigned":              ResourceIdentityTypeSystemAssigned,
		"systemassigned,userassigned": ResourceIdentityTypeSystemAssignedUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceIdentityType(input)
	return &out, nil
}

type WorkspacePublicNetworkAccess string

const (
	WorkspacePublicNetworkAccessDisabled WorkspacePublicNetworkAccess = "Disabled"
	WorkspacePublicNetworkAccessEnabled  WorkspacePublicNetworkAccess = "Enabled"
)

func PossibleValuesForWorkspacePublicNetworkAccess() []string {
	return []string{
		string(WorkspacePublicNetworkAccessDisabled),
		string(WorkspacePublicNetworkAccessEnabled),
	}
}

func (s *WorkspacePublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkspacePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkspacePublicNetworkAccess(input string) (*WorkspacePublicNetworkAccess, error) {
	vals := map[string]WorkspacePublicNetworkAccess{
		"disabled": WorkspacePublicNetworkAccessDisabled,
		"enabled":  WorkspacePublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkspacePublicNetworkAccess(input)
	return &out, nil
}
