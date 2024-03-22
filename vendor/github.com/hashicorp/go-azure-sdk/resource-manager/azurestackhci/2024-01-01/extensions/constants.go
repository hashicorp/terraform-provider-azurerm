package extensions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionAggregateState string

const (
	ExtensionAggregateStateAccepted                       ExtensionAggregateState = "Accepted"
	ExtensionAggregateStateCanceled                       ExtensionAggregateState = "Canceled"
	ExtensionAggregateStateConnected                      ExtensionAggregateState = "Connected"
	ExtensionAggregateStateCreating                       ExtensionAggregateState = "Creating"
	ExtensionAggregateStateDeleted                        ExtensionAggregateState = "Deleted"
	ExtensionAggregateStateDeleting                       ExtensionAggregateState = "Deleting"
	ExtensionAggregateStateDisconnected                   ExtensionAggregateState = "Disconnected"
	ExtensionAggregateStateError                          ExtensionAggregateState = "Error"
	ExtensionAggregateStateFailed                         ExtensionAggregateState = "Failed"
	ExtensionAggregateStateInProgress                     ExtensionAggregateState = "InProgress"
	ExtensionAggregateStateMoving                         ExtensionAggregateState = "Moving"
	ExtensionAggregateStateNotSpecified                   ExtensionAggregateState = "NotSpecified"
	ExtensionAggregateStatePartiallyConnected             ExtensionAggregateState = "PartiallyConnected"
	ExtensionAggregateStatePartiallySucceeded             ExtensionAggregateState = "PartiallySucceeded"
	ExtensionAggregateStateProvisioning                   ExtensionAggregateState = "Provisioning"
	ExtensionAggregateStateSucceeded                      ExtensionAggregateState = "Succeeded"
	ExtensionAggregateStateUpdating                       ExtensionAggregateState = "Updating"
	ExtensionAggregateStateUpgradeFailedRollbackSucceeded ExtensionAggregateState = "UpgradeFailedRollbackSucceeded"
)

func PossibleValuesForExtensionAggregateState() []string {
	return []string{
		string(ExtensionAggregateStateAccepted),
		string(ExtensionAggregateStateCanceled),
		string(ExtensionAggregateStateConnected),
		string(ExtensionAggregateStateCreating),
		string(ExtensionAggregateStateDeleted),
		string(ExtensionAggregateStateDeleting),
		string(ExtensionAggregateStateDisconnected),
		string(ExtensionAggregateStateError),
		string(ExtensionAggregateStateFailed),
		string(ExtensionAggregateStateInProgress),
		string(ExtensionAggregateStateMoving),
		string(ExtensionAggregateStateNotSpecified),
		string(ExtensionAggregateStatePartiallyConnected),
		string(ExtensionAggregateStatePartiallySucceeded),
		string(ExtensionAggregateStateProvisioning),
		string(ExtensionAggregateStateSucceeded),
		string(ExtensionAggregateStateUpdating),
		string(ExtensionAggregateStateUpgradeFailedRollbackSucceeded),
	}
}

