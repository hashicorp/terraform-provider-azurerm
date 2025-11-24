package account

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountProvisioningState string

const (
	AccountProvisioningStateCanceled     AccountProvisioningState = "Canceled"
	AccountProvisioningStateCreating     AccountProvisioningState = "Creating"
	AccountProvisioningStateDeleting     AccountProvisioningState = "Deleting"
	AccountProvisioningStateFailed       AccountProvisioningState = "Failed"
	AccountProvisioningStateMoving       AccountProvisioningState = "Moving"
	AccountProvisioningStateSoftDeleted  AccountProvisioningState = "SoftDeleted"
	AccountProvisioningStateSoftDeleting AccountProvisioningState = "SoftDeleting"
	AccountProvisioningStateSucceeded    AccountProvisioningState = "Succeeded"
	AccountProvisioningStateUnknown      AccountProvisioningState = "Unknown"
	AccountProvisioningStateUpdating     AccountProvisioningState = "Updating"
)

func PossibleValuesForAccountProvisioningState() []string {
	return []string{
		string(AccountProvisioningStateCanceled),
		string(AccountProvisioningStateCreating),
		string(AccountProvisioningStateDeleting),
		string(AccountProvisioningStateFailed),
		string(AccountProvisioningStateMoving),
		string(AccountProvisioningStateSoftDeleted),
		string(AccountProvisioningStateSoftDeleting),
		string(AccountProvisioningStateSucceeded),
		string(AccountProvisioningStateUnknown),
		string(AccountProvisioningStateUpdating),
	}
}

func (s *AccountProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccountProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccountProvisioningState(input string) (*AccountProvisioningState, error) {
	vals := map[string]AccountProvisioningState{
		"canceled":     AccountProvisioningStateCanceled,
		"creating":     AccountProvisioningStateCreating,
		"deleting":     AccountProvisioningStateDeleting,
		"failed":       AccountProvisioningStateFailed,
		"moving":       AccountProvisioningStateMoving,
		"softdeleted":  AccountProvisioningStateSoftDeleted,
		"softdeleting": AccountProvisioningStateSoftDeleting,
		"succeeded":    AccountProvisioningStateSucceeded,
		"unknown":      AccountProvisioningStateUnknown,
		"updating":     AccountProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountProvisioningState(input)
	return &out, nil
}

type ManagedEventHubState string

const (
	ManagedEventHubStateDisabled     ManagedEventHubState = "Disabled"
	ManagedEventHubStateEnabled      ManagedEventHubState = "Enabled"
	ManagedEventHubStateNotSpecified ManagedEventHubState = "NotSpecified"
)

func PossibleValuesForManagedEventHubState() []string {
	return []string{
		string(ManagedEventHubStateDisabled),
		string(ManagedEventHubStateEnabled),
		string(ManagedEventHubStateNotSpecified),
	}
}

func (s *ManagedEventHubState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedEventHubState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedEventHubState(input string) (*ManagedEventHubState, error) {
	vals := map[string]ManagedEventHubState{
		"disabled":     ManagedEventHubStateDisabled,
		"enabled":      ManagedEventHubStateEnabled,
		"notspecified": ManagedEventHubStateNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedEventHubState(input)
	return &out, nil
}

type ManagedResourcesPublicNetworkAccess string

const (
	ManagedResourcesPublicNetworkAccessDisabled     ManagedResourcesPublicNetworkAccess = "Disabled"
	ManagedResourcesPublicNetworkAccessEnabled      ManagedResourcesPublicNetworkAccess = "Enabled"
	ManagedResourcesPublicNetworkAccessNotSpecified ManagedResourcesPublicNetworkAccess = "NotSpecified"
)

func PossibleValuesForManagedResourcesPublicNetworkAccess() []string {
	return []string{
		string(ManagedResourcesPublicNetworkAccessDisabled),
		string(ManagedResourcesPublicNetworkAccessEnabled),
		string(ManagedResourcesPublicNetworkAccessNotSpecified),
	}
}

func (s *ManagedResourcesPublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedResourcesPublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedResourcesPublicNetworkAccess(input string) (*ManagedResourcesPublicNetworkAccess, error) {
	vals := map[string]ManagedResourcesPublicNetworkAccess{
		"disabled":     ManagedResourcesPublicNetworkAccessDisabled,
		"enabled":      ManagedResourcesPublicNetworkAccessEnabled,
		"notspecified": ManagedResourcesPublicNetworkAccessNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedResourcesPublicNetworkAccess(input)
	return &out, nil
}

type Name string

const (
	NameStandard Name = "Standard"
)

func PossibleValuesForName() []string {
	return []string{
		string(NameStandard),
	}
}

func (s *Name) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseName(input string) (*Name, error) {
	vals := map[string]Name{
		"standard": NameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Name(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateMoving       ProvisioningState = "Moving"
	ProvisioningStateSoftDeleted  ProvisioningState = "SoftDeleted"
	ProvisioningStateSoftDeleting ProvisioningState = "SoftDeleting"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUnknown      ProvisioningState = "Unknown"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateSoftDeleted),
		string(ProvisioningStateSoftDeleting),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
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
		"canceled":     ProvisioningStateCanceled,
		"creating":     ProvisioningStateCreating,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"moving":       ProvisioningStateMoving,
		"softdeleted":  ProvisioningStateSoftDeleted,
		"softdeleting": ProvisioningStateSoftDeleting,
		"succeeded":    ProvisioningStateSucceeded,
		"unknown":      ProvisioningStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled     PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled      PublicNetworkAccess = "Enabled"
	PublicNetworkAccessNotSpecified PublicNetworkAccess = "NotSpecified"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
		string(PublicNetworkAccessNotSpecified),
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
		"disabled":     PublicNetworkAccessDisabled,
		"enabled":      PublicNetworkAccessEnabled,
		"notspecified": PublicNetworkAccessNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type Status string

const (
	StatusApproved     Status = "Approved"
	StatusDisconnected Status = "Disconnected"
	StatusPending      Status = "Pending"
	StatusRejected     Status = "Rejected"
	StatusUnknown      Status = "Unknown"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusApproved),
		string(StatusDisconnected),
		string(StatusPending),
		string(StatusRejected),
		string(StatusUnknown),
	}
}

func (s *Status) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"approved":     StatusApproved,
		"disconnected": StatusDisconnected,
		"pending":      StatusPending,
		"rejected":     StatusRejected,
		"unknown":      StatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}
