package replicas

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GeoRedundantBackup string

const (
	GeoRedundantBackupDisabled GeoRedundantBackup = "Disabled"
	GeoRedundantBackupEnabled  GeoRedundantBackup = "Enabled"
)

func PossibleValuesForGeoRedundantBackup() []string {
	return []string{
		string(GeoRedundantBackupDisabled),
		string(GeoRedundantBackupEnabled),
	}
}

func (s *GeoRedundantBackup) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGeoRedundantBackup(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGeoRedundantBackup(input string) (*GeoRedundantBackup, error) {
	vals := map[string]GeoRedundantBackup{
		"disabled": GeoRedundantBackupDisabled,
		"enabled":  GeoRedundantBackupEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GeoRedundantBackup(input)
	return &out, nil
}

type InfrastructureEncryption string

const (
	InfrastructureEncryptionDisabled InfrastructureEncryption = "Disabled"
	InfrastructureEncryptionEnabled  InfrastructureEncryption = "Enabled"
)

func PossibleValuesForInfrastructureEncryption() []string {
	return []string{
		string(InfrastructureEncryptionDisabled),
		string(InfrastructureEncryptionEnabled),
	}
}

func (s *InfrastructureEncryption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInfrastructureEncryption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInfrastructureEncryption(input string) (*InfrastructureEncryption, error) {
	vals := map[string]InfrastructureEncryption{
		"disabled": InfrastructureEncryptionDisabled,
		"enabled":  InfrastructureEncryptionEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InfrastructureEncryption(input)
	return &out, nil
}

type MinimalTlsVersionEnum string

const (
	MinimalTlsVersionEnumTLSEnforcementDisabled MinimalTlsVersionEnum = "TLSEnforcementDisabled"
	MinimalTlsVersionEnumTLSOneOne              MinimalTlsVersionEnum = "TLS1_1"
	MinimalTlsVersionEnumTLSOneTwo              MinimalTlsVersionEnum = "TLS1_2"
	MinimalTlsVersionEnumTLSOneZero             MinimalTlsVersionEnum = "TLS1_0"
)

func PossibleValuesForMinimalTlsVersionEnum() []string {
	return []string{
		string(MinimalTlsVersionEnumTLSEnforcementDisabled),
		string(MinimalTlsVersionEnumTLSOneOne),
		string(MinimalTlsVersionEnumTLSOneTwo),
		string(MinimalTlsVersionEnumTLSOneZero),
	}
}

func (s *MinimalTlsVersionEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMinimalTlsVersionEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMinimalTlsVersionEnum(input string) (*MinimalTlsVersionEnum, error) {
	vals := map[string]MinimalTlsVersionEnum{
		"tlsenforcementdisabled": MinimalTlsVersionEnumTLSEnforcementDisabled,
		"tls1_1":                 MinimalTlsVersionEnumTLSOneOne,
		"tls1_2":                 MinimalTlsVersionEnumTLSOneTwo,
		"tls1_0":                 MinimalTlsVersionEnumTLSOneZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MinimalTlsVersionEnum(input)
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

type PublicNetworkAccessEnum string

const (
	PublicNetworkAccessEnumDisabled PublicNetworkAccessEnum = "Disabled"
	PublicNetworkAccessEnumEnabled  PublicNetworkAccessEnum = "Enabled"
)

func PossibleValuesForPublicNetworkAccessEnum() []string {
	return []string{
		string(PublicNetworkAccessEnumDisabled),
		string(PublicNetworkAccessEnumEnabled),
	}
}

func (s *PublicNetworkAccessEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccessEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccessEnum(input string) (*PublicNetworkAccessEnum, error) {
	vals := map[string]PublicNetworkAccessEnum{
		"disabled": PublicNetworkAccessEnumDisabled,
		"enabled":  PublicNetworkAccessEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccessEnum(input)
	return &out, nil
}

type ServerState string

const (
	ServerStateDisabled     ServerState = "Disabled"
	ServerStateDropping     ServerState = "Dropping"
	ServerStateInaccessible ServerState = "Inaccessible"
	ServerStateReady        ServerState = "Ready"
)

func PossibleValuesForServerState() []string {
	return []string{
		string(ServerStateDisabled),
		string(ServerStateDropping),
		string(ServerStateInaccessible),
		string(ServerStateReady),
	}
}

func (s *ServerState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerState(input string) (*ServerState, error) {
	vals := map[string]ServerState{
		"disabled":     ServerStateDisabled,
		"dropping":     ServerStateDropping,
		"inaccessible": ServerStateInaccessible,
		"ready":        ServerStateReady,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerState(input)
	return &out, nil
}

type ServerVersion string

const (
	ServerVersionEightPointZero ServerVersion = "8.0"
	ServerVersionFivePointSeven ServerVersion = "5.7"
	ServerVersionFivePointSix   ServerVersion = "5.6"
)

func PossibleValuesForServerVersion() []string {
	return []string{
		string(ServerVersionEightPointZero),
		string(ServerVersionFivePointSeven),
		string(ServerVersionFivePointSix),
	}
}

func (s *ServerVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerVersion(input string) (*ServerVersion, error) {
	vals := map[string]ServerVersion{
		"8.0": ServerVersionEightPointZero,
		"5.7": ServerVersionFivePointSeven,
		"5.6": ServerVersionFivePointSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerVersion(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic           SkuTier = "Basic"
	SkuTierGeneralPurpose  SkuTier = "GeneralPurpose"
	SkuTierMemoryOptimized SkuTier = "MemoryOptimized"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierGeneralPurpose),
		string(SkuTierMemoryOptimized),
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
		"basic":           SkuTierBasic,
		"generalpurpose":  SkuTierGeneralPurpose,
		"memoryoptimized": SkuTierMemoryOptimized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type SslEnforcementEnum string

const (
	SslEnforcementEnumDisabled SslEnforcementEnum = "Disabled"
	SslEnforcementEnumEnabled  SslEnforcementEnum = "Enabled"
)

func PossibleValuesForSslEnforcementEnum() []string {
	return []string{
		string(SslEnforcementEnumDisabled),
		string(SslEnforcementEnumEnabled),
	}
}

func (s *SslEnforcementEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSslEnforcementEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSslEnforcementEnum(input string) (*SslEnforcementEnum, error) {
	vals := map[string]SslEnforcementEnum{
		"disabled": SslEnforcementEnumDisabled,
		"enabled":  SslEnforcementEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SslEnforcementEnum(input)
	return &out, nil
}

type StorageAutogrow string

const (
	StorageAutogrowDisabled StorageAutogrow = "Disabled"
	StorageAutogrowEnabled  StorageAutogrow = "Enabled"
)

func PossibleValuesForStorageAutogrow() []string {
	return []string{
		string(StorageAutogrowDisabled),
		string(StorageAutogrowEnabled),
	}
}

func (s *StorageAutogrow) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageAutogrow(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageAutogrow(input string) (*StorageAutogrow, error) {
	vals := map[string]StorageAutogrow{
		"disabled": StorageAutogrowDisabled,
		"enabled":  StorageAutogrowEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAutogrow(input)
	return &out, nil
}
