package authenticationpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OnUnauthenticatedRequest string

const (
	OnUnauthenticatedRequestAllow        OnUnauthenticatedRequest = "allow"
	OnUnauthenticatedRequestAuthenticate OnUnauthenticatedRequest = "authenticate"
	OnUnauthenticatedRequestDeny         OnUnauthenticatedRequest = "deny"
)

func PossibleValuesForOnUnauthenticatedRequest() []string {
	return []string{
		string(OnUnauthenticatedRequestAllow),
		string(OnUnauthenticatedRequestAuthenticate),
		string(OnUnauthenticatedRequestDeny),
	}
}

func (s *OnUnauthenticatedRequest) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOnUnauthenticatedRequest(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOnUnauthenticatedRequest(input string) (*OnUnauthenticatedRequest, error) {
	vals := map[string]OnUnauthenticatedRequest{
		"allow":        OnUnauthenticatedRequestAllow,
		"authenticate": OnUnauthenticatedRequestAuthenticate,
		"deny":         OnUnauthenticatedRequestDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OnUnauthenticatedRequest(input)
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

type UserTrustProviderType string

const (
	UserTrustProviderTypeEntra UserTrustProviderType = "entra"
)

func PossibleValuesForUserTrustProviderType() []string {
	return []string{
		string(UserTrustProviderTypeEntra),
	}
}

func (s *UserTrustProviderType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUserTrustProviderType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUserTrustProviderType(input string) (*UserTrustProviderType, error) {
	vals := map[string]UserTrustProviderType{
		"entra": UserTrustProviderTypeEntra,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UserTrustProviderType(input)
	return &out, nil
}
