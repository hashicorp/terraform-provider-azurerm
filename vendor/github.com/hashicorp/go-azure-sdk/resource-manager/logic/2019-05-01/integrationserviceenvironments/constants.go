package integrationserviceenvironments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationServiceEnvironmentAccessEndpointType string

const (
	IntegrationServiceEnvironmentAccessEndpointTypeExternal     IntegrationServiceEnvironmentAccessEndpointType = "External"
	IntegrationServiceEnvironmentAccessEndpointTypeInternal     IntegrationServiceEnvironmentAccessEndpointType = "Internal"
	IntegrationServiceEnvironmentAccessEndpointTypeNotSpecified IntegrationServiceEnvironmentAccessEndpointType = "NotSpecified"
)

func PossibleValuesForIntegrationServiceEnvironmentAccessEndpointType() []string {
	return []string{
		string(IntegrationServiceEnvironmentAccessEndpointTypeExternal),
		string(IntegrationServiceEnvironmentAccessEndpointTypeInternal),
		string(IntegrationServiceEnvironmentAccessEndpointTypeNotSpecified),
	}
}

func (s *IntegrationServiceEnvironmentAccessEndpointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationServiceEnvironmentAccessEndpointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationServiceEnvironmentAccessEndpointType(input string) (*IntegrationServiceEnvironmentAccessEndpointType, error) {
	vals := map[string]IntegrationServiceEnvironmentAccessEndpointType{
		"external":     IntegrationServiceEnvironmentAccessEndpointTypeExternal,
		"internal":     IntegrationServiceEnvironmentAccessEndpointTypeInternal,
		"notspecified": IntegrationServiceEnvironmentAccessEndpointTypeNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationServiceEnvironmentAccessEndpointType(input)
	return &out, nil
}

type IntegrationServiceEnvironmentSkuName string

const (
	IntegrationServiceEnvironmentSkuNameDeveloper    IntegrationServiceEnvironmentSkuName = "Developer"
	IntegrationServiceEnvironmentSkuNameNotSpecified IntegrationServiceEnvironmentSkuName = "NotSpecified"
	IntegrationServiceEnvironmentSkuNamePremium      IntegrationServiceEnvironmentSkuName = "Premium"
)

func PossibleValuesForIntegrationServiceEnvironmentSkuName() []string {
	return []string{
		string(IntegrationServiceEnvironmentSkuNameDeveloper),
		string(IntegrationServiceEnvironmentSkuNameNotSpecified),
		string(IntegrationServiceEnvironmentSkuNamePremium),
	}
}

func (s *IntegrationServiceEnvironmentSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationServiceEnvironmentSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationServiceEnvironmentSkuName(input string) (*IntegrationServiceEnvironmentSkuName, error) {
	vals := map[string]IntegrationServiceEnvironmentSkuName{
		"developer":    IntegrationServiceEnvironmentSkuNameDeveloper,
		"notspecified": IntegrationServiceEnvironmentSkuNameNotSpecified,
		"premium":      IntegrationServiceEnvironmentSkuNamePremium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationServiceEnvironmentSkuName(input)
	return &out, nil
}

type WorkflowProvisioningState string

const (
	WorkflowProvisioningStateAccepted      WorkflowProvisioningState = "Accepted"
	WorkflowProvisioningStateCanceled      WorkflowProvisioningState = "Canceled"
	WorkflowProvisioningStateCompleted     WorkflowProvisioningState = "Completed"
	WorkflowProvisioningStateCreated       WorkflowProvisioningState = "Created"
	WorkflowProvisioningStateCreating      WorkflowProvisioningState = "Creating"
	WorkflowProvisioningStateDeleted       WorkflowProvisioningState = "Deleted"
	WorkflowProvisioningStateDeleting      WorkflowProvisioningState = "Deleting"
	WorkflowProvisioningStateFailed        WorkflowProvisioningState = "Failed"
	WorkflowProvisioningStateInProgress    WorkflowProvisioningState = "InProgress"
	WorkflowProvisioningStateMoving        WorkflowProvisioningState = "Moving"
	WorkflowProvisioningStateNotSpecified  WorkflowProvisioningState = "NotSpecified"
	WorkflowProvisioningStatePending       WorkflowProvisioningState = "Pending"
	WorkflowProvisioningStateReady         WorkflowProvisioningState = "Ready"
	WorkflowProvisioningStateRegistered    WorkflowProvisioningState = "Registered"
	WorkflowProvisioningStateRegistering   WorkflowProvisioningState = "Registering"
	WorkflowProvisioningStateRenewing      WorkflowProvisioningState = "Renewing"
	WorkflowProvisioningStateRunning       WorkflowProvisioningState = "Running"
	WorkflowProvisioningStateSucceeded     WorkflowProvisioningState = "Succeeded"
	WorkflowProvisioningStateUnregistered  WorkflowProvisioningState = "Unregistered"
	WorkflowProvisioningStateUnregistering WorkflowProvisioningState = "Unregistering"
	WorkflowProvisioningStateUpdating      WorkflowProvisioningState = "Updating"
	WorkflowProvisioningStateWaiting       WorkflowProvisioningState = "Waiting"
)

func PossibleValuesForWorkflowProvisioningState() []string {
	return []string{
		string(WorkflowProvisioningStateAccepted),
		string(WorkflowProvisioningStateCanceled),
		string(WorkflowProvisioningStateCompleted),
		string(WorkflowProvisioningStateCreated),
		string(WorkflowProvisioningStateCreating),
		string(WorkflowProvisioningStateDeleted),
		string(WorkflowProvisioningStateDeleting),
		string(WorkflowProvisioningStateFailed),
		string(WorkflowProvisioningStateInProgress),
		string(WorkflowProvisioningStateMoving),
		string(WorkflowProvisioningStateNotSpecified),
		string(WorkflowProvisioningStatePending),
		string(WorkflowProvisioningStateReady),
		string(WorkflowProvisioningStateRegistered),
		string(WorkflowProvisioningStateRegistering),
		string(WorkflowProvisioningStateRenewing),
		string(WorkflowProvisioningStateRunning),
		string(WorkflowProvisioningStateSucceeded),
		string(WorkflowProvisioningStateUnregistered),
		string(WorkflowProvisioningStateUnregistering),
		string(WorkflowProvisioningStateUpdating),
		string(WorkflowProvisioningStateWaiting),
	}
}

func (s *WorkflowProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkflowProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkflowProvisioningState(input string) (*WorkflowProvisioningState, error) {
	vals := map[string]WorkflowProvisioningState{
		"accepted":      WorkflowProvisioningStateAccepted,
		"canceled":      WorkflowProvisioningStateCanceled,
		"completed":     WorkflowProvisioningStateCompleted,
		"created":       WorkflowProvisioningStateCreated,
		"creating":      WorkflowProvisioningStateCreating,
		"deleted":       WorkflowProvisioningStateDeleted,
		"deleting":      WorkflowProvisioningStateDeleting,
		"failed":        WorkflowProvisioningStateFailed,
		"inprogress":    WorkflowProvisioningStateInProgress,
		"moving":        WorkflowProvisioningStateMoving,
		"notspecified":  WorkflowProvisioningStateNotSpecified,
		"pending":       WorkflowProvisioningStatePending,
		"ready":         WorkflowProvisioningStateReady,
		"registered":    WorkflowProvisioningStateRegistered,
		"registering":   WorkflowProvisioningStateRegistering,
		"renewing":      WorkflowProvisioningStateRenewing,
		"running":       WorkflowProvisioningStateRunning,
		"succeeded":     WorkflowProvisioningStateSucceeded,
		"unregistered":  WorkflowProvisioningStateUnregistered,
		"unregistering": WorkflowProvisioningStateUnregistering,
		"updating":      WorkflowProvisioningStateUpdating,
		"waiting":       WorkflowProvisioningStateWaiting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkflowProvisioningState(input)
	return &out, nil
}

type WorkflowState string

const (
	WorkflowStateCompleted    WorkflowState = "Completed"
	WorkflowStateDeleted      WorkflowState = "Deleted"
	WorkflowStateDisabled     WorkflowState = "Disabled"
	WorkflowStateEnabled      WorkflowState = "Enabled"
	WorkflowStateNotSpecified WorkflowState = "NotSpecified"
	WorkflowStateSuspended    WorkflowState = "Suspended"
)

func PossibleValuesForWorkflowState() []string {
	return []string{
		string(WorkflowStateCompleted),
		string(WorkflowStateDeleted),
		string(WorkflowStateDisabled),
		string(WorkflowStateEnabled),
		string(WorkflowStateNotSpecified),
		string(WorkflowStateSuspended),
	}
}

func (s *WorkflowState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkflowState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkflowState(input string) (*WorkflowState, error) {
	vals := map[string]WorkflowState{
		"completed":    WorkflowStateCompleted,
		"deleted":      WorkflowStateDeleted,
		"disabled":     WorkflowStateDisabled,
		"enabled":      WorkflowStateEnabled,
		"notspecified": WorkflowStateNotSpecified,
		"suspended":    WorkflowStateSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkflowState(input)
	return &out, nil
}
