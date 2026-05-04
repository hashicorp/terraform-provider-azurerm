package servers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdministratorType string

const (
	AdministratorTypeActiveDirectory AdministratorType = "ActiveDirectory"
)

func PossibleValuesForAdministratorType() []string {
	return []string{
		string(AdministratorTypeActiveDirectory),
	}
}

func (s *AdministratorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdministratorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdministratorType(input string) (*AdministratorType, error) {
	vals := map[string]AdministratorType{
		"activedirectory": AdministratorTypeActiveDirectory,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdministratorType(input)
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

type CheckNameAvailabilityResourceType string

const (
	CheckNameAvailabilityResourceTypeMicrosoftPointSqlServers CheckNameAvailabilityResourceType = "Microsoft.Sql/servers"
)

func PossibleValuesForCheckNameAvailabilityResourceType() []string {
	return []string{
		string(CheckNameAvailabilityResourceTypeMicrosoftPointSqlServers),
	}
}

func (s *CheckNameAvailabilityResourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCheckNameAvailabilityResourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCheckNameAvailabilityResourceType(input string) (*CheckNameAvailabilityResourceType, error) {
	vals := map[string]CheckNameAvailabilityResourceType{
		"microsoft.sql/servers": CheckNameAvailabilityResourceTypeMicrosoftPointSqlServers,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CheckNameAvailabilityResourceType(input)
	return &out, nil
}

type ExternalGovernanceStatus string

const (
	ExternalGovernanceStatusDisabled ExternalGovernanceStatus = "Disabled"
	ExternalGovernanceStatusEnabled  ExternalGovernanceStatus = "Enabled"
)

func PossibleValuesForExternalGovernanceStatus() []string {
	return []string{
		string(ExternalGovernanceStatusDisabled),
		string(ExternalGovernanceStatusEnabled),
	}
}

func (s *ExternalGovernanceStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExternalGovernanceStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExternalGovernanceStatus(input string) (*ExternalGovernanceStatus, error) {
	vals := map[string]ExternalGovernanceStatus{
		"disabled": ExternalGovernanceStatusDisabled,
		"enabled":  ExternalGovernanceStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExternalGovernanceStatus(input)
	return &out, nil
}

type MinimalTlsVersion string

const (
	MinimalTlsVersionNone          MinimalTlsVersion = "None"
	MinimalTlsVersionOnePointOne   MinimalTlsVersion = "1.1"
	MinimalTlsVersionOnePointThree MinimalTlsVersion = "1.3"
	MinimalTlsVersionOnePointTwo   MinimalTlsVersion = "1.2"
	MinimalTlsVersionOnePointZero  MinimalTlsVersion = "1.0"
)

func PossibleValuesForMinimalTlsVersion() []string {
	return []string{
		string(MinimalTlsVersionNone),
		string(MinimalTlsVersionOnePointOne),
		string(MinimalTlsVersionOnePointThree),
		string(MinimalTlsVersionOnePointTwo),
		string(MinimalTlsVersionOnePointZero),
	}
}

func (s *MinimalTlsVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMinimalTlsVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMinimalTlsVersion(input string) (*MinimalTlsVersion, error) {
	vals := map[string]MinimalTlsVersion{
		"none": MinimalTlsVersionNone,
		"1.1":  MinimalTlsVersionOnePointOne,
		"1.3":  MinimalTlsVersionOnePointThree,
		"1.2":  MinimalTlsVersionOnePointTwo,
		"1.0":  MinimalTlsVersionOnePointZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MinimalTlsVersion(input)
	return &out, nil
}

type PrincipalType string

const (
	PrincipalTypeApplication PrincipalType = "Application"
	PrincipalTypeGroup       PrincipalType = "Group"
	PrincipalTypeUser        PrincipalType = "User"
)

func PossibleValuesForPrincipalType() []string {
	return []string{
		string(PrincipalTypeApplication),
		string(PrincipalTypeGroup),
		string(PrincipalTypeUser),
	}
}

func (s *PrincipalType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrincipalType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrincipalType(input string) (*PrincipalType, error) {
	vals := map[string]PrincipalType{
		"application": PrincipalTypeApplication,
		"group":       PrincipalTypeGroup,
		"user":        PrincipalTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrincipalType(input)
	return &out, nil
}

type PrivateEndpointProvisioningState string

const (
	PrivateEndpointProvisioningStateApproving PrivateEndpointProvisioningState = "Approving"
	PrivateEndpointProvisioningStateDropping  PrivateEndpointProvisioningState = "Dropping"
	PrivateEndpointProvisioningStateFailed    PrivateEndpointProvisioningState = "Failed"
	PrivateEndpointProvisioningStateReady     PrivateEndpointProvisioningState = "Ready"
	PrivateEndpointProvisioningStateRejecting PrivateEndpointProvisioningState = "Rejecting"
)

func PossibleValuesForPrivateEndpointProvisioningState() []string {
	return []string{
		string(PrivateEndpointProvisioningStateApproving),
		string(PrivateEndpointProvisioningStateDropping),
		string(PrivateEndpointProvisioningStateFailed),
		string(PrivateEndpointProvisioningStateReady),
		string(PrivateEndpointProvisioningStateRejecting),
	}
}

func (s *PrivateEndpointProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateEndpointProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateEndpointProvisioningState(input string) (*PrivateEndpointProvisioningState, error) {
	vals := map[string]PrivateEndpointProvisioningState{
		"approving": PrivateEndpointProvisioningStateApproving,
		"dropping":  PrivateEndpointProvisioningStateDropping,
		"failed":    PrivateEndpointProvisioningStateFailed,
		"ready":     PrivateEndpointProvisioningStateReady,
		"rejecting": PrivateEndpointProvisioningStateRejecting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointProvisioningState(input)
	return &out, nil
}

type PrivateLinkServiceConnectionStateActionsRequire string

const (
	PrivateLinkServiceConnectionStateActionsRequireNone PrivateLinkServiceConnectionStateActionsRequire = "None"
)

func PossibleValuesForPrivateLinkServiceConnectionStateActionsRequire() []string {
	return []string{
		string(PrivateLinkServiceConnectionStateActionsRequireNone),
	}
}

func (s *PrivateLinkServiceConnectionStateActionsRequire) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateLinkServiceConnectionStateActionsRequire(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateLinkServiceConnectionStateActionsRequire(input string) (*PrivateLinkServiceConnectionStateActionsRequire, error) {
	vals := map[string]PrivateLinkServiceConnectionStateActionsRequire{
		"none": PrivateLinkServiceConnectionStateActionsRequireNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateLinkServiceConnectionStateActionsRequire(input)
	return &out, nil
}

type PrivateLinkServiceConnectionStateStatus string

const (
	PrivateLinkServiceConnectionStateStatusApproved     PrivateLinkServiceConnectionStateStatus = "Approved"
	PrivateLinkServiceConnectionStateStatusDisconnected PrivateLinkServiceConnectionStateStatus = "Disconnected"
	PrivateLinkServiceConnectionStateStatusPending      PrivateLinkServiceConnectionStateStatus = "Pending"
	PrivateLinkServiceConnectionStateStatusRejected     PrivateLinkServiceConnectionStateStatus = "Rejected"
)

func PossibleValuesForPrivateLinkServiceConnectionStateStatus() []string {
	return []string{
		string(PrivateLinkServiceConnectionStateStatusApproved),
		string(PrivateLinkServiceConnectionStateStatusDisconnected),
		string(PrivateLinkServiceConnectionStateStatusPending),
		string(PrivateLinkServiceConnectionStateStatusRejected),
	}
}

func (s *PrivateLinkServiceConnectionStateStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateLinkServiceConnectionStateStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateLinkServiceConnectionStateStatus(input string) (*PrivateLinkServiceConnectionStateStatus, error) {
	vals := map[string]PrivateLinkServiceConnectionStateStatus{
		"approved":     PrivateLinkServiceConnectionStateStatusApproved,
		"disconnected": PrivateLinkServiceConnectionStateStatusDisconnected,
		"pending":      PrivateLinkServiceConnectionStateStatusPending,
		"rejected":     PrivateLinkServiceConnectionStateStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateLinkServiceConnectionStateStatus(input)
	return &out, nil
}

type ServerNetworkAccessFlag string

const (
	ServerNetworkAccessFlagDisabled ServerNetworkAccessFlag = "Disabled"
	ServerNetworkAccessFlagEnabled  ServerNetworkAccessFlag = "Enabled"
)

func PossibleValuesForServerNetworkAccessFlag() []string {
	return []string{
		string(ServerNetworkAccessFlagDisabled),
		string(ServerNetworkAccessFlagEnabled),
	}
}

func (s *ServerNetworkAccessFlag) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerNetworkAccessFlag(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerNetworkAccessFlag(input string) (*ServerNetworkAccessFlag, error) {
	vals := map[string]ServerNetworkAccessFlag{
		"disabled": ServerNetworkAccessFlagDisabled,
		"enabled":  ServerNetworkAccessFlagEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerNetworkAccessFlag(input)
	return &out, nil
}

type ServerPublicNetworkAccessFlag string

const (
	ServerPublicNetworkAccessFlagDisabled           ServerPublicNetworkAccessFlag = "Disabled"
	ServerPublicNetworkAccessFlagEnabled            ServerPublicNetworkAccessFlag = "Enabled"
	ServerPublicNetworkAccessFlagSecuredByPerimeter ServerPublicNetworkAccessFlag = "SecuredByPerimeter"
)

func PossibleValuesForServerPublicNetworkAccessFlag() []string {
	return []string{
		string(ServerPublicNetworkAccessFlagDisabled),
		string(ServerPublicNetworkAccessFlagEnabled),
		string(ServerPublicNetworkAccessFlagSecuredByPerimeter),
	}
}

func (s *ServerPublicNetworkAccessFlag) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerPublicNetworkAccessFlag(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerPublicNetworkAccessFlag(input string) (*ServerPublicNetworkAccessFlag, error) {
	vals := map[string]ServerPublicNetworkAccessFlag{
		"disabled":           ServerPublicNetworkAccessFlagDisabled,
		"enabled":            ServerPublicNetworkAccessFlagEnabled,
		"securedbyperimeter": ServerPublicNetworkAccessFlagSecuredByPerimeter,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerPublicNetworkAccessFlag(input)
	return &out, nil
}

type ServerWorkspaceFeature string

const (
	ServerWorkspaceFeatureConnected    ServerWorkspaceFeature = "Connected"
	ServerWorkspaceFeatureDisconnected ServerWorkspaceFeature = "Disconnected"
)

func PossibleValuesForServerWorkspaceFeature() []string {
	return []string{
		string(ServerWorkspaceFeatureConnected),
		string(ServerWorkspaceFeatureDisconnected),
	}
}

func (s *ServerWorkspaceFeature) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerWorkspaceFeature(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerWorkspaceFeature(input string) (*ServerWorkspaceFeature, error) {
	vals := map[string]ServerWorkspaceFeature{
		"connected":    ServerWorkspaceFeatureConnected,
		"disconnected": ServerWorkspaceFeatureDisconnected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerWorkspaceFeature(input)
	return &out, nil
}

type StorageKeyType string

const (
	StorageKeyTypeSharedAccessKey  StorageKeyType = "SharedAccessKey"
	StorageKeyTypeStorageAccessKey StorageKeyType = "StorageAccessKey"
)

func PossibleValuesForStorageKeyType() []string {
	return []string{
		string(StorageKeyTypeSharedAccessKey),
		string(StorageKeyTypeStorageAccessKey),
	}
}

func (s *StorageKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageKeyType(input string) (*StorageKeyType, error) {
	vals := map[string]StorageKeyType{
		"sharedaccesskey":  StorageKeyTypeSharedAccessKey,
		"storageaccesskey": StorageKeyTypeStorageAccessKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageKeyType(input)
	return &out, nil
}
