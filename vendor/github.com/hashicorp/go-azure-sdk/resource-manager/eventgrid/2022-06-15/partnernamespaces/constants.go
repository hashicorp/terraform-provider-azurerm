package partnernamespaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPActionType string

const (
	IPActionTypeAllow IPActionType = "Allow"
)

func PossibleValuesForIPActionType() []string {
	return []string{
		string(IPActionTypeAllow),
	}
}

func (s *IPActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPActionType(input string) (*IPActionType, error) {
	vals := map[string]IPActionType{
		"allow": IPActionTypeAllow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPActionType(input)
	return &out, nil
}

type PartnerNamespaceProvisioningState string

const (
	PartnerNamespaceProvisioningStateCanceled  PartnerNamespaceProvisioningState = "Canceled"
	PartnerNamespaceProvisioningStateCreating  PartnerNamespaceProvisioningState = "Creating"
	PartnerNamespaceProvisioningStateDeleting  PartnerNamespaceProvisioningState = "Deleting"
	PartnerNamespaceProvisioningStateFailed    PartnerNamespaceProvisioningState = "Failed"
	PartnerNamespaceProvisioningStateSucceeded PartnerNamespaceProvisioningState = "Succeeded"
	PartnerNamespaceProvisioningStateUpdating  PartnerNamespaceProvisioningState = "Updating"
)

func PossibleValuesForPartnerNamespaceProvisioningState() []string {
	return []string{
		string(PartnerNamespaceProvisioningStateCanceled),
		string(PartnerNamespaceProvisioningStateCreating),
		string(PartnerNamespaceProvisioningStateDeleting),
		string(PartnerNamespaceProvisioningStateFailed),
		string(PartnerNamespaceProvisioningStateSucceeded),
		string(PartnerNamespaceProvisioningStateUpdating),
	}
}

func (s *PartnerNamespaceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePartnerNamespaceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePartnerNamespaceProvisioningState(input string) (*PartnerNamespaceProvisioningState, error) {
	vals := map[string]PartnerNamespaceProvisioningState{
		"canceled":  PartnerNamespaceProvisioningStateCanceled,
		"creating":  PartnerNamespaceProvisioningStateCreating,
		"deleting":  PartnerNamespaceProvisioningStateDeleting,
		"failed":    PartnerNamespaceProvisioningStateFailed,
		"succeeded": PartnerNamespaceProvisioningStateSucceeded,
		"updating":  PartnerNamespaceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerNamespaceProvisioningState(input)
	return &out, nil
}

type PartnerTopicRoutingMode string

const (
	PartnerTopicRoutingModeChannelNameHeader    PartnerTopicRoutingMode = "ChannelNameHeader"
	PartnerTopicRoutingModeSourceEventAttribute PartnerTopicRoutingMode = "SourceEventAttribute"
)

func PossibleValuesForPartnerTopicRoutingMode() []string {
	return []string{
		string(PartnerTopicRoutingModeChannelNameHeader),
		string(PartnerTopicRoutingModeSourceEventAttribute),
	}
}

func (s *PartnerTopicRoutingMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePartnerTopicRoutingMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePartnerTopicRoutingMode(input string) (*PartnerTopicRoutingMode, error) {
	vals := map[string]PartnerTopicRoutingMode{
		"channelnameheader":    PartnerTopicRoutingModeChannelNameHeader,
		"sourceeventattribute": PartnerTopicRoutingModeSourceEventAttribute,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartnerTopicRoutingMode(input)
	return &out, nil
}

type PersistedConnectionStatus string

const (
	PersistedConnectionStatusApproved     PersistedConnectionStatus = "Approved"
	PersistedConnectionStatusDisconnected PersistedConnectionStatus = "Disconnected"
	PersistedConnectionStatusPending      PersistedConnectionStatus = "Pending"
	PersistedConnectionStatusRejected     PersistedConnectionStatus = "Rejected"
)

func PossibleValuesForPersistedConnectionStatus() []string {
	return []string{
		string(PersistedConnectionStatusApproved),
		string(PersistedConnectionStatusDisconnected),
		string(PersistedConnectionStatusPending),
		string(PersistedConnectionStatusRejected),
	}
}

func (s *PersistedConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePersistedConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePersistedConnectionStatus(input string) (*PersistedConnectionStatus, error) {
	vals := map[string]PersistedConnectionStatus{
		"approved":     PersistedConnectionStatusApproved,
		"disconnected": PersistedConnectionStatusDisconnected,
		"pending":      PersistedConnectionStatusPending,
		"rejected":     PersistedConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PersistedConnectionStatus(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateCreating  ResourceProvisioningState = "Creating"
	ResourceProvisioningStateDeleting  ResourceProvisioningState = "Deleting"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
	ResourceProvisioningStateUpdating  ResourceProvisioningState = "Updating"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateCreating),
		string(ResourceProvisioningStateDeleting),
		string(ResourceProvisioningStateFailed),
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
		"canceled":  ResourceProvisioningStateCanceled,
		"creating":  ResourceProvisioningStateCreating,
		"deleting":  ResourceProvisioningStateDeleting,
		"failed":    ResourceProvisioningStateFailed,
		"succeeded": ResourceProvisioningStateSucceeded,
		"updating":  ResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}
