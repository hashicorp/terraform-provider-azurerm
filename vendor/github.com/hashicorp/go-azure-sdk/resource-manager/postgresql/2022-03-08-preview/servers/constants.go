package servers

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArmServerKeyType string

const (
	ArmServerKeyTypeAzureKeyVault  ArmServerKeyType = "AzureKeyVault"
	ArmServerKeyTypeSystemAssigned ArmServerKeyType = "SystemAssigned"
)

func PossibleValuesForArmServerKeyType() []string {
	return []string{
		string(ArmServerKeyTypeAzureKeyVault),
		string(ArmServerKeyTypeSystemAssigned),
	}
}

func parseArmServerKeyType(input string) (*ArmServerKeyType, error) {
	vals := map[string]ArmServerKeyType{
		"azurekeyvault":  ArmServerKeyTypeAzureKeyVault,
		"systemassigned": ArmServerKeyTypeSystemAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ArmServerKeyType(input)
	return &out, nil
}

type CreateMode string

const (
	CreateModeCreate             CreateMode = "Create"
	CreateModeDefault            CreateMode = "Default"
	CreateModeGeoRestore         CreateMode = "GeoRestore"
	CreateModePointInTimeRestore CreateMode = "PointInTimeRestore"
	CreateModeReplica            CreateMode = "Replica"
	CreateModeUpdate             CreateMode = "Update"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeCreate),
		string(CreateModeDefault),
		string(CreateModeGeoRestore),
		string(CreateModePointInTimeRestore),
		string(CreateModeReplica),
		string(CreateModeUpdate),
	}
}