func (s *ExtensionAggregateState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtensionAggregateState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtensionAggregateState(input string) (*ExtensionAggregateState, error) {
	vals := map[string]ExtensionAggregateState{
		"accepted":                       ExtensionAggregateStateAccepted,
		"canceled":                       ExtensionAggregateStateCanceled,
		"connected":                      ExtensionAggregateStateConnected,
		"creating":                       ExtensionAggregateStateCreating,
		"deleted":                        ExtensionAggregateStateDeleted,
		"deleting":                       ExtensionAggregateStateDeleting,
		"disconnected":                   ExtensionAggregateStateDisconnected,
		"error":                          ExtensionAggregateStateError,
		"failed":                         ExtensionAggregateStateFailed,
		"inprogress":                     ExtensionAggregateStateInProgress,
		"moving":                         ExtensionAggregateStateMoving,
		"notspecified":                   ExtensionAggregateStateNotSpecified,
		"partiallyconnected":             ExtensionAggregateStatePartiallyConnected,
		"partiallysucceeded":             ExtensionAggregateStatePartiallySucceeded,
		"provisioning":                   ExtensionAggregateStateProvisioning,
		"succeeded":                      ExtensionAggregateStateSucceeded,
		"updating":                       ExtensionAggregateStateUpdating,
		"upgradefailedrollbacksucceeded": ExtensionAggregateStateUpgradeFailedRollbackSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtensionAggregateState(input)
	return &out, nil
}

type ExtensionManagedBy string

const (
	ExtensionManagedByAzure ExtensionManagedBy = "Azure"
	ExtensionManagedByUser  ExtensionManagedBy = "User"
)

func PossibleValuesForExtensionManagedBy() []string {
	return []string{
		string(ExtensionManagedByAzure),
		string(ExtensionManagedByUser),
	}
}

func (s *ExtensionManagedBy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtensionManagedBy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtensionManagedBy(input string) (*ExtensionManagedBy, error) {
	vals := map[string]ExtensionManagedBy{
		"azure": ExtensionManagedByAzure,
		"user":  ExtensionManagedByUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtensionManagedBy(input)
	return &out, nil
}

type NodeExtensionState string

const (
	NodeExtensionStateAccepted           NodeExtensionState = "Accepted"
	NodeExtensionStateCanceled           NodeExtensionState = "Canceled"
	NodeExtensionStateConnected          NodeExtensionState = "Connected"
	NodeExtensionStateCreating           NodeExtensionState = "Creating"
	NodeExtensionStateDeleted            NodeExtensionState = "Deleted"
	NodeExtensionStateDeleting           NodeExtensionState = "Deleting"
	NodeExtensionStateDisconnected       NodeExtensionState = "Disconnected"
	NodeExtensionStateError              NodeExtensionState = "Error"
	NodeExtensionStateFailed             NodeExtensionState = "Failed"
	NodeExtensionStateInProgress         NodeExtensionState = "InProgress"
	NodeExtensionStateMoving             NodeExtensionState = "Moving"
	NodeExtensionStateNotSpecified       NodeExtensionState = "NotSpecified"
	NodeExtensionStatePartiallyConnected NodeExtensionState = "PartiallyConnected"
	NodeExtensionStatePartiallySucceeded NodeExtensionState = "PartiallySucceeded"
	NodeExtensionStateProvisioning       NodeExtensionState = "Provisioning"
	NodeExtensionStateSucceeded          NodeExtensionState = "Succeeded"
	NodeExtensionStateUpdating           NodeExtensionState = "Updating"
)

func PossibleValuesForNodeExtensionState() []string {
	return []string{
		string(NodeExtensionStateAccepted),
		string(NodeExtensionStateCanceled),
		string(NodeExtensionStateConnected),
		string(NodeExtensionStateCreating),
		string(NodeExtensionStateDeleted),
		string(NodeExtensionStateDeleting),
		string(NodeExtensionStateDisconnected),
		string(NodeExtensionStateError),
		string(NodeExtensionStateFailed),
		string(NodeExtensionStateInProgress),
		string(NodeExtensionStateMoving),
		string(NodeExtensionStateNotSpecified),
		string(NodeExtensionStatePartiallyConnected),
		string(NodeExtensionStatePartiallySucceeded),
		string(NodeExtensionStateProvisioning),
		string(NodeExtensionStateSucceeded),
		string(NodeExtensionStateUpdating),
	}
}

func (s *NodeExtensionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNodeExtensionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNodeExtensionState(input string) (*NodeExtensionState, error) {
	vals := map[string]NodeExtensionState{
		"accepted":           NodeExtensionStateAccepted,
		"canceled":           NodeExtensionStateCanceled,
		"connected":          NodeExtensionStateConnected,
		"creating":           NodeExtensionStateCreating,
		"deleted":            NodeExtensionStateDeleted,
		"deleting":           NodeExtensionStateDeleting,
		"disconnected":       NodeExtensionStateDisconnected,
		"error":              NodeExtensionStateError,
		"failed":             NodeExtensionStateFailed,
		"inprogress":         NodeExtensionStateInProgress,
		"moving":             NodeExtensionStateMoving,
		"notspecified":       NodeExtensionStateNotSpecified,
		"partiallyconnected": NodeExtensionStatePartiallyConnected,
		"partiallysucceeded": NodeExtensionStatePartiallySucceeded,
		"provisioning":       NodeExtensionStateProvisioning,
		"succeeded":          NodeExtensionStateSucceeded,
		"updating":           NodeExtensionStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodeExtensionState(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted           ProvisioningState = "Accepted"
	ProvisioningStateCanceled           ProvisioningState = "Canceled"
	ProvisioningStateConnected          ProvisioningState = "Connected"
	ProvisioningStateCreating           ProvisioningState = "Creating"
	ProvisioningStateDeleted            ProvisioningState = "Deleted"
	ProvisioningStateDeleting           ProvisioningState = "Deleting"
	ProvisioningStateDisconnected       ProvisioningState = "Disconnected"
	ProvisioningStateError              ProvisioningState = "Error"
	ProvisioningStateFailed             ProvisioningState = "Failed"
	ProvisioningStateInProgress         ProvisioningState = "InProgress"
	ProvisioningStateMoving             ProvisioningState = "Moving"
	ProvisioningStateNotSpecified       ProvisioningState = "NotSpecified"
	ProvisioningStatePartiallyConnected ProvisioningState = "PartiallyConnected"
	ProvisioningStatePartiallySucceeded ProvisioningState = "PartiallySucceeded"
	ProvisioningStateProvisioning       ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded          ProvisioningState = "Succeeded"
	ProvisioningStateUpdating           ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateConnected),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateDisconnected),
		string(ProvisioningStateError),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
		string(ProvisioningStateMoving),
		string(ProvisioningStateNotSpecified),
		string(ProvisioningStatePartiallyConnected),
		string(ProvisioningStatePartiallySucceeded),
		string(ProvisioningStateProvisioning),
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
		"accepted":           ProvisioningStateAccepted,
		"canceled":           ProvisioningStateCanceled,
		"connected":          ProvisioningStateConnected,
		"creating":           ProvisioningStateCreating,
		"deleted":            ProvisioningStateDeleted,
		"deleting":           ProvisioningStateDeleting,
		"disconnected":       ProvisioningStateDisconnected,
		"error":              ProvisioningStateError,
		"failed":             ProvisioningStateFailed,
		"inprogress":         ProvisioningStateInProgress,
		"moving":             ProvisioningStateMoving,
		"notspecified":       ProvisioningStateNotSpecified,
		"partiallyconnected": ProvisioningStatePartiallyConnected,
		"partiallysucceeded": ProvisioningStatePartiallySucceeded,
		"provisioning":       ProvisioningStateProvisioning,
		"succeeded":          ProvisioningStateSucceeded,
		"updating":           ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type StatusLevelTypes string

const (
	StatusLevelTypesError   StatusLevelTypes = "Error"
	StatusLevelTypesInfo    StatusLevelTypes = "Info"
	StatusLevelTypesWarning StatusLevelTypes = "Warning"
)

func PossibleValuesForStatusLevelTypes() []string {
	return []string{
		string(StatusLevelTypesError),
		string(StatusLevelTypesInfo),
		string(StatusLevelTypesWarning),
	}
}

func (s *StatusLevelTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusLevelTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusLevelTypes(input string) (*StatusLevelTypes, error) {
	vals := map[string]StatusLevelTypes{
		"error":   StatusLevelTypesError,
		"info":    StatusLevelTypesInfo,
		"warning": StatusLevelTypesWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusLevelTypes(input)
	return &out, nil
}
