package namespaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessRights string

const (
	AccessRightsListen AccessRights = "Listen"
	AccessRightsManage AccessRights = "Manage"
	AccessRightsSend   AccessRights = "Send"
)

func PossibleValuesForAccessRights() []string {
	return []string{
		string(AccessRightsListen),
		string(AccessRightsManage),
		string(AccessRightsSend),
	}
}

func (s *AccessRights) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessRights(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessRights(input string) (*AccessRights, error) {
	vals := map[string]AccessRights{
		"listen": AccessRightsListen,
		"manage": AccessRightsManage,
		"send":   AccessRightsSend,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessRights(input)
	return &out, nil
}

type NamespaceStatus string

const (
	NamespaceStatusCreated   NamespaceStatus = "Created"
	NamespaceStatusCreating  NamespaceStatus = "Creating"
	NamespaceStatusDeleting  NamespaceStatus = "Deleting"
	NamespaceStatusSuspended NamespaceStatus = "Suspended"
)

func PossibleValuesForNamespaceStatus() []string {
	return []string{
		string(NamespaceStatusCreated),
		string(NamespaceStatusCreating),
		string(NamespaceStatusDeleting),
		string(NamespaceStatusSuspended),
	}
}

func (s *NamespaceStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNamespaceStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNamespaceStatus(input string) (*NamespaceStatus, error) {
	vals := map[string]NamespaceStatus{
		"created":   NamespaceStatusCreated,
		"creating":  NamespaceStatusCreating,
		"deleting":  NamespaceStatusDeleting,
		"suspended": NamespaceStatusSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NamespaceStatus(input)
	return &out, nil
}

type NamespaceType string

const (
	NamespaceTypeMessaging       NamespaceType = "Messaging"
	NamespaceTypeNotificationHub NamespaceType = "NotificationHub"
)

func PossibleValuesForNamespaceType() []string {
	return []string{
		string(NamespaceTypeMessaging),
		string(NamespaceTypeNotificationHub),
	}
}

func (s *NamespaceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNamespaceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNamespaceType(input string) (*NamespaceType, error) {
	vals := map[string]NamespaceType{
		"messaging":       NamespaceTypeMessaging,
		"notificationhub": NamespaceTypeNotificationHub,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NamespaceType(input)
	return &out, nil
}

type OperationProvisioningState string

const (
	OperationProvisioningStateCanceled   OperationProvisioningState = "Canceled"
	OperationProvisioningStateDisabled   OperationProvisioningState = "Disabled"
	OperationProvisioningStateFailed     OperationProvisioningState = "Failed"
	OperationProvisioningStateInProgress OperationProvisioningState = "InProgress"
	OperationProvisioningStatePending    OperationProvisioningState = "Pending"
	OperationProvisioningStateSucceeded  OperationProvisioningState = "Succeeded"
	OperationProvisioningStateUnknown    OperationProvisioningState = "Unknown"
)

func PossibleValuesForOperationProvisioningState() []string {
	return []string{
		string(OperationProvisioningStateCanceled),
		string(OperationProvisioningStateDisabled),
		string(OperationProvisioningStateFailed),
		string(OperationProvisioningStateInProgress),
		string(OperationProvisioningStatePending),
		string(OperationProvisioningStateSucceeded),
		string(OperationProvisioningStateUnknown),
	}
}

func (s *OperationProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationProvisioningState(input string) (*OperationProvisioningState, error) {
	vals := map[string]OperationProvisioningState{
		"canceled":   OperationProvisioningStateCanceled,
		"disabled":   OperationProvisioningStateDisabled,
		"failed":     OperationProvisioningStateFailed,
		"inprogress": OperationProvisioningStateInProgress,
		"pending":    OperationProvisioningStatePending,
		"succeeded":  OperationProvisioningStateSucceeded,
		"unknown":    OperationProvisioningStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationProvisioningState(input)
	return &out, nil
}

type PolicyKeyType string

const (
	PolicyKeyTypePrimaryKey   PolicyKeyType = "PrimaryKey"
	PolicyKeyTypeSecondaryKey PolicyKeyType = "SecondaryKey"
)

func PossibleValuesForPolicyKeyType() []string {
	return []string{
		string(PolicyKeyTypePrimaryKey),
		string(PolicyKeyTypeSecondaryKey),
	}
}

func (s *PolicyKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyKeyType(input string) (*PolicyKeyType, error) {
	vals := map[string]PolicyKeyType{
		"primarykey":   PolicyKeyTypePrimaryKey,
		"secondarykey": PolicyKeyTypeSecondaryKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyKeyType(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating        PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleted         PrivateEndpointConnectionProvisioningState = "Deleted"
	PrivateEndpointConnectionProvisioningStateDeleting        PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateDeletingByProxy PrivateEndpointConnectionProvisioningState = "DeletingByProxy"
	PrivateEndpointConnectionProvisioningStateSucceeded       PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUnknown         PrivateEndpointConnectionProvisioningState = "Unknown"
	PrivateEndpointConnectionProvisioningStateUpdating        PrivateEndpointConnectionProvisioningState = "Updating"
	PrivateEndpointConnectionProvisioningStateUpdatingByProxy PrivateEndpointConnectionProvisioningState = "UpdatingByProxy"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleted),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateDeletingByProxy),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
		string(PrivateEndpointConnectionProvisioningStateUnknown),
		string(PrivateEndpointConnectionProvisioningStateUpdating),
		string(PrivateEndpointConnectionProvisioningStateUpdatingByProxy),
	}
}

func (s *PrivateEndpointConnectionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointConnectionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":        PrivateEndpointConnectionProvisioningStateCreating,
		"deleted":         PrivateEndpointConnectionProvisioningStateDeleted,
		"deleting":        PrivateEndpointConnectionProvisioningStateDeleting,
		"deletingbyproxy": PrivateEndpointConnectionProvisioningStateDeletingByProxy,
		"succeeded":       PrivateEndpointConnectionProvisioningStateSucceeded,
		"unknown":         PrivateEndpointConnectionProvisioningStateUnknown,
		"updating":        PrivateEndpointConnectionProvisioningStateUpdating,
		"updatingbyproxy": PrivateEndpointConnectionProvisioningStateUpdatingByProxy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateLinkConnectionStatus string

const (
	PrivateLinkConnectionStatusApproved     PrivateLinkConnectionStatus = "Approved"
	PrivateLinkConnectionStatusDisconnected PrivateLinkConnectionStatus = "Disconnected"
	PrivateLinkConnectionStatusPending      PrivateLinkConnectionStatus = "Pending"
	PrivateLinkConnectionStatusRejected     PrivateLinkConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateLinkConnectionStatus() []string {
	return []string{
		string(PrivateLinkConnectionStatusApproved),
		string(PrivateLinkConnectionStatusDisconnected),
		string(PrivateLinkConnectionStatusPending),
		string(PrivateLinkConnectionStatusRejected),
	}
}

func (s *PrivateLinkConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateLinkConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateLinkConnectionStatus(input string) (*PrivateLinkConnectionStatus, error) {
	vals := map[string]PrivateLinkConnectionStatus{
		"approved":     PrivateLinkConnectionStatusApproved,
		"disconnected": PrivateLinkConnectionStatusDisconnected,
		"pending":      PrivateLinkConnectionStatusPending,
		"rejected":     PrivateLinkConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateLinkConnectionStatus(input)
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

type ReplicationRegion string

const (
	ReplicationRegionAustraliaEast    ReplicationRegion = "AustraliaEast"
	ReplicationRegionBrazilSouth      ReplicationRegion = "BrazilSouth"
	ReplicationRegionDefault          ReplicationRegion = "Default"
	ReplicationRegionNone             ReplicationRegion = "None"
	ReplicationRegionNorthEurope      ReplicationRegion = "NorthEurope"
	ReplicationRegionSouthAfricaNorth ReplicationRegion = "SouthAfricaNorth"
	ReplicationRegionSouthEastAsia    ReplicationRegion = "SouthEastAsia"
	ReplicationRegionWestUsTwo        ReplicationRegion = "WestUs2"
)

func PossibleValuesForReplicationRegion() []string {
	return []string{
		string(ReplicationRegionAustraliaEast),
		string(ReplicationRegionBrazilSouth),
		string(ReplicationRegionDefault),
		string(ReplicationRegionNone),
		string(ReplicationRegionNorthEurope),
		string(ReplicationRegionSouthAfricaNorth),
		string(ReplicationRegionSouthEastAsia),
		string(ReplicationRegionWestUsTwo),
	}
}

func (s *ReplicationRegion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationRegion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationRegion(input string) (*ReplicationRegion, error) {
	vals := map[string]ReplicationRegion{
		"australiaeast":    ReplicationRegionAustraliaEast,
		"brazilsouth":      ReplicationRegionBrazilSouth,
		"default":          ReplicationRegionDefault,
		"none":             ReplicationRegionNone,
		"northeurope":      ReplicationRegionNorthEurope,
		"southafricanorth": ReplicationRegionSouthAfricaNorth,
		"southeastasia":    ReplicationRegionSouthEastAsia,
		"westus2":          ReplicationRegionWestUsTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationRegion(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameBasic    SkuName = "Basic"
	SkuNameFree     SkuName = "Free"
	SkuNameStandard SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameBasic),
		string(SkuNameFree),
		string(SkuNameStandard),
	}
}

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"basic":    SkuNameBasic,
		"free":     SkuNameFree,
		"standard": SkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}

type ZoneRedundancyPreference string

const (
	ZoneRedundancyPreferenceDisabled ZoneRedundancyPreference = "Disabled"
	ZoneRedundancyPreferenceEnabled  ZoneRedundancyPreference = "Enabled"
)

func PossibleValuesForZoneRedundancyPreference() []string {
	return []string{
		string(ZoneRedundancyPreferenceDisabled),
		string(ZoneRedundancyPreferenceEnabled),
	}
}

func (s *ZoneRedundancyPreference) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseZoneRedundancyPreference(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseZoneRedundancyPreference(input string) (*ZoneRedundancyPreference, error) {
	vals := map[string]ZoneRedundancyPreference{
		"disabled": ZoneRedundancyPreferenceDisabled,
		"enabled":  ZoneRedundancyPreferenceEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ZoneRedundancyPreference(input)
	return &out, nil
}
