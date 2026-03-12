package mongoclusters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationMode string

const (
	AuthenticationModeMicrosoftEntraID AuthenticationMode = "MicrosoftEntraID"
	AuthenticationModeNativeAuth       AuthenticationMode = "NativeAuth"
)

func PossibleValuesForAuthenticationMode() []string {
	return []string{
		string(AuthenticationModeMicrosoftEntraID),
		string(AuthenticationModeNativeAuth),
	}
}

func (s *AuthenticationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthenticationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthenticationMode(input string) (*AuthenticationMode, error) {
	vals := map[string]AuthenticationMode{
		"microsoftentraid": AuthenticationModeMicrosoftEntraID,
		"nativeauth":       AuthenticationModeNativeAuth,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationMode(input)
	return &out, nil
}

type CheckNameAvailabilityReason string

const (
	CheckNameAvailabilityReasonAlreadyExists CheckNameAvailabilityReason = "AlreadyExists"
	CheckNameAvailabilityReasonInvalid       CheckNameAvailabilityReason = "Invalid"
)

func PossibleValuesForCheckNameAvailabilityReason() []string {
	return []string{
		string(CheckNameAvailabilityReasonAlreadyExists),
		string(CheckNameAvailabilityReasonInvalid),
	}
}

func (s *CheckNameAvailabilityReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCheckNameAvailabilityReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCheckNameAvailabilityReason(input string) (*CheckNameAvailabilityReason, error) {
	vals := map[string]CheckNameAvailabilityReason{
		"alreadyexists": CheckNameAvailabilityReasonAlreadyExists,
		"invalid":       CheckNameAvailabilityReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CheckNameAvailabilityReason(input)
	return &out, nil
}

type CreateMode string

const (
	CreateModeDefault            CreateMode = "Default"
	CreateModeGeoReplica         CreateMode = "GeoReplica"
	CreateModePointInTimeRestore CreateMode = "PointInTimeRestore"
	CreateModeReplica            CreateMode = "Replica"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeDefault),
		string(CreateModeGeoReplica),
		string(CreateModePointInTimeRestore),
		string(CreateModeReplica),
	}
}

func (s *CreateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCreateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCreateMode(input string) (*CreateMode, error) {
	vals := map[string]CreateMode{
		"default":            CreateModeDefault,
		"georeplica":         CreateModeGeoReplica,
		"pointintimerestore": CreateModePointInTimeRestore,
		"replica":            CreateModeReplica,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateMode(input)
	return &out, nil
}

type DataApiMode string

const (
	DataApiModeDisabled DataApiMode = "Disabled"
	DataApiModeEnabled  DataApiMode = "Enabled"
)

func PossibleValuesForDataApiMode() []string {
	return []string{
		string(DataApiModeDisabled),
		string(DataApiModeEnabled),
	}
}

func (s *DataApiMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataApiMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataApiMode(input string) (*DataApiMode, error) {
	vals := map[string]DataApiMode{
		"disabled": DataApiModeDisabled,
		"enabled":  DataApiModeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataApiMode(input)
	return &out, nil
}

type HighAvailabilityMode string

const (
	HighAvailabilityModeDisabled               HighAvailabilityMode = "Disabled"
	HighAvailabilityModeSameZone               HighAvailabilityMode = "SameZone"
	HighAvailabilityModeZoneRedundantPreferred HighAvailabilityMode = "ZoneRedundantPreferred"
)

func PossibleValuesForHighAvailabilityMode() []string {
	return []string{
		string(HighAvailabilityModeDisabled),
		string(HighAvailabilityModeSameZone),
		string(HighAvailabilityModeZoneRedundantPreferred),
	}
}

func (s *HighAvailabilityMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHighAvailabilityMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHighAvailabilityMode(input string) (*HighAvailabilityMode, error) {
	vals := map[string]HighAvailabilityMode{
		"disabled":               HighAvailabilityModeDisabled,
		"samezone":               HighAvailabilityModeSameZone,
		"zoneredundantpreferred": HighAvailabilityModeZoneRedundantPreferred,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HighAvailabilityMode(input)
	return &out, nil
}

type KeyEncryptionKeyIdentityType string

const (
	KeyEncryptionKeyIdentityTypeUserAssignedIdentity KeyEncryptionKeyIdentityType = "UserAssignedIdentity"
)

func PossibleValuesForKeyEncryptionKeyIdentityType() []string {
	return []string{
		string(KeyEncryptionKeyIdentityTypeUserAssignedIdentity),
	}
}

func (s *KeyEncryptionKeyIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyEncryptionKeyIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyEncryptionKeyIdentityType(input string) (*KeyEncryptionKeyIdentityType, error) {
	vals := map[string]KeyEncryptionKeyIdentityType{
		"userassignedidentity": KeyEncryptionKeyIdentityTypeUserAssignedIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyEncryptionKeyIdentityType(input)
	return &out, nil
}

type MongoClusterStatus string

const (
	MongoClusterStatusDropping     MongoClusterStatus = "Dropping"
	MongoClusterStatusProvisioning MongoClusterStatus = "Provisioning"
	MongoClusterStatusReady        MongoClusterStatus = "Ready"
	MongoClusterStatusStarting     MongoClusterStatus = "Starting"
	MongoClusterStatusStopped      MongoClusterStatus = "Stopped"
	MongoClusterStatusStopping     MongoClusterStatus = "Stopping"
	MongoClusterStatusUpdating     MongoClusterStatus = "Updating"
)

func PossibleValuesForMongoClusterStatus() []string {
	return []string{
		string(MongoClusterStatusDropping),
		string(MongoClusterStatusProvisioning),
		string(MongoClusterStatusReady),
		string(MongoClusterStatusStarting),
		string(MongoClusterStatusStopped),
		string(MongoClusterStatusStopping),
		string(MongoClusterStatusUpdating),
	}
}

func (s *MongoClusterStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMongoClusterStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMongoClusterStatus(input string) (*MongoClusterStatus, error) {
	vals := map[string]MongoClusterStatus{
		"dropping":     MongoClusterStatusDropping,
		"provisioning": MongoClusterStatusProvisioning,
		"ready":        MongoClusterStatusReady,
		"starting":     MongoClusterStatusStarting,
		"stopped":      MongoClusterStatusStopped,
		"stopping":     MongoClusterStatusStopping,
		"updating":     MongoClusterStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MongoClusterStatus(input)
	return &out, nil
}

type PreviewFeature string

const (
	PreviewFeatureGeoReplicas PreviewFeature = "GeoReplicas"
)

func PossibleValuesForPreviewFeature() []string {
	return []string{
		string(PreviewFeatureGeoReplicas),
	}
}

func (s *PreviewFeature) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePreviewFeature(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePreviewFeature(input string) (*PreviewFeature, error) {
	vals := map[string]PreviewFeature{
		"georeplicas": PreviewFeatureGeoReplicas,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreviewFeature(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
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
		"creating":  PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusPending  PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		string(PrivateEndpointServiceConnectionStatusApproved),
		string(PrivateEndpointServiceConnectionStatusPending),
		string(PrivateEndpointServiceConnectionStatusRejected),
	}
}

func (s *PrivateEndpointServiceConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointServiceConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointServiceConnectionStatus(input string) (*PrivateEndpointServiceConnectionStatus, error) {
	vals := map[string]PrivateEndpointServiceConnectionStatus{
		"approved": PrivateEndpointServiceConnectionStatusApproved,
		"pending":  PrivateEndpointServiceConnectionStatusPending,
		"rejected": PrivateEndpointServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointServiceConnectionStatus(input)
	return &out, nil
}

type PromoteMode string

const (
	PromoteModeSwitchover PromoteMode = "Switchover"
)

func PossibleValuesForPromoteMode() []string {
	return []string{
		string(PromoteModeSwitchover),
	}
}

func (s *PromoteMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePromoteMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePromoteMode(input string) (*PromoteMode, error) {
	vals := map[string]PromoteMode{
		"switchover": PromoteModeSwitchover,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PromoteMode(input)
	return &out, nil
}

type PromoteOption string

const (
	PromoteOptionForced PromoteOption = "Forced"
)

func PossibleValuesForPromoteOption() []string {
	return []string{
		string(PromoteOptionForced),
	}
}

func (s *PromoteOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePromoteOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePromoteOption(input string) (*PromoteOption, error) {
	vals := map[string]PromoteOption{
		"forced": PromoteOptionForced,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PromoteOption(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled   ProvisioningState = "Canceled"
	ProvisioningStateDropping   ProvisioningState = "Dropping"
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateInProgress ProvisioningState = "InProgress"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
	ProvisioningStateUpdating   ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDropping),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
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
		"canceled":   ProvisioningStateCanceled,
		"dropping":   ProvisioningStateDropping,
		"failed":     ProvisioningStateFailed,
		"inprogress": ProvisioningStateInProgress,
		"succeeded":  ProvisioningStateSucceeded,
		"updating":   ProvisioningStateUpdating,
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

type ReplicationRole string

const (
	ReplicationRoleAsyncReplica    ReplicationRole = "AsyncReplica"
	ReplicationRoleGeoAsyncReplica ReplicationRole = "GeoAsyncReplica"
	ReplicationRolePrimary         ReplicationRole = "Primary"
)

func PossibleValuesForReplicationRole() []string {
	return []string{
		string(ReplicationRoleAsyncReplica),
		string(ReplicationRoleGeoAsyncReplica),
		string(ReplicationRolePrimary),
	}
}

func (s *ReplicationRole) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationRole(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationRole(input string) (*ReplicationRole, error) {
	vals := map[string]ReplicationRole{
		"asyncreplica":    ReplicationRoleAsyncReplica,
		"geoasyncreplica": ReplicationRoleGeoAsyncReplica,
		"primary":         ReplicationRolePrimary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationRole(input)
	return &out, nil
}

type ReplicationState string

const (
	ReplicationStateActive        ReplicationState = "Active"
	ReplicationStateBroken        ReplicationState = "Broken"
	ReplicationStateCatchup       ReplicationState = "Catchup"
	ReplicationStateProvisioning  ReplicationState = "Provisioning"
	ReplicationStateReconfiguring ReplicationState = "Reconfiguring"
	ReplicationStateUpdating      ReplicationState = "Updating"
)

func PossibleValuesForReplicationState() []string {
	return []string{
		string(ReplicationStateActive),
		string(ReplicationStateBroken),
		string(ReplicationStateCatchup),
		string(ReplicationStateProvisioning),
		string(ReplicationStateReconfiguring),
		string(ReplicationStateUpdating),
	}
}

func (s *ReplicationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationState(input string) (*ReplicationState, error) {
	vals := map[string]ReplicationState{
		"active":        ReplicationStateActive,
		"broken":        ReplicationStateBroken,
		"catchup":       ReplicationStateCatchup,
		"provisioning":  ReplicationStateProvisioning,
		"reconfiguring": ReplicationStateReconfiguring,
		"updating":      ReplicationStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationState(input)
	return &out, nil
}

type StorageType string

const (
	StorageTypePremiumSSD     StorageType = "PremiumSSD"
	StorageTypePremiumSSDvTwo StorageType = "PremiumSSDv2"
)

func PossibleValuesForStorageType() []string {
	return []string{
		string(StorageTypePremiumSSD),
		string(StorageTypePremiumSSDvTwo),
	}
}

func (s *StorageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageType(input string) (*StorageType, error) {
	vals := map[string]StorageType{
		"premiumssd":   StorageTypePremiumSSD,
		"premiumssdv2": StorageTypePremiumSSDvTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageType(input)
	return &out, nil
}
