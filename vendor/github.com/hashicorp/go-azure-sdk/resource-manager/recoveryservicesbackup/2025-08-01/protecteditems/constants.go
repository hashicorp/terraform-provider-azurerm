package protecteditems

import (
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupManagementType string

const (
	BackupManagementTypeAzureBackupServer BackupManagementType = "AzureBackupServer"
	BackupManagementTypeAzureIaasVM       BackupManagementType = "AzureIaasVM"
	BackupManagementTypeAzureSql          BackupManagementType = "AzureSql"
	BackupManagementTypeAzureStorage      BackupManagementType = "AzureStorage"
	BackupManagementTypeAzureWorkload     BackupManagementType = "AzureWorkload"
	BackupManagementTypeDPM               BackupManagementType = "DPM"
	BackupManagementTypeDefaultBackup     BackupManagementType = "DefaultBackup"
	BackupManagementTypeInvalid           BackupManagementType = "Invalid"
	BackupManagementTypeMAB               BackupManagementType = "MAB"
)

func PossibleValuesForBackupManagementType() []string {
	return []string{
		string(BackupManagementTypeAzureBackupServer),
		string(BackupManagementTypeAzureIaasVM),
		string(BackupManagementTypeAzureSql),
		string(BackupManagementTypeAzureStorage),
		string(BackupManagementTypeAzureWorkload),
		string(BackupManagementTypeDPM),
		string(BackupManagementTypeDefaultBackup),
		string(BackupManagementTypeInvalid),
		string(BackupManagementTypeMAB),
	}
}

func parseBackupManagementType(input string) (*BackupManagementType, error) {
	vals := map[string]BackupManagementType{
		"azurebackupserver": BackupManagementTypeAzureBackupServer,
		"azureiaasvm":       BackupManagementTypeAzureIaasVM,
		"azuresql":          BackupManagementTypeAzureSql,
		"azurestorage":      BackupManagementTypeAzureStorage,
		"azureworkload":     BackupManagementTypeAzureWorkload,
		"dpm":               BackupManagementTypeDPM,
		"defaultbackup":     BackupManagementTypeDefaultBackup,
		"invalid":           BackupManagementTypeInvalid,
		"mab":               BackupManagementTypeMAB,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupManagementType(input)
	return &out, nil
}

type CreateMode string

const (
	CreateModeDefault CreateMode = "Default"
	CreateModeInvalid CreateMode = "Invalid"
	CreateModeRecover CreateMode = "Recover"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeDefault),
		string(CreateModeInvalid),
		string(CreateModeRecover),
	}
}

