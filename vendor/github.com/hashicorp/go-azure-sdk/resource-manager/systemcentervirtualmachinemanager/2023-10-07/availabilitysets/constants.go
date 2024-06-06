package availabilitysets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForceDelete string

const (
	ForceDeleteFalse ForceDelete = "false"
	ForceDeleteTrue  ForceDelete = "true"
)

func PossibleValuesForForceDelete() []string {
	return []string{
		string(ForceDeleteFalse),
		string(ForceDeleteTrue),
	}
}

func (s *ForceDelete) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseForceDelete(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseForceDelete(input string) (*ForceDelete, error) {
	vals := map[string]ForceDelete{
		"false": ForceDeleteFalse,
		"true":  ForceDeleteTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ForceDelete(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateAccepted     ResourceProvisioningState = "Accepted"
	ResourceProvisioningStateCanceled     ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateCreated      ResourceProvisioningState = "Created"
	ResourceProvisioningStateDeleting     ResourceProvisioningState = "Deleting"
	ResourceProvisioningStateFailed       ResourceProvisioningState = "Failed"
	ResourceProvisioningStateProvisioning ResourceProvisioningState = "Provisioning"
	ResourceProvisioningStateSucceeded    ResourceProvisioningState = "Succeeded"
	ResourceProvisioningStateUpdating     ResourceProvisioningState = "Updating"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateAccepted),
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateCreated),
		string(ResourceProvisioningStateDeleting),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateProvisioning),
		string(ResourceProvisioningStateSucceeded),
		string(ResourceProvisioningStateUpdating),
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
		"accepted":     ResourceProvisioningStateAccepted,
		"canceled":     ResourceProvisioningStateCanceled,
		"created":      ResourceProvisioningStateCreated,
		"deleting":     ResourceProvisioningStateDeleting,
		"failed":       ResourceProvisioningStateFailed,
		"provisioning": ResourceProvisioningStateProvisioning,
		"succeeded":    ResourceProvisioningStateSucceeded,
		"updating":     ResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}
