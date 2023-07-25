package protectioncontainers

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AcquireStorageAccountLock string

const (
	AcquireStorageAccountLockAcquire    AcquireStorageAccountLock = "Acquire"
	AcquireStorageAccountLockNotAcquire AcquireStorageAccountLock = "NotAcquire"
)

func PossibleValuesForAcquireStorageAccountLock() []string {
	return []string{
		string(AcquireStorageAccountLockAcquire),
		string(AcquireStorageAccountLockNotAcquire),
	}
}

func parseAcquireStorageAccountLock(input string) (*AcquireStorageAccountLock, error) {
	vals := map[string]AcquireStorageAccountLock{
		"acquire":    AcquireStorageAccountLockAcquire,
		"notacquire": AcquireStorageAccountLockNotAcquire,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AcquireStorageAccountLock(input)
	return &out, nil
}

type BackupItemType string

const (
	BackupItemTypeAzureFileShare    BackupItemType = "AzureFileShare"
	BackupItemTypeAzureSqlDb        BackupItemType = "AzureSqlDb"
	BackupItemTypeClient            BackupItemType = "Client"
	BackupItemTypeExchange          BackupItemType = "Exchange"
	BackupItemTypeFileFolder        BackupItemType = "FileFolder"
	BackupItemTypeGenericDataSource BackupItemType = "GenericDataSource"
	BackupItemTypeInvalid           BackupItemType = "Invalid"
	BackupItemTypeSAPAseDatabase    BackupItemType = "SAPAseDatabase"
	BackupItemTypeSAPHanaDBInstance BackupItemType = "SAPHanaDBInstance"
	BackupItemTypeSAPHanaDatabase   BackupItemType = "SAPHanaDatabase"
	BackupItemTypeSQLDB             BackupItemType = "SQLDB"
	BackupItemTypeSQLDataBase       BackupItemType = "SQLDataBase"
	BackupItemTypeSharepoint        BackupItemType = "Sharepoint"
	BackupItemTypeSystemState       BackupItemType = "SystemState"
	BackupItemTypeVM                BackupItemType = "VM"
	BackupItemTypeVMwareVM          BackupItemType = "VMwareVM"
)

func PossibleValuesForBackupItemType() []string {
	return []string{
		string(BackupItemTypeAzureFileShare),
		string(BackupItemTypeAzureSqlDb),
		string(BackupItemTypeClient),
		string(BackupItemTypeExchange),
		string(BackupItemTypeFileFolder),
		string(BackupItemTypeGenericDataSource),
		string(BackupItemTypeInvalid),
		string(BackupItemTypeSAPAseDatabase),
		string(BackupItemTypeSAPHanaDBInstance),
		string(BackupItemTypeSAPHanaDatabase),
		string(BackupItemTypeSQLDB),
		string(BackupItemTypeSQLDataBase),
		string(BackupItemTypeSharepoint),
		string(BackupItemTypeSystemState),
		string(BackupItemTypeVM),
		string(BackupItemTypeVMwareVM),
	}
}

func parseBackupItemType(input string) (*BackupItemType, error) {
	vals := map[string]BackupItemType{
		"azurefileshare":    BackupItemTypeAzureFileShare,
		"azuresqldb":        BackupItemTypeAzureSqlDb,
		"client":            BackupItemTypeClient,
		"exchange":          BackupItemTypeExchange,
		"filefolder":        BackupItemTypeFileFolder,
		"genericdatasource": BackupItemTypeGenericDataSource,
		"invalid":           BackupItemTypeInvalid,
		"sapasedatabase":    BackupItemTypeSAPAseDatabase,
		"saphanadbinstance": BackupItemTypeSAPHanaDBInstance,
		"saphanadatabase":   BackupItemTypeSAPHanaDatabase,
		"sqldb":             BackupItemTypeSQLDB,
		"sqldatabase":       BackupItemTypeSQLDataBase,
		"sharepoint":        BackupItemTypeSharepoint,
		"systemstate":       BackupItemTypeSystemState,
		"vm":                BackupItemTypeVM,
		"vmwarevm":          BackupItemTypeVMwareVM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupItemType(input)
	return &out, nil
}

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

type OperationType string

const (
	OperationTypeInvalid    OperationType = "Invalid"
	OperationTypeRegister   OperationType = "Register"
	OperationTypeReregister OperationType = "Reregister"
)

func PossibleValuesForOperationType() []string {
	return []string{
		string(OperationTypeInvalid),
		string(OperationTypeRegister),
		string(OperationTypeReregister),
	}
}

func parseOperationType(input string) (*OperationType, error) {
	vals := map[string]OperationType{
		"invalid":    OperationTypeInvalid,
		"register":   OperationTypeRegister,
		"reregister": OperationTypeReregister,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationType(input)
	return &out, nil
}

type ProtectableContainerType string

const (
	ProtectableContainerTypeAzureBackupServerContainer                  ProtectableContainerType = "AzureBackupServerContainer"
	ProtectableContainerTypeAzureSqlContainer                           ProtectableContainerType = "AzureSqlContainer"
	ProtectableContainerTypeAzureWorkloadContainer                      ProtectableContainerType = "AzureWorkloadContainer"
	ProtectableContainerTypeCluster                                     ProtectableContainerType = "Cluster"
	ProtectableContainerTypeDPMContainer                                ProtectableContainerType = "DPMContainer"
	ProtectableContainerTypeGenericContainer                            ProtectableContainerType = "GenericContainer"
	ProtectableContainerTypeIaasVMContainer                             ProtectableContainerType = "IaasVMContainer"
	ProtectableContainerTypeIaasVMServiceContainer                      ProtectableContainerType = "IaasVMServiceContainer"
	ProtectableContainerTypeInvalid                                     ProtectableContainerType = "Invalid"
	ProtectableContainerTypeMABContainer                                ProtectableContainerType = "MABContainer"
	ProtectableContainerTypeMicrosoftPointClassicComputeVirtualMachines ProtectableContainerType = "Microsoft.ClassicCompute/virtualMachines"
	ProtectableContainerTypeMicrosoftPointComputeVirtualMachines        ProtectableContainerType = "Microsoft.Compute/virtualMachines"
	ProtectableContainerTypeSQLAGWorkLoadContainer                      ProtectableContainerType = "SQLAGWorkLoadContainer"
	ProtectableContainerTypeStorageContainer                            ProtectableContainerType = "StorageContainer"
	ProtectableContainerTypeUnknown                                     ProtectableContainerType = "Unknown"
	ProtectableContainerTypeVCenter                                     ProtectableContainerType = "VCenter"
	ProtectableContainerTypeVMAppContainer                              ProtectableContainerType = "VMAppContainer"
	ProtectableContainerTypeWindows                                     ProtectableContainerType = "Windows"
)

func PossibleValuesForProtectableContainerType() []string {
	return []string{
		string(ProtectableContainerTypeAzureBackupServerContainer),
		string(ProtectableContainerTypeAzureSqlContainer),
		string(ProtectableContainerTypeAzureWorkloadContainer),
		string(ProtectableContainerTypeCluster),
		string(ProtectableContainerTypeDPMContainer),
		string(ProtectableContainerTypeGenericContainer),
		string(ProtectableContainerTypeIaasVMContainer),
		string(ProtectableContainerTypeIaasVMServiceContainer),
		string(ProtectableContainerTypeInvalid),
		string(ProtectableContainerTypeMABContainer),
		string(ProtectableContainerTypeMicrosoftPointClassicComputeVirtualMachines),
		string(ProtectableContainerTypeMicrosoftPointComputeVirtualMachines),
		string(ProtectableContainerTypeSQLAGWorkLoadContainer),
		string(ProtectableContainerTypeStorageContainer),
		string(ProtectableContainerTypeUnknown),
		string(ProtectableContainerTypeVCenter),
		string(ProtectableContainerTypeVMAppContainer),
		string(ProtectableContainerTypeWindows),
	}
}

func parseProtectableContainerType(input string) (*ProtectableContainerType, error) {
	vals := map[string]ProtectableContainerType{
		"azurebackupservercontainer": ProtectableContainerTypeAzureBackupServerContainer,
		"azuresqlcontainer":          ProtectableContainerTypeAzureSqlContainer,
		"azureworkloadcontainer":     ProtectableContainerTypeAzureWorkloadContainer,
		"cluster":                    ProtectableContainerTypeCluster,
		"dpmcontainer":               ProtectableContainerTypeDPMContainer,
		"genericcontainer":           ProtectableContainerTypeGenericContainer,
		"iaasvmcontainer":            ProtectableContainerTypeIaasVMContainer,
		"iaasvmservicecontainer":     ProtectableContainerTypeIaasVMServiceContainer,
		"invalid":                    ProtectableContainerTypeInvalid,
		"mabcontainer":               ProtectableContainerTypeMABContainer,
		"microsoft.classiccompute/virtualmachines": ProtectableContainerTypeMicrosoftPointClassicComputeVirtualMachines,
		"microsoft.compute/virtualmachines":        ProtectableContainerTypeMicrosoftPointComputeVirtualMachines,
		"sqlagworkloadcontainer":                   ProtectableContainerTypeSQLAGWorkLoadContainer,
		"storagecontainer":                         ProtectableContainerTypeStorageContainer,
		"unknown":                                  ProtectableContainerTypeUnknown,
		"vcenter":                                  ProtectableContainerTypeVCenter,
		"vmappcontainer":                           ProtectableContainerTypeVMAppContainer,
		"windows":                                  ProtectableContainerTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProtectableContainerType(input)
	return &out, nil
}

type WorkloadType string

const (
	WorkloadTypeAzureFileShare    WorkloadType = "AzureFileShare"
	WorkloadTypeAzureSqlDb        WorkloadType = "AzureSqlDb"
	WorkloadTypeClient            WorkloadType = "Client"
	WorkloadTypeExchange          WorkloadType = "Exchange"
	WorkloadTypeFileFolder        WorkloadType = "FileFolder"
	WorkloadTypeGenericDataSource WorkloadType = "GenericDataSource"
	WorkloadTypeInvalid           WorkloadType = "Invalid"
	WorkloadTypeSAPAseDatabase    WorkloadType = "SAPAseDatabase"
	WorkloadTypeSAPHanaDBInstance WorkloadType = "SAPHanaDBInstance"
	WorkloadTypeSAPHanaDatabase   WorkloadType = "SAPHanaDatabase"
	WorkloadTypeSQLDB             WorkloadType = "SQLDB"
	WorkloadTypeSQLDataBase       WorkloadType = "SQLDataBase"
	WorkloadTypeSharepoint        WorkloadType = "Sharepoint"
	WorkloadTypeSystemState       WorkloadType = "SystemState"
	WorkloadTypeVM                WorkloadType = "VM"
	WorkloadTypeVMwareVM          WorkloadType = "VMwareVM"
)

func PossibleValuesForWorkloadType() []string {
	return []string{
		string(WorkloadTypeAzureFileShare),
		string(WorkloadTypeAzureSqlDb),
		string(WorkloadTypeClient),
		string(WorkloadTypeExchange),
		string(WorkloadTypeFileFolder),
		string(WorkloadTypeGenericDataSource),
		string(WorkloadTypeInvalid),
		string(WorkloadTypeSAPAseDatabase),
		string(WorkloadTypeSAPHanaDBInstance),
		string(WorkloadTypeSAPHanaDatabase),
		string(WorkloadTypeSQLDB),
		string(WorkloadTypeSQLDataBase),
		string(WorkloadTypeSharepoint),
		string(WorkloadTypeSystemState),
		string(WorkloadTypeVM),
		string(WorkloadTypeVMwareVM),
	}
}

func parseWorkloadType(input string) (*WorkloadType, error) {
	vals := map[string]WorkloadType{
		"azurefileshare":    WorkloadTypeAzureFileShare,
		"azuresqldb":        WorkloadTypeAzureSqlDb,
		"client":            WorkloadTypeClient,
		"exchange":          WorkloadTypeExchange,
		"filefolder":        WorkloadTypeFileFolder,
		"genericdatasource": WorkloadTypeGenericDataSource,
		"invalid":           WorkloadTypeInvalid,
		"sapasedatabase":    WorkloadTypeSAPAseDatabase,
		"saphanadbinstance": WorkloadTypeSAPHanaDBInstance,
		"saphanadatabase":   WorkloadTypeSAPHanaDatabase,
		"sqldb":             WorkloadTypeSQLDB,
		"sqldatabase":       WorkloadTypeSQLDataBase,
		"sharepoint":        WorkloadTypeSharepoint,
		"systemstate":       WorkloadTypeSystemState,
		"vm":                WorkloadTypeVM,
		"vmwarevm":          WorkloadTypeVMwareVM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkloadType(input)
	return &out, nil
}
