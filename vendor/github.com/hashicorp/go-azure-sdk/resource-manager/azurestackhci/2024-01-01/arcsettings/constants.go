package arcsettings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArcSettingAggregateState string

const (
	ArcSettingAggregateStateAccepted           ArcSettingAggregateState = "Accepted"
	ArcSettingAggregateStateCanceled           ArcSettingAggregateState = "Canceled"
	ArcSettingAggregateStateConnected          ArcSettingAggregateState = "Connected"
	ArcSettingAggregateStateCreating           ArcSettingAggregateState = "Creating"
	ArcSettingAggregateStateDeleted            ArcSettingAggregateState = "Deleted"
	ArcSettingAggregateStateDeleting           ArcSettingAggregateState = "Deleting"
	ArcSettingAggregateStateDisableInProgress  ArcSettingAggregateState = "DisableInProgress"
	ArcSettingAggregateStateDisconnected       ArcSettingAggregateState = "Disconnected"
	ArcSettingAggregateStateError              ArcSettingAggregateState = "Error"
	ArcSettingAggregateStateFailed             ArcSettingAggregateState = "Failed"
	ArcSettingAggregateStateInProgress         ArcSettingAggregateState = "InProgress"
	ArcSettingAggregateStateMoving             ArcSettingAggregateState = "Moving"
	ArcSettingAggregateStateNotSpecified       ArcSettingAggregateState = "NotSpecified"
	ArcSettingAggregateStatePartiallyConnected ArcSettingAggregateState = "PartiallyConnected"
	ArcSettingAggregateStatePartiallySucceeded ArcSettingAggregateState = "PartiallySucceeded"
	ArcSettingAggregateStateProvisioning       ArcSettingAggregateState = "Provisioning"
	ArcSettingAggregateStateSucceeded          ArcSettingAggregateState = "Succeeded"
	ArcSettingAggregateStateUpdating           ArcSettingAggregateState = "Updating"
)

func PossibleValuesForArcSettingAggregateState() []string {
	return []string{
		string(ArcSettingAggregateStateAccepted),
		string(ArcSettingAggregateStateCanceled),
		string(ArcSettingAggregateStateConnected),
		string(ArcSettingAggregateStateCreating),
		string(ArcSettingAggregateStateDeleted),
		string(ArcSettingAggregateStateDeleting),
		string(ArcSettingAggregateStateDisableInProgress),
		string(ArcSettingAggregateStateDisconnected),
		string(ArcSettingAggregateStateError),
		string(ArcSettingAggregateStateFailed),
		string(ArcSettingAggregateStateInProgress),
		string(ArcSettingAggregateStateMoving),
		string(ArcSettingAggregateStateNotSpecified),
		string(ArcSettingAggregateStatePartiallyConnected),
		string(ArcSettingAggregateStatePartiallySucceeded),
		string(ArcSettingAggregateStateProvisioning),
		string(ArcSettingAggregateStateSucceeded),
		string(ArcSettingAggregateStateUpdating),
	}
}

