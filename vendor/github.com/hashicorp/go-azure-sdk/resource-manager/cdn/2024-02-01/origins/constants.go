package origins

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OriginProvisioningState string

const (
	OriginProvisioningStateCreating  OriginProvisioningState = "Creating"
	OriginProvisioningStateDeleting  OriginProvisioningState = "Deleting"
	OriginProvisioningStateFailed    OriginProvisioningState = "Failed"
	OriginProvisioningStateSucceeded OriginProvisioningState = "Succeeded"
	OriginProvisioningStateUpdating  OriginProvisioningState = "Updating"
)

func PossibleValuesForOriginProvisioningState() []string {
	return []string{
		string(OriginProvisioningStateCreating),
		string(OriginProvisioningStateDeleting),
		string(OriginProvisioningStateFailed),
		string(OriginProvisioningStateSucceeded),
		string(OriginProvisioningStateUpdating),
	}
}

func (s *OriginProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOriginProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOriginProvisioningState(input string) (*OriginProvisioningState, error) {
	vals := map[string]OriginProvisioningState{
		"creating":  OriginProvisioningStateCreating,
		"deleting":  OriginProvisioningStateDeleting,
		"failed":    OriginProvisioningStateFailed,
		"succeeded": OriginProvisioningStateSucceeded,
		"updating":  OriginProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OriginProvisioningState(input)
	return &out, nil
}

type OriginResourceState string

const (
	OriginResourceStateActive   OriginResourceState = "Active"
	OriginResourceStateCreating OriginResourceState = "Creating"
	OriginResourceStateDeleting OriginResourceState = "Deleting"
)

func PossibleValuesForOriginResourceState() []string {
	return []string{
		string(OriginResourceStateActive),
		string(OriginResourceStateCreating),
		string(OriginResourceStateDeleting),
	}
}

func (s *OriginResourceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOriginResourceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOriginResourceState(input string) (*OriginResourceState, error) {
	vals := map[string]OriginResourceState{
		"active":   OriginResourceStateActive,
		"creating": OriginResourceStateCreating,
		"deleting": OriginResourceStateDeleting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OriginResourceState(input)
	return &out, nil
}

type PrivateEndpointStatus string

const (
	PrivateEndpointStatusApproved     PrivateEndpointStatus = "Approved"
	PrivateEndpointStatusDisconnected PrivateEndpointStatus = "Disconnected"
	PrivateEndpointStatusPending      PrivateEndpointStatus = "Pending"
	PrivateEndpointStatusRejected     PrivateEndpointStatus = "Rejected"
	PrivateEndpointStatusTimeout      PrivateEndpointStatus = "Timeout"
)

func PossibleValuesForPrivateEndpointStatus() []string {
	return []string{
		string(PrivateEndpointStatusApproved),
		string(PrivateEndpointStatusDisconnected),
		string(PrivateEndpointStatusPending),
		string(PrivateEndpointStatusRejected),
		string(PrivateEndpointStatusTimeout),
	}
}

func (s *PrivateEndpointStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointStatus(input string) (*PrivateEndpointStatus, error) {
	vals := map[string]PrivateEndpointStatus{
		"approved":     PrivateEndpointStatusApproved,
		"disconnected": PrivateEndpointStatusDisconnected,
		"pending":      PrivateEndpointStatusPending,
		"rejected":     PrivateEndpointStatusRejected,
		"timeout":      PrivateEndpointStatusTimeout,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointStatus(input)
	return &out, nil
}
