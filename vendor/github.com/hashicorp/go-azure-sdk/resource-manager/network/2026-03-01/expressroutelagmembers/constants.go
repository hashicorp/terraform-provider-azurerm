package expressroutelagmembers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkAdminState string

const (
	ExpressRouteLinkAdminStateDisabled ExpressRouteLinkAdminState = "Disabled"
	ExpressRouteLinkAdminStateEnabled  ExpressRouteLinkAdminState = "Enabled"
)

func PossibleValuesForExpressRouteLinkAdminState() []string {
	return []string{
		string(ExpressRouteLinkAdminStateDisabled),
		string(ExpressRouteLinkAdminStateEnabled),
	}
}

func (s *ExpressRouteLinkAdminState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRouteLinkAdminState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRouteLinkAdminState(input string) (*ExpressRouteLinkAdminState, error) {
	vals := map[string]ExpressRouteLinkAdminState{
		"disabled": ExpressRouteLinkAdminStateDisabled,
		"enabled":  ExpressRouteLinkAdminStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRouteLinkAdminState(input)
	return &out, nil
}

type ExpressRouteLinkConnectorType string

const (
	ExpressRouteLinkConnectorTypeLC ExpressRouteLinkConnectorType = "LC"
	ExpressRouteLinkConnectorTypeSC ExpressRouteLinkConnectorType = "SC"
)

func PossibleValuesForExpressRouteLinkConnectorType() []string {
	return []string{
		string(ExpressRouteLinkConnectorTypeLC),
		string(ExpressRouteLinkConnectorTypeSC),
	}
}

func (s *ExpressRouteLinkConnectorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpressRouteLinkConnectorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpressRouteLinkConnectorType(input string) (*ExpressRouteLinkConnectorType, error) {
	vals := map[string]ExpressRouteLinkConnectorType{
		"lc": ExpressRouteLinkConnectorTypeLC,
		"sc": ExpressRouteLinkConnectorTypeSC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpressRouteLinkConnectorType(input)
	return &out, nil
}

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