func (s *ArcSettingAggregateState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseArcSettingAggregateState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseArcSettingAggregateState(input string) (*ArcSettingAggregateState, error) {
	vals := map[string]ArcSettingAggregateState{
		"accepted":           ArcSettingAggregateStateAccepted,
		"canceled":           ArcSettingAggregateStateCanceled,
		"connected":          ArcSettingAggregateStateConnected,
		"creating":           ArcSettingAggregateStateCreating,
		"deleted":            ArcSettingAggregateStateDeleted,
		"deleting":           ArcSettingAggregateStateDeleting,
		"disableinprogress":  ArcSettingAggregateStateDisableInProgress,
		"disconnected":       ArcSettingAggregateStateDisconnected,
		"error":              ArcSettingAggregateStateError,
		"failed":             ArcSettingAggregateStateFailed,
		"inprogress":         ArcSettingAggregateStateInProgress,
		"moving":             ArcSettingAggregateStateMoving,
		"notspecified":       ArcSettingAggregateStateNotSpecified,
		"partiallyconnected": ArcSettingAggregateStatePartiallyConnected,
		"partiallysucceeded": ArcSettingAggregateStatePartiallySucceeded,
		"provisioning":       ArcSettingAggregateStateProvisioning,
		"succeeded":          ArcSettingAggregateStateSucceeded,
		"updating":           ArcSettingAggregateStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ArcSettingAggregateState(input)
	return &out, nil
}

type NodeArcState string

const (
	NodeArcStateAccepted           NodeArcState = "Accepted"
	NodeArcStateCanceled           NodeArcState = "Canceled"
	NodeArcStateConnected          NodeArcState = "Connected"
	NodeArcStateCreating           NodeArcState = "Creating"
	NodeArcStateDeleted            NodeArcState = "Deleted"
	NodeArcStateDeleting           NodeArcState = "Deleting"
	NodeArcStateDisableInProgress  NodeArcState = "DisableInProgress"
	NodeArcStateDisconnected       NodeArcState = "Disconnected"
	NodeArcStateError              NodeArcState = "Error"
	NodeArcStateFailed             NodeArcState = "Failed"
	NodeArcStateInProgress         NodeArcState = "InProgress"
	NodeArcStateMoving             NodeArcState = "Moving"
	NodeArcStateNotSpecified       NodeArcState = "NotSpecified"
	NodeArcStatePartiallyConnected NodeArcState = "PartiallyConnected"
	NodeArcStatePartiallySucceeded NodeArcState = "PartiallySucceeded"
	NodeArcStateProvisioning       NodeArcState = "Provisioning"
	NodeArcStateSucceeded          NodeArcState = "Succeeded"
	NodeArcStateUpdating           NodeArcState = "Updating"
)

func PossibleValuesForNodeArcState() []string {
	return []string{
		string(NodeArcStateAccepted),
		string(NodeArcStateCanceled),
		string(NodeArcStateConnected),
		string(NodeArcStateCreating),
		string(NodeArcStateDeleted),
		string(NodeArcStateDeleting),
		string(NodeArcStateDisableInProgress),
		string(NodeArcStateDisconnected),
		string(NodeArcStateError),
		string(NodeArcStateFailed),
		string(NodeArcStateInProgress),
		string(NodeArcStateMoving),
		string(NodeArcStateNotSpecified),
		string(NodeArcStatePartiallyConnected),
		string(NodeArcStatePartiallySucceeded),
		string(NodeArcStateProvisioning),
		string(NodeArcStateSucceeded),
		string(NodeArcStateUpdating),
	}
}

func (s *NodeArcState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNodeArcState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNodeArcState(input string) (*NodeArcState, error) {
	vals := map[string]NodeArcState{
		"accepted":           NodeArcStateAccepted,
		"canceled":           NodeArcStateCanceled,
		"connected":          NodeArcStateConnected,
		"creating":           NodeArcStateCreating,
		"deleted":            NodeArcStateDeleted,
		"deleting":           NodeArcStateDeleting,
		"disableinprogress":  NodeArcStateDisableInProgress,
		"disconnected":       NodeArcStateDisconnected,
		"error":              NodeArcStateError,
		"failed":             NodeArcStateFailed,
		"inprogress":         NodeArcStateInProgress,
		"moving":             NodeArcStateMoving,
		"notspecified":       NodeArcStateNotSpecified,
		"partiallyconnected": NodeArcStatePartiallyConnected,
		"partiallysucceeded": NodeArcStatePartiallySucceeded,
		"provisioning":       NodeArcStateProvisioning,
		"succeeded":          NodeArcStateSucceeded,
		"updating":           NodeArcStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodeArcState(input)
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
	ProvisioningStateDisableInProgress  ProvisioningState = "DisableInProgress"
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
		string(ProvisioningStateDisableInProgress),
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
		"disableinprogress":  ProvisioningStateDisableInProgress,
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
