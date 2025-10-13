package projectpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevCenterResourceType string

const (
	DevCenterResourceTypeAttachedNetworks DevCenterResourceType = "AttachedNetworks"
	DevCenterResourceTypeImages           DevCenterResourceType = "Images"
	DevCenterResourceTypeSkus             DevCenterResourceType = "Skus"
)

func PossibleValuesForDevCenterResourceType() []string {
	return []string{
		string(DevCenterResourceTypeAttachedNetworks),
		string(DevCenterResourceTypeImages),
		string(DevCenterResourceTypeSkus),
	}
}

func (s *DevCenterResourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDevCenterResourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDevCenterResourceType(input string) (*DevCenterResourceType, error) {
	vals := map[string]DevCenterResourceType{
		"attachednetworks": DevCenterResourceTypeAttachedNetworks,
		"images":           DevCenterResourceTypeImages,
		"skus":             DevCenterResourceTypeSkus,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DevCenterResourceType(input)
	return &out, nil
}

type PolicyAction string

const (
	PolicyActionAllow PolicyAction = "Allow"
	PolicyActionDeny  PolicyAction = "Deny"
)

func PossibleValuesForPolicyAction() []string {
	return []string{
		string(PolicyActionAllow),
		string(PolicyActionDeny),
	}
}

func (s *PolicyAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyAction(input string) (*PolicyAction, error) {
	vals := map[string]PolicyAction{
		"allow": PolicyActionAllow,
		"deny":  PolicyActionDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyAction(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted                  ProvisioningState = "Accepted"
	ProvisioningStateCanceled                  ProvisioningState = "Canceled"
	ProvisioningStateCreated                   ProvisioningState = "Created"
	ProvisioningStateCreating                  ProvisioningState = "Creating"
	ProvisioningStateDeleted                   ProvisioningState = "Deleted"
	ProvisioningStateDeleting                  ProvisioningState = "Deleting"
	ProvisioningStateFailed                    ProvisioningState = "Failed"
	ProvisioningStateMovingResources           ProvisioningState = "MovingResources"
	ProvisioningStateNotSpecified              ProvisioningState = "NotSpecified"
	ProvisioningStateRolloutInProgress         ProvisioningState = "RolloutInProgress"
	ProvisioningStateRunning                   ProvisioningState = "Running"
	ProvisioningStateStorageProvisioningFailed ProvisioningState = "StorageProvisioningFailed"
	ProvisioningStateSucceeded                 ProvisioningState = "Succeeded"
	ProvisioningStateTransientFailure          ProvisioningState = "TransientFailure"
	ProvisioningStateUpdated                   ProvisioningState = "Updated"
	ProvisioningStateUpdating                  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreated),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMovingResources),
		string(ProvisioningStateNotSpecified),
		string(ProvisioningStateRolloutInProgress),
		string(ProvisioningStateRunning),
		string(ProvisioningStateStorageProvisioningFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateTransientFailure),
		string(ProvisioningStateUpdated),
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
		"accepted":                  ProvisioningStateAccepted,
		"canceled":                  ProvisioningStateCanceled,
		"created":                   ProvisioningStateCreated,
		"creating":                  ProvisioningStateCreating,
		"deleted":                   ProvisioningStateDeleted,
		"deleting":                  ProvisioningStateDeleting,
		"failed":                    ProvisioningStateFailed,
		"movingresources":           ProvisioningStateMovingResources,
		"notspecified":              ProvisioningStateNotSpecified,
		"rolloutinprogress":         ProvisioningStateRolloutInProgress,
		"running":                   ProvisioningStateRunning,
		"storageprovisioningfailed": ProvisioningStateStorageProvisioningFailed,
		"succeeded":                 ProvisioningStateSucceeded,
		"transientfailure":          ProvisioningStateTransientFailure,
		"updated":                   ProvisioningStateUpdated,
		"updating":                  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
