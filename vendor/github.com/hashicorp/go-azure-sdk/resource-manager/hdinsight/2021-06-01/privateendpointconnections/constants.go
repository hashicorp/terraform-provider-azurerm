package privateendpointconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCanceled   PrivateEndpointConnectionProvisioningState = "Canceled"
	PrivateEndpointConnectionProvisioningStateDeleting   PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed     PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateInProgress PrivateEndpointConnectionProvisioningState = "InProgress"
	PrivateEndpointConnectionProvisioningStateSucceeded  PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUpdating   PrivateEndpointConnectionProvisioningState = "Updating"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCanceled),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateInProgress),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
		string(PrivateEndpointConnectionProvisioningStateUpdating),
	}
}

func (s *PrivateEndpointConnectionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointConnectionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"canceled":   PrivateEndpointConnectionProvisioningStateCanceled,
		"deleting":   PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":     PrivateEndpointConnectionProvisioningStateFailed,
		"inprogress": PrivateEndpointConnectionProvisioningStateInProgress,
		"succeeded":  PrivateEndpointConnectionProvisioningStateSucceeded,
		"updating":   PrivateEndpointConnectionProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateLinkServiceConnectionStatus string

const (
	PrivateLinkServiceConnectionStatusApproved PrivateLinkServiceConnectionStatus = "Approved"
	PrivateLinkServiceConnectionStatusPending  PrivateLinkServiceConnectionStatus = "Pending"
	PrivateLinkServiceConnectionStatusRejected PrivateLinkServiceConnectionStatus = "Rejected"
	PrivateLinkServiceConnectionStatusRemoved  PrivateLinkServiceConnectionStatus = "Removed"
)

func PossibleValuesForPrivateLinkServiceConnectionStatus() []string {
	return []string{
		string(PrivateLinkServiceConnectionStatusApproved),
		string(PrivateLinkServiceConnectionStatusPending),
		string(PrivateLinkServiceConnectionStatusRejected),
		string(PrivateLinkServiceConnectionStatusRemoved),
	}
}

func (s *PrivateLinkServiceConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateLinkServiceConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateLinkServiceConnectionStatus(input string) (*PrivateLinkServiceConnectionStatus, error) {
	vals := map[string]PrivateLinkServiceConnectionStatus{
		"approved": PrivateLinkServiceConnectionStatusApproved,
		"pending":  PrivateLinkServiceConnectionStatusPending,
		"rejected": PrivateLinkServiceConnectionStatusRejected,
		"removed":  PrivateLinkServiceConnectionStatusRemoved,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateLinkServiceConnectionStatus(input)
	return &out, nil
}
