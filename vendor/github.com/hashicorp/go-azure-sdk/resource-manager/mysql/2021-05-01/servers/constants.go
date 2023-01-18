package servers

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateMode string

const (
	CreateModeDefault            CreateMode = "Default"
	CreateModeGeoRestore         CreateMode = "GeoRestore"
	CreateModePointInTimeRestore CreateMode = "PointInTimeRestore"
	CreateModeReplica            CreateMode = "Replica"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeDefault),
		string(CreateModeGeoRestore),
		string(CreateModePointInTimeRestore),
		string(CreateModeReplica),
	}
}

func parseCreateMode(input string) (*CreateMode, error) {
	vals := map[string]CreateMode{
		"default":            CreateModeDefault,
		"georestore":         CreateModeGeoRestore,
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

type DataEncryptionType string

const (
	DataEncryptionTypeAzureKeyVault DataEncryptionType = "AzureKeyVault"
	DataEncryptionTypeSystemManaged DataEncryptionType = "SystemManaged"
)

func PossibleValuesForDataEncryptionType() []string {
	return []string{
		string(DataEncryptionTypeAzureKeyVault),
		string(DataEncryptionTypeSystemManaged),
	}
}

func parseDataEncryptionType(input string) (*DataEncryptionType, error) {
	vals := map[string]DataEncryptionType{
		"azurekeyvault": DataEncryptionTypeAzureKeyVault,
		"systemmanaged": DataEncryptionTypeSystemManaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataEncryptionType(input)
	return &out, nil
}

type EnableStatusEnum string

const (
	EnableStatusEnumDisabled EnableStatusEnum = "Disabled"
	EnableStatusEnumEnabled  EnableStatusEnum = "Enabled"
)

func PossibleValuesForEnableStatusEnum() []string {
	return []string{
		string(EnableStatusEnumDisabled),
		string(EnableStatusEnumEnabled),
	}
}

func parseEnableStatusEnum(input string) (*EnableStatusEnum, error) {
	vals := map[string]EnableStatusEnum{
		"disabled": EnableStatusEnumDisabled,
		"enabled":  EnableStatusEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnableStatusEnum(input)
	return &out, nil
}

type HighAvailabilityMode string

const (
	HighAvailabilityModeDisabled      HighAvailabilityMode = "Disabled"
	HighAvailabilityModeSameZone      HighAvailabilityMode = "SameZone"
	HighAvailabilityModeZoneRedundant HighAvailabilityMode = "ZoneRedundant"
)

func PossibleValuesForHighAvailabilityMode() []string {
	return []string{
		string(HighAvailabilityModeDisabled),
		string(HighAvailabilityModeSameZone),
		string(HighAvailabilityModeZoneRedundant),
	}
}

func parseHighAvailabilityMode(input string) (*HighAvailabilityMode, error) {
	vals := map[string]HighAvailabilityMode{
		"disabled":      HighAvailabilityModeDisabled,
		"samezone":      HighAvailabilityModeSameZone,
		"zoneredundant": HighAvailabilityModeZoneRedundant,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HighAvailabilityMode(input)
	return &out, nil
}

type HighAvailabilityState string

const (
	HighAvailabilityStateCreatingStandby HighAvailabilityState = "CreatingStandby"
	HighAvailabilityStateFailingOver     HighAvailabilityState = "FailingOver"
	HighAvailabilityStateHealthy         HighAvailabilityState = "Healthy"
	HighAvailabilityStateNotEnabled      HighAvailabilityState = "NotEnabled"
	HighAvailabilityStateRemovingStandby HighAvailabilityState = "RemovingStandby"
)

func PossibleValuesForHighAvailabilityState() []string {
	return []string{
		string(HighAvailabilityStateCreatingStandby),
		string(HighAvailabilityStateFailingOver),
		string(HighAvailabilityStateHealthy),
		string(HighAvailabilityStateNotEnabled),
		string(HighAvailabilityStateRemovingStandby),
	}
}

func parseHighAvailabilityState(input string) (*HighAvailabilityState, error) {
	vals := map[string]HighAvailabilityState{
		"creatingstandby": HighAvailabilityStateCreatingStandby,
		"failingover":     HighAvailabilityStateFailingOver,
		"healthy":         HighAvailabilityStateHealthy,
		"notenabled":      HighAvailabilityStateNotEnabled,
		"removingstandby": HighAvailabilityStateRemovingStandby,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HighAvailabilityState(input)
	return &out, nil
}

type ManagedServiceIdentityType string

const (
	ManagedServiceIdentityTypeUserAssigned ManagedServiceIdentityType = "UserAssigned"
)

func PossibleValuesForManagedServiceIdentityType() []string {
	return []string{
		string(ManagedServiceIdentityTypeUserAssigned),
	}
}

func parseManagedServiceIdentityType(input string) (*ManagedServiceIdentityType, error) {
	vals := map[string]ManagedServiceIdentityType{
		"userassigned": ManagedServiceIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedServiceIdentityType(input)
	return &out, nil
}

type ReplicationRole string

const (
	ReplicationRoleNone    ReplicationRole = "None"
	ReplicationRoleReplica ReplicationRole = "Replica"
	ReplicationRoleSource  ReplicationRole = "Source"
)

func PossibleValuesForReplicationRole() []string {
	return []string{
		string(ReplicationRoleNone),
		string(ReplicationRoleReplica),
		string(ReplicationRoleSource),
	}
}

func parseReplicationRole(input string) (*ReplicationRole, error) {
	vals := map[string]ReplicationRole{
		"none":    ReplicationRoleNone,
		"replica": ReplicationRoleReplica,
		"source":  ReplicationRoleSource,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationRole(input)
	return &out, nil
}

type ServerState string

const (
	ServerStateDisabled ServerState = "Disabled"
	ServerStateDropping ServerState = "Dropping"
	ServerStateReady    ServerState = "Ready"
	ServerStateStarting ServerState = "Starting"
	ServerStateStopped  ServerState = "Stopped"
	ServerStateStopping ServerState = "Stopping"
	ServerStateUpdating ServerState = "Updating"
)

func PossibleValuesForServerState() []string {
	return []string{
		string(ServerStateDisabled),
		string(ServerStateDropping),
		string(ServerStateReady),
		string(ServerStateStarting),
		string(ServerStateStopped),
		string(ServerStateStopping),
		string(ServerStateUpdating),
	}
}

func parseServerState(input string) (*ServerState, error) {
	vals := map[string]ServerState{
		"disabled": ServerStateDisabled,
		"dropping": ServerStateDropping,
		"ready":    ServerStateReady,
		"starting": ServerStateStarting,
		"stopped":  ServerStateStopped,
		"stopping": ServerStateStopping,
		"updating": ServerStateUpdating,
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
	ServerVersionEightPointZeroPointTwoOne ServerVersion = "8.0.21"
	ServerVersionFivePointSeven            ServerVersion = "5.7"
)

func PossibleValuesForServerVersion() []string {
	return []string{
		string(ServerVersionEightPointZeroPointTwoOne),
		string(ServerVersionFivePointSeven),
	}
}

func parseServerVersion(input string) (*ServerVersion, error) {
	vals := map[string]ServerVersion{
		"8.0.21": ServerVersionEightPointZeroPointTwoOne,
		"5.7":    ServerVersionFivePointSeven,
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
	SkuTierBurstable       SkuTier = "Burstable"
	SkuTierGeneralPurpose  SkuTier = "GeneralPurpose"
	SkuTierMemoryOptimized SkuTier = "MemoryOptimized"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBurstable),
		string(SkuTierGeneralPurpose),
		string(SkuTierMemoryOptimized),
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"burstable":       SkuTierBurstable,
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