func parseCreateMode(input string) (*CreateMode, error) {
	vals := map[string]CreateMode{
		"create":             CreateModeCreate,
		"default":            CreateModeDefault,
		"georestore":         CreateModeGeoRestore,
		"pointintimerestore": CreateModePointInTimeRestore,
		"replica":            CreateModeReplica,
		"update":             CreateModeUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateMode(input)
	return &out, nil
}

type CreateModeForUpdate string

const (
	CreateModeForUpdateDefault CreateModeForUpdate = "Default"
	CreateModeForUpdateUpdate  CreateModeForUpdate = "Update"
)

func PossibleValuesForCreateModeForUpdate() []string {
	return []string{
		string(CreateModeForUpdateDefault),
		string(CreateModeForUpdateUpdate),
	}
}

func parseCreateModeForUpdate(input string) (*CreateModeForUpdate, error) {
	vals := map[string]CreateModeForUpdate{
		"default": CreateModeForUpdateDefault,
		"update":  CreateModeForUpdateUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateModeForUpdate(input)
	return &out, nil
}

type GeoRedundantBackupEnum string

const (
	GeoRedundantBackupEnumDisabled GeoRedundantBackupEnum = "Disabled"
	GeoRedundantBackupEnumEnabled  GeoRedundantBackupEnum = "Enabled"
)

func PossibleValuesForGeoRedundantBackupEnum() []string {
	return []string{
		string(GeoRedundantBackupEnumDisabled),
		string(GeoRedundantBackupEnumEnabled),
	}
}

func parseGeoRedundantBackupEnum(input string) (*GeoRedundantBackupEnum, error) {
	vals := map[string]GeoRedundantBackupEnum{
		"disabled": GeoRedundantBackupEnumDisabled,
		"enabled":  GeoRedundantBackupEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GeoRedundantBackupEnum(input)
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

type IdentityType string

const (
	IdentityTypeNone           IdentityType = "None"
	IdentityTypeSystemAssigned IdentityType = "SystemAssigned"
	IdentityTypeUserAssigned   IdentityType = "UserAssigned"
)

func PossibleValuesForIdentityType() []string {
	return []string{
		string(IdentityTypeNone),
		string(IdentityTypeSystemAssigned),
		string(IdentityTypeUserAssigned),
	}
}

func parseIdentityType(input string) (*IdentityType, error) {
	vals := map[string]IdentityType{
		"none":           IdentityTypeNone,
		"systemassigned": IdentityTypeSystemAssigned,
		"userassigned":   IdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityType(input)
	return &out, nil
}

type ReplicationRole string

const (
	ReplicationRoleAsyncReplica    ReplicationRole = "AsyncReplica"
	ReplicationRoleGeoAsyncReplica ReplicationRole = "GeoAsyncReplica"
	ReplicationRoleGeoSyncReplica  ReplicationRole = "GeoSyncReplica"
	ReplicationRolePrimary         ReplicationRole = "Primary"
	ReplicationRoleSecondary       ReplicationRole = "Secondary"
	ReplicationRoleSyncReplica     ReplicationRole = "SyncReplica"
	ReplicationRoleWalReplica      ReplicationRole = "WalReplica"
)

func PossibleValuesForReplicationRole() []string {
	return []string{
		string(ReplicationRoleAsyncReplica),
		string(ReplicationRoleGeoAsyncReplica),
		string(ReplicationRoleGeoSyncReplica),
		string(ReplicationRolePrimary),
		string(ReplicationRoleSecondary),
		string(ReplicationRoleSyncReplica),
		string(ReplicationRoleWalReplica),
	}
}

func parseReplicationRole(input string) (*ReplicationRole, error) {
	vals := map[string]ReplicationRole{
		"asyncreplica":    ReplicationRoleAsyncReplica,
		"geoasyncreplica": ReplicationRoleGeoAsyncReplica,
		"geosyncreplica":  ReplicationRoleGeoSyncReplica,
		"primary":         ReplicationRolePrimary,
		"secondary":       ReplicationRoleSecondary,
		"syncreplica":     ReplicationRoleSyncReplica,
		"walreplica":      ReplicationRoleWalReplica,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationRole(input)
	return &out, nil
}

type ServerHAState string

const (
	ServerHAStateCreatingStandby ServerHAState = "CreatingStandby"
	ServerHAStateFailingOver     ServerHAState = "FailingOver"
	ServerHAStateHealthy         ServerHAState = "Healthy"
	ServerHAStateNotEnabled      ServerHAState = "NotEnabled"
	ServerHAStateRemovingStandby ServerHAState = "RemovingStandby"
	ServerHAStateReplicatingData ServerHAState = "ReplicatingData"
)

func PossibleValuesForServerHAState() []string {
	return []string{
		string(ServerHAStateCreatingStandby),
		string(ServerHAStateFailingOver),
		string(ServerHAStateHealthy),
		string(ServerHAStateNotEnabled),
		string(ServerHAStateRemovingStandby),
		string(ServerHAStateReplicatingData),
	}
}

func parseServerHAState(input string) (*ServerHAState, error) {
	vals := map[string]ServerHAState{
		"creatingstandby": ServerHAStateCreatingStandby,
		"failingover":     ServerHAStateFailingOver,
		"healthy":         ServerHAStateHealthy,
		"notenabled":      ServerHAStateNotEnabled,
		"removingstandby": ServerHAStateRemovingStandby,
		"replicatingdata": ServerHAStateReplicatingData,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerHAState(input)
	return &out, nil
}

type ServerPublicNetworkAccessState string

const (
	ServerPublicNetworkAccessStateDisabled ServerPublicNetworkAccessState = "Disabled"
	ServerPublicNetworkAccessStateEnabled  ServerPublicNetworkAccessState = "Enabled"
)

func PossibleValuesForServerPublicNetworkAccessState() []string {
	return []string{
		string(ServerPublicNetworkAccessStateDisabled),
		string(ServerPublicNetworkAccessStateEnabled),
	}
}

func parseServerPublicNetworkAccessState(input string) (*ServerPublicNetworkAccessState, error) {
	vals := map[string]ServerPublicNetworkAccessState{
		"disabled": ServerPublicNetworkAccessStateDisabled,
		"enabled":  ServerPublicNetworkAccessStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerPublicNetworkAccessState(input)
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
	ServerVersionOneFour  ServerVersion = "14"
	ServerVersionOneOne   ServerVersion = "11"
	ServerVersionOneThree ServerVersion = "13"
	ServerVersionOneTwo   ServerVersion = "12"
)

func PossibleValuesForServerVersion() []string {
	return []string{
		string(ServerVersionOneFour),
		string(ServerVersionOneOne),
		string(ServerVersionOneThree),
		string(ServerVersionOneTwo),
	}
}

func parseServerVersion(input string) (*ServerVersion, error) {
	vals := map[string]ServerVersion{
		"14": ServerVersionOneFour,
		"11": ServerVersionOneOne,
		"13": ServerVersionOneThree,
		"12": ServerVersionOneTwo,
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
