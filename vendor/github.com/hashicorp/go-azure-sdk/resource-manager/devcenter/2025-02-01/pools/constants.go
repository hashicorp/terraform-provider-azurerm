package pools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "Healthy"
	HealthStatusPending   HealthStatus = "Pending"
	HealthStatusUnhealthy HealthStatus = "Unhealthy"
	HealthStatusUnknown   HealthStatus = "Unknown"
	HealthStatusWarning   HealthStatus = "Warning"
)

func PossibleValuesForHealthStatus() []string {
	return []string{
		string(HealthStatusHealthy),
		string(HealthStatusPending),
		string(HealthStatusUnhealthy),
		string(HealthStatusUnknown),
		string(HealthStatusWarning),
	}
}

func (s *HealthStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthStatus(input string) (*HealthStatus, error) {
	vals := map[string]HealthStatus{
		"healthy":   HealthStatusHealthy,
		"pending":   HealthStatusPending,
		"unhealthy": HealthStatusUnhealthy,
		"unknown":   HealthStatusUnknown,
		"warning":   HealthStatusWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthStatus(input)
	return &out, nil
}

type LicenseType string

const (
	LicenseTypeWindowsClient LicenseType = "Windows_Client"
)

func PossibleValuesForLicenseType() []string {
	return []string{
		string(LicenseTypeWindowsClient),
	}
}

func (s *LicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseType(input string) (*LicenseType, error) {
	vals := map[string]LicenseType{
		"windows_client": LicenseTypeWindowsClient,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseType(input)
	return &out, nil
}

type LocalAdminStatus string

const (
	LocalAdminStatusDisabled LocalAdminStatus = "Disabled"
	LocalAdminStatusEnabled  LocalAdminStatus = "Enabled"
)

func PossibleValuesForLocalAdminStatus() []string {
	return []string{
		string(LocalAdminStatusDisabled),
		string(LocalAdminStatusEnabled),
	}
}

func (s *LocalAdminStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalAdminStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalAdminStatus(input string) (*LocalAdminStatus, error) {
	vals := map[string]LocalAdminStatus{
		"disabled": LocalAdminStatusDisabled,
		"enabled":  LocalAdminStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalAdminStatus(input)
	return &out, nil
}

type PoolDevBoxDefinitionType string

const (
	PoolDevBoxDefinitionTypeReference PoolDevBoxDefinitionType = "Reference"
	PoolDevBoxDefinitionTypeValue     PoolDevBoxDefinitionType = "Value"
)

func PossibleValuesForPoolDevBoxDefinitionType() []string {
	return []string{
		string(PoolDevBoxDefinitionTypeReference),
		string(PoolDevBoxDefinitionTypeValue),
	}
}

func (s *PoolDevBoxDefinitionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePoolDevBoxDefinitionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePoolDevBoxDefinitionType(input string) (*PoolDevBoxDefinitionType, error) {
	vals := map[string]PoolDevBoxDefinitionType{
		"reference": PoolDevBoxDefinitionTypeReference,
		"value":     PoolDevBoxDefinitionTypeValue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PoolDevBoxDefinitionType(input)
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

type SingleSignOnStatus string

const (
	SingleSignOnStatusDisabled SingleSignOnStatus = "Disabled"
	SingleSignOnStatusEnabled  SingleSignOnStatus = "Enabled"
)

func PossibleValuesForSingleSignOnStatus() []string {
	return []string{
		string(SingleSignOnStatusDisabled),
		string(SingleSignOnStatusEnabled),
	}
}

func (s *SingleSignOnStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSingleSignOnStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSingleSignOnStatus(input string) (*SingleSignOnStatus, error) {
	vals := map[string]SingleSignOnStatus{
		"disabled": SingleSignOnStatusDisabled,
		"enabled":  SingleSignOnStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SingleSignOnStatus(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierFree),
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    SkuTierBasic,
		"free":     SkuTierFree,
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type StopOnDisconnectEnableStatus string

const (
	StopOnDisconnectEnableStatusDisabled StopOnDisconnectEnableStatus = "Disabled"
	StopOnDisconnectEnableStatusEnabled  StopOnDisconnectEnableStatus = "Enabled"
)

func PossibleValuesForStopOnDisconnectEnableStatus() []string {
	return []string{
		string(StopOnDisconnectEnableStatusDisabled),
		string(StopOnDisconnectEnableStatusEnabled),
	}
}

func (s *StopOnDisconnectEnableStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStopOnDisconnectEnableStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStopOnDisconnectEnableStatus(input string) (*StopOnDisconnectEnableStatus, error) {
	vals := map[string]StopOnDisconnectEnableStatus{
		"disabled": StopOnDisconnectEnableStatusDisabled,
		"enabled":  StopOnDisconnectEnableStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StopOnDisconnectEnableStatus(input)
	return &out, nil
}

type StopOnNoConnectEnableStatus string

const (
	StopOnNoConnectEnableStatusDisabled StopOnNoConnectEnableStatus = "Disabled"
	StopOnNoConnectEnableStatusEnabled  StopOnNoConnectEnableStatus = "Enabled"
)

func PossibleValuesForStopOnNoConnectEnableStatus() []string {
	return []string{
		string(StopOnNoConnectEnableStatusDisabled),
		string(StopOnNoConnectEnableStatusEnabled),
	}
}

func (s *StopOnNoConnectEnableStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStopOnNoConnectEnableStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStopOnNoConnectEnableStatus(input string) (*StopOnNoConnectEnableStatus, error) {
	vals := map[string]StopOnNoConnectEnableStatus{
		"disabled": StopOnNoConnectEnableStatusDisabled,
		"enabled":  StopOnNoConnectEnableStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StopOnNoConnectEnableStatus(input)
	return &out, nil
}

type VirtualNetworkType string

const (
	VirtualNetworkTypeManaged   VirtualNetworkType = "Managed"
	VirtualNetworkTypeUnmanaged VirtualNetworkType = "Unmanaged"
)

func PossibleValuesForVirtualNetworkType() []string {
	return []string{
		string(VirtualNetworkTypeManaged),
		string(VirtualNetworkTypeUnmanaged),
	}
}

func (s *VirtualNetworkType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkType(input string) (*VirtualNetworkType, error) {
	vals := map[string]VirtualNetworkType{
		"managed":   VirtualNetworkTypeManaged,
		"unmanaged": VirtualNetworkTypeUnmanaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkType(input)
	return &out, nil
}
