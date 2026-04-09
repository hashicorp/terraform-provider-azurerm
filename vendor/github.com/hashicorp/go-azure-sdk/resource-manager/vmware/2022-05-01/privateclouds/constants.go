package privateclouds

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailabilityStrategy string

const (
	AvailabilityStrategyDualZone   AvailabilityStrategy = "DualZone"
	AvailabilityStrategySingleZone AvailabilityStrategy = "SingleZone"
)

func PossibleValuesForAvailabilityStrategy() []string {
	return []string{
		string(AvailabilityStrategyDualZone),
		string(AvailabilityStrategySingleZone),
	}
}

func (s *AvailabilityStrategy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAvailabilityStrategy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAvailabilityStrategy(input string) (*AvailabilityStrategy, error) {
	vals := map[string]AvailabilityStrategy{
		"dualzone":   AvailabilityStrategyDualZone,
		"singlezone": AvailabilityStrategySingleZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AvailabilityStrategy(input)
	return &out, nil
}

type ClusterProvisioningState string

const (
	ClusterProvisioningStateCanceled  ClusterProvisioningState = "Canceled"
	ClusterProvisioningStateCancelled ClusterProvisioningState = "Cancelled"
	ClusterProvisioningStateDeleting  ClusterProvisioningState = "Deleting"
	ClusterProvisioningStateFailed    ClusterProvisioningState = "Failed"
	ClusterProvisioningStateSucceeded ClusterProvisioningState = "Succeeded"
	ClusterProvisioningStateUpdating  ClusterProvisioningState = "Updating"
)

func PossibleValuesForClusterProvisioningState() []string {
	return []string{
		string(ClusterProvisioningStateCanceled),
		string(ClusterProvisioningStateCancelled),
		string(ClusterProvisioningStateDeleting),
		string(ClusterProvisioningStateFailed),
		string(ClusterProvisioningStateSucceeded),
		string(ClusterProvisioningStateUpdating),
	}
}

func (s *ClusterProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterProvisioningState(input string) (*ClusterProvisioningState, error) {
	vals := map[string]ClusterProvisioningState{
		"canceled":  ClusterProvisioningStateCanceled,
		"cancelled": ClusterProvisioningStateCancelled,
		"deleting":  ClusterProvisioningStateDeleting,
		"failed":    ClusterProvisioningStateFailed,
		"succeeded": ClusterProvisioningStateSucceeded,
		"updating":  ClusterProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterProvisioningState(input)
	return &out, nil
}

type EncryptionKeyStatus string

const (
	EncryptionKeyStatusAccessDenied EncryptionKeyStatus = "AccessDenied"
	EncryptionKeyStatusConnected    EncryptionKeyStatus = "Connected"
)

func PossibleValuesForEncryptionKeyStatus() []string {
	return []string{
		string(EncryptionKeyStatusAccessDenied),
		string(EncryptionKeyStatusConnected),
	}
}

func (s *EncryptionKeyStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionKeyStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionKeyStatus(input string) (*EncryptionKeyStatus, error) {
	vals := map[string]EncryptionKeyStatus{
		"accessdenied": EncryptionKeyStatusAccessDenied,
		"connected":    EncryptionKeyStatusConnected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionKeyStatus(input)
	return &out, nil
}

type EncryptionState string

const (
	EncryptionStateDisabled EncryptionState = "Disabled"
	EncryptionStateEnabled  EncryptionState = "Enabled"
)

func PossibleValuesForEncryptionState() []string {
	return []string{
		string(EncryptionStateDisabled),
		string(EncryptionStateEnabled),
	}
}

func (s *EncryptionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionState(input string) (*EncryptionState, error) {
	vals := map[string]EncryptionState{
		"disabled": EncryptionStateDisabled,
		"enabled":  EncryptionStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionState(input)
	return &out, nil
}

type EncryptionVersionType string

const (
	EncryptionVersionTypeAutoDetected EncryptionVersionType = "AutoDetected"
	EncryptionVersionTypeFixed        EncryptionVersionType = "Fixed"
)

func PossibleValuesForEncryptionVersionType() []string {
	return []string{
		string(EncryptionVersionTypeAutoDetected),
		string(EncryptionVersionTypeFixed),
	}
}

func (s *EncryptionVersionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionVersionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionVersionType(input string) (*EncryptionVersionType, error) {
	vals := map[string]EncryptionVersionType{
		"autodetected": EncryptionVersionTypeAutoDetected,
		"fixed":        EncryptionVersionTypeFixed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionVersionType(input)
	return &out, nil
}

type InternetEnum string

const (
	InternetEnumDisabled InternetEnum = "Disabled"
	InternetEnumEnabled  InternetEnum = "Enabled"
)

func PossibleValuesForInternetEnum() []string {
	return []string{
		string(InternetEnumDisabled),
		string(InternetEnumEnabled),
	}
}

func (s *InternetEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInternetEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInternetEnum(input string) (*InternetEnum, error) {
	vals := map[string]InternetEnum{
		"disabled": InternetEnumDisabled,
		"enabled":  InternetEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InternetEnum(input)
	return &out, nil
}

type NsxPublicIPQuotaRaisedEnum string

const (
	NsxPublicIPQuotaRaisedEnumDisabled NsxPublicIPQuotaRaisedEnum = "Disabled"
	NsxPublicIPQuotaRaisedEnumEnabled  NsxPublicIPQuotaRaisedEnum = "Enabled"
)

func PossibleValuesForNsxPublicIPQuotaRaisedEnum() []string {
	return []string{
		string(NsxPublicIPQuotaRaisedEnumDisabled),
		string(NsxPublicIPQuotaRaisedEnumEnabled),
	}
}

func (s *NsxPublicIPQuotaRaisedEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNsxPublicIPQuotaRaisedEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNsxPublicIPQuotaRaisedEnum(input string) (*NsxPublicIPQuotaRaisedEnum, error) {
	vals := map[string]NsxPublicIPQuotaRaisedEnum{
		"disabled": NsxPublicIPQuotaRaisedEnumDisabled,
		"enabled":  NsxPublicIPQuotaRaisedEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NsxPublicIPQuotaRaisedEnum(input)
	return &out, nil
}

type PrivateCloudProvisioningState string

const (
	PrivateCloudProvisioningStateBuilding  PrivateCloudProvisioningState = "Building"
	PrivateCloudProvisioningStateCanceled  PrivateCloudProvisioningState = "Canceled"
	PrivateCloudProvisioningStateCancelled PrivateCloudProvisioningState = "Cancelled"
	PrivateCloudProvisioningStateDeleting  PrivateCloudProvisioningState = "Deleting"
	PrivateCloudProvisioningStateFailed    PrivateCloudProvisioningState = "Failed"
	PrivateCloudProvisioningStatePending   PrivateCloudProvisioningState = "Pending"
	PrivateCloudProvisioningStateSucceeded PrivateCloudProvisioningState = "Succeeded"
	PrivateCloudProvisioningStateUpdating  PrivateCloudProvisioningState = "Updating"
)

func PossibleValuesForPrivateCloudProvisioningState() []string {
	return []string{
		string(PrivateCloudProvisioningStateBuilding),
		string(PrivateCloudProvisioningStateCanceled),
		string(PrivateCloudProvisioningStateCancelled),
		string(PrivateCloudProvisioningStateDeleting),
		string(PrivateCloudProvisioningStateFailed),
		string(PrivateCloudProvisioningStatePending),
		string(PrivateCloudProvisioningStateSucceeded),
		string(PrivateCloudProvisioningStateUpdating),
	}
}

func (s *PrivateCloudProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateCloudProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateCloudProvisioningState(input string) (*PrivateCloudProvisioningState, error) {
	vals := map[string]PrivateCloudProvisioningState{
		"building":  PrivateCloudProvisioningStateBuilding,
		"canceled":  PrivateCloudProvisioningStateCanceled,
		"cancelled": PrivateCloudProvisioningStateCancelled,
		"deleting":  PrivateCloudProvisioningStateDeleting,
		"failed":    PrivateCloudProvisioningStateFailed,
		"pending":   PrivateCloudProvisioningStatePending,
		"succeeded": PrivateCloudProvisioningStateSucceeded,
		"updating":  PrivateCloudProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateCloudProvisioningState(input)
	return &out, nil
}

type SslEnum string

const (
	SslEnumDisabled SslEnum = "Disabled"
	SslEnumEnabled  SslEnum = "Enabled"
)

func PossibleValuesForSslEnum() []string {
	return []string{
		string(SslEnumDisabled),
		string(SslEnumEnabled),
	}
}

func (s *SslEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSslEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSslEnum(input string) (*SslEnum, error) {
	vals := map[string]SslEnum{
		"disabled": SslEnumDisabled,
		"enabled":  SslEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SslEnum(input)
	return &out, nil
}
