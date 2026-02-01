package sharedprivatelinkresources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedPrivateLinkResourceProvisioningState string

const (
	SharedPrivateLinkResourceProvisioningStateDeleting   SharedPrivateLinkResourceProvisioningState = "Deleting"
	SharedPrivateLinkResourceProvisioningStateFailed     SharedPrivateLinkResourceProvisioningState = "Failed"
	SharedPrivateLinkResourceProvisioningStateIncomplete SharedPrivateLinkResourceProvisioningState = "Incomplete"
	SharedPrivateLinkResourceProvisioningStateSucceeded  SharedPrivateLinkResourceProvisioningState = "Succeeded"
	SharedPrivateLinkResourceProvisioningStateUpdating   SharedPrivateLinkResourceProvisioningState = "Updating"
)

func PossibleValuesForSharedPrivateLinkResourceProvisioningState() []string {
	return []string{
		string(SharedPrivateLinkResourceProvisioningStateDeleting),
		string(SharedPrivateLinkResourceProvisioningStateFailed),
		string(SharedPrivateLinkResourceProvisioningStateIncomplete),
		string(SharedPrivateLinkResourceProvisioningStateSucceeded),
		string(SharedPrivateLinkResourceProvisioningStateUpdating),
	}
}

func (s *SharedPrivateLinkResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSharedPrivateLinkResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSharedPrivateLinkResourceProvisioningState(input string) (*SharedPrivateLinkResourceProvisioningState, error) {
	vals := map[string]SharedPrivateLinkResourceProvisioningState{
		"deleting":   SharedPrivateLinkResourceProvisioningStateDeleting,
		"failed":     SharedPrivateLinkResourceProvisioningStateFailed,
		"incomplete": SharedPrivateLinkResourceProvisioningStateIncomplete,
		"succeeded":  SharedPrivateLinkResourceProvisioningStateSucceeded,
		"updating":   SharedPrivateLinkResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharedPrivateLinkResourceProvisioningState(input)
	return &out, nil
}

type SharedPrivateLinkResourceStatus string

const (
	SharedPrivateLinkResourceStatusApproved     SharedPrivateLinkResourceStatus = "Approved"
	SharedPrivateLinkResourceStatusDisconnected SharedPrivateLinkResourceStatus = "Disconnected"
	SharedPrivateLinkResourceStatusPending      SharedPrivateLinkResourceStatus = "Pending"
	SharedPrivateLinkResourceStatusRejected     SharedPrivateLinkResourceStatus = "Rejected"
)

func PossibleValuesForSharedPrivateLinkResourceStatus() []string {
	return []string{
		string(SharedPrivateLinkResourceStatusApproved),
		string(SharedPrivateLinkResourceStatusDisconnected),
		string(SharedPrivateLinkResourceStatusPending),
		string(SharedPrivateLinkResourceStatusRejected),
	}
}

func (s *SharedPrivateLinkResourceStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSharedPrivateLinkResourceStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSharedPrivateLinkResourceStatus(input string) (*SharedPrivateLinkResourceStatus, error) {
	vals := map[string]SharedPrivateLinkResourceStatus{
		"approved":     SharedPrivateLinkResourceStatusApproved,
		"disconnected": SharedPrivateLinkResourceStatusDisconnected,
		"pending":      SharedPrivateLinkResourceStatusPending,
		"rejected":     SharedPrivateLinkResourceStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharedPrivateLinkResourceStatus(input)
	return &out, nil
}