func parseCreateMode(input string) (*CreateMode, error) {
	vals := map[string]CreateMode{
		"default": CreateModeDefault,
		"invalid": CreateModeInvalid,
		"recover": CreateModeRecover,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateMode(input)
	return &out, nil
}

type DataSourceType string

const (
	DataSourceTypeAzureFileShare    DataSourceType = "AzureFileShare"
	DataSourceTypeAzureSqlDb        DataSourceType = "AzureSqlDb"
	DataSourceTypeClient            DataSourceType = "Client"
	DataSourceTypeExchange          DataSourceType = "Exchange"
	DataSourceTypeFileFolder        DataSourceType = "FileFolder"
	DataSourceTypeGenericDataSource DataSourceType = "GenericDataSource"
	DataSourceTypeInvalid           DataSourceType = "Invalid"
	DataSourceTypeSAPAseDatabase    DataSourceType = "SAPAseDatabase"
	DataSourceTypeSAPHanaDBInstance DataSourceType = "SAPHanaDBInstance"
	DataSourceTypeSAPHanaDatabase   DataSourceType = "SAPHanaDatabase"
	DataSourceTypeSQLDB             DataSourceType = "SQLDB"
	DataSourceTypeSQLDataBase       DataSourceType = "SQLDataBase"
	DataSourceTypeSharepoint        DataSourceType = "Sharepoint"
	DataSourceTypeSystemState       DataSourceType = "SystemState"
	DataSourceTypeVM                DataSourceType = "VM"
	DataSourceTypeVMwareVM          DataSourceType = "VMwareVM"
)

func PossibleValuesForDataSourceType() []string {
	return []string{
		string(DataSourceTypeAzureFileShare),
		string(DataSourceTypeAzureSqlDb),
		string(DataSourceTypeClient),
		string(DataSourceTypeExchange),
		string(DataSourceTypeFileFolder),
		string(DataSourceTypeGenericDataSource),
		string(DataSourceTypeInvalid),
		string(DataSourceTypeSAPAseDatabase),
		string(DataSourceTypeSAPHanaDBInstance),
		string(DataSourceTypeSAPHanaDatabase),
		string(DataSourceTypeSQLDB),
		string(DataSourceTypeSQLDataBase),
		string(DataSourceTypeSharepoint),
		string(DataSourceTypeSystemState),
		string(DataSourceTypeVM),
		string(DataSourceTypeVMwareVM),
	}
}

func parseDataSourceType(input string) (*DataSourceType, error) {
	vals := map[string]DataSourceType{
		"azurefileshare":    DataSourceTypeAzureFileShare,
		"azuresqldb":        DataSourceTypeAzureSqlDb,
		"client":            DataSourceTypeClient,
		"exchange":          DataSourceTypeExchange,
		"filefolder":        DataSourceTypeFileFolder,
		"genericdatasource": DataSourceTypeGenericDataSource,
		"invalid":           DataSourceTypeInvalid,
		"sapasedatabase":    DataSourceTypeSAPAseDatabase,
		"saphanadbinstance": DataSourceTypeSAPHanaDBInstance,
		"saphanadatabase":   DataSourceTypeSAPHanaDatabase,
		"sqldb":             DataSourceTypeSQLDB,
		"sqldatabase":       DataSourceTypeSQLDataBase,
		"sharepoint":        DataSourceTypeSharepoint,
		"systemstate":       DataSourceTypeSystemState,
		"vm":                DataSourceTypeVM,
		"vmwarevm":          DataSourceTypeVMwareVM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataSourceType(input)
	return &out, nil
}

type HealthStatus string

const (
	HealthStatusActionRequired  HealthStatus = "ActionRequired"
	HealthStatusActionSuggested HealthStatus = "ActionSuggested"
	HealthStatusInvalid         HealthStatus = "Invalid"
	HealthStatusPassed          HealthStatus = "Passed"
)

func PossibleValuesForHealthStatus() []string {
	return []string{
		string(HealthStatusActionRequired),
		string(HealthStatusActionSuggested),
		string(HealthStatusInvalid),
		string(HealthStatusPassed),
	}
}

func parseHealthStatus(input string) (*HealthStatus, error) {
	vals := map[string]HealthStatus{
		"actionrequired":  HealthStatusActionRequired,
		"actionsuggested": HealthStatusActionSuggested,
		"invalid":         HealthStatusInvalid,
		"passed":          HealthStatusPassed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthStatus(input)
	return &out, nil
}

type LastBackupStatus string

const (
	LastBackupStatusHealthy   LastBackupStatus = "Healthy"
	LastBackupStatusIRPending LastBackupStatus = "IRPending"
	LastBackupStatusInvalid   LastBackupStatus = "Invalid"
	LastBackupStatusUnhealthy LastBackupStatus = "Unhealthy"
)

func PossibleValuesForLastBackupStatus() []string {
	return []string{
		string(LastBackupStatusHealthy),
		string(LastBackupStatusIRPending),
		string(LastBackupStatusInvalid),
		string(LastBackupStatusUnhealthy),
	}
}

func parseLastBackupStatus(input string) (*LastBackupStatus, error) {
	vals := map[string]LastBackupStatus{
		"healthy":   LastBackupStatusHealthy,
		"irpending": LastBackupStatusIRPending,
		"invalid":   LastBackupStatusInvalid,
		"unhealthy": LastBackupStatusUnhealthy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LastBackupStatus(input)
	return &out, nil
}

type ProtectedItemHealthStatus string

const (
	ProtectedItemHealthStatusHealthy      ProtectedItemHealthStatus = "Healthy"
	ProtectedItemHealthStatusIRPending    ProtectedItemHealthStatus = "IRPending"
	ProtectedItemHealthStatusInvalid      ProtectedItemHealthStatus = "Invalid"
	ProtectedItemHealthStatusNotReachable ProtectedItemHealthStatus = "NotReachable"
	ProtectedItemHealthStatusUnhealthy    ProtectedItemHealthStatus = "Unhealthy"
)

func PossibleValuesForProtectedItemHealthStatus() []string {
	return []string{
		string(ProtectedItemHealthStatusHealthy),
		string(ProtectedItemHealthStatusIRPending),
		string(ProtectedItemHealthStatusInvalid),
		string(ProtectedItemHealthStatusNotReachable),
		string(ProtectedItemHealthStatusUnhealthy),
	}
}

func parseProtectedItemHealthStatus(input string) (*ProtectedItemHealthStatus, error) {
	vals := map[string]ProtectedItemHealthStatus{
		"healthy":      ProtectedItemHealthStatusHealthy,
		"irpending":    ProtectedItemHealthStatusIRPending,
		"invalid":      ProtectedItemHealthStatusInvalid,
		"notreachable": ProtectedItemHealthStatusNotReachable,
		"unhealthy":    ProtectedItemHealthStatusUnhealthy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProtectedItemHealthStatus(input)
	return &out, nil
}

type ProtectedItemState string

const (
	ProtectedItemStateBackupsSuspended  ProtectedItemState = "BackupsSuspended"
	ProtectedItemStateIRPending         ProtectedItemState = "IRPending"
	ProtectedItemStateInvalid           ProtectedItemState = "Invalid"
	ProtectedItemStateProtected         ProtectedItemState = "Protected"
	ProtectedItemStateProtectionError   ProtectedItemState = "ProtectionError"
	ProtectedItemStateProtectionPaused  ProtectedItemState = "ProtectionPaused"
	ProtectedItemStateProtectionStopped ProtectedItemState = "ProtectionStopped"
)

func PossibleValuesForProtectedItemState() []string {
	return []string{
		string(ProtectedItemStateBackupsSuspended),
		string(ProtectedItemStateIRPending),
		string(ProtectedItemStateInvalid),
		string(ProtectedItemStateProtected),
		string(ProtectedItemStateProtectionError),
		string(ProtectedItemStateProtectionPaused),
		string(ProtectedItemStateProtectionStopped),
	}
}

func parseProtectedItemState(input string) (*ProtectedItemState, error) {
	vals := map[string]ProtectedItemState{
		"backupssuspended":  ProtectedItemStateBackupsSuspended,
		"irpending":         ProtectedItemStateIRPending,
		"invalid":           ProtectedItemStateInvalid,
		"protected":         ProtectedItemStateProtected,
		"protectionerror":   ProtectedItemStateProtectionError,
		"protectionpaused":  ProtectedItemStateProtectionPaused,
		"protectionstopped": ProtectedItemStateProtectionStopped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProtectedItemState(input)
	return &out, nil
}

type ProtectionState string

const (
	ProtectionStateBackupsSuspended  ProtectionState = "BackupsSuspended"
	ProtectionStateIRPending         ProtectionState = "IRPending"
	ProtectionStateInvalid           ProtectionState = "Invalid"
	ProtectionStateProtected         ProtectionState = "Protected"
	ProtectionStateProtectionError   ProtectionState = "ProtectionError"
	ProtectionStateProtectionPaused  ProtectionState = "ProtectionPaused"
	ProtectionStateProtectionStopped ProtectionState = "ProtectionStopped"
)

func PossibleValuesForProtectionState() []string {
	return []string{
		string(ProtectionStateBackupsSuspended),
		string(ProtectionStateIRPending),
		string(ProtectionStateInvalid),
		string(ProtectionStateProtected),
		string(ProtectionStateProtectionError),
		string(ProtectionStateProtectionPaused),
		string(ProtectionStateProtectionStopped),
	}
}

func parseProtectionState(input string) (*ProtectionState, error) {
	vals := map[string]ProtectionState{
		"backupssuspended":  ProtectionStateBackupsSuspended,
		"irpending":         ProtectionStateIRPending,
		"invalid":           ProtectionStateInvalid,
		"protected":         ProtectionStateProtected,
		"protectionerror":   ProtectionStateProtectionError,
		"protectionpaused":  ProtectionStateProtectionPaused,
		"protectionstopped": ProtectionStateProtectionStopped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProtectionState(input)
	return &out, nil
}

type ResourceHealthStatus string

const (
	ResourceHealthStatusHealthy             ResourceHealthStatus = "Healthy"
	ResourceHealthStatusInvalid             ResourceHealthStatus = "Invalid"
	ResourceHealthStatusPersistentDegraded  ResourceHealthStatus = "PersistentDegraded"
	ResourceHealthStatusPersistentUnhealthy ResourceHealthStatus = "PersistentUnhealthy"
	ResourceHealthStatusTransientDegraded   ResourceHealthStatus = "TransientDegraded"
	ResourceHealthStatusTransientUnhealthy  ResourceHealthStatus = "TransientUnhealthy"
)

func PossibleValuesForResourceHealthStatus() []string {
	return []string{
		string(ResourceHealthStatusHealthy),
		string(ResourceHealthStatusInvalid),
		string(ResourceHealthStatusPersistentDegraded),
		string(ResourceHealthStatusPersistentUnhealthy),
		string(ResourceHealthStatusTransientDegraded),
		string(ResourceHealthStatusTransientUnhealthy),
	}
}

func parseResourceHealthStatus(input string) (*ResourceHealthStatus, error) {
	vals := map[string]ResourceHealthStatus{
		"healthy":             ResourceHealthStatusHealthy,
		"invalid":             ResourceHealthStatusInvalid,
		"persistentdegraded":  ResourceHealthStatusPersistentDegraded,
		"persistentunhealthy": ResourceHealthStatusPersistentUnhealthy,
		"transientdegraded":   ResourceHealthStatusTransientDegraded,
		"transientunhealthy":  ResourceHealthStatusTransientUnhealthy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceHealthStatus(input)
	return &out, nil
}
