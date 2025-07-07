package dnsprivateviews

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsPrivateViewsLifecycleState string

const (
	DnsPrivateViewsLifecycleStateActive   DnsPrivateViewsLifecycleState = "Active"
	DnsPrivateViewsLifecycleStateDeleted  DnsPrivateViewsLifecycleState = "Deleted"
	DnsPrivateViewsLifecycleStateDeleting DnsPrivateViewsLifecycleState = "Deleting"
	DnsPrivateViewsLifecycleStateUpdating DnsPrivateViewsLifecycleState = "Updating"
)

func PossibleValuesForDnsPrivateViewsLifecycleState() []string {
	return []string{
		string(DnsPrivateViewsLifecycleStateActive),
		string(DnsPrivateViewsLifecycleStateDeleted),
		string(DnsPrivateViewsLifecycleStateDeleting),
		string(DnsPrivateViewsLifecycleStateUpdating),
	}
}

func (s *DnsPrivateViewsLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDnsPrivateViewsLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDnsPrivateViewsLifecycleState(input string) (*DnsPrivateViewsLifecycleState, error) {
	vals := map[string]DnsPrivateViewsLifecycleState{
		"active":   DnsPrivateViewsLifecycleStateActive,
		"deleted":  DnsPrivateViewsLifecycleStateDeleted,
		"deleting": DnsPrivateViewsLifecycleStateDeleting,
		"updating": DnsPrivateViewsLifecycleStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DnsPrivateViewsLifecycleState(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateSucceeded),
	}
}

func (s *ResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceProvisioningState(input string) (*ResourceProvisioningState, error) {
	vals := map[string]ResourceProvisioningState{
		"canceled":  ResourceProvisioningStateCanceled,
		"failed":    ResourceProvisioningStateFailed,
		"succeeded": ResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}
