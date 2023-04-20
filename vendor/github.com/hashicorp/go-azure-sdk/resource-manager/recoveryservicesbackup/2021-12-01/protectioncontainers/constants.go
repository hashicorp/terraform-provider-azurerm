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

type ContainerType string

const (
	ContainerTypeAzureBackupServerContainer                  ContainerType = "AzureBackupServerContainer"
	ContainerTypeAzureSqlContainer                           ContainerType = "AzureSqlContainer"
	ContainerTypeAzureWorkloadContainer                      ContainerType = "AzureWorkloadContainer"
	ContainerTypeCluster                                     ContainerType = "Cluster"
	ContainerTypeDPMContainer                                ContainerType = "DPMContainer"
	ContainerTypeGenericContainer                            ContainerType = "GenericContainer"
	ContainerTypeIaasVMContainer                             ContainerType = "IaasVMContainer"
	ContainerTypeIaasVMServiceContainer                      ContainerType = "IaasVMServiceContainer"
	ContainerTypeInvalid                                     ContainerType = "Invalid"
	ContainerTypeMABContainer                                ContainerType = "MABContainer"
	ContainerTypeMicrosoftPointClassicComputeVirtualMachines ContainerType = "Microsoft.ClassicCompute/virtualMachines"
	ContainerTypeMicrosoftPointComputeVirtualMachines        ContainerType = "Microsoft.Compute/virtualMachines"
	ContainerTypeSQLAGWorkLoadContainer                      ContainerType = "SQLAGWorkLoadContainer"
	ContainerTypeStorageContainer                            ContainerType = "StorageContainer"
	ContainerTypeUnknown                                     ContainerType = "Unknown"
	ContainerTypeVCenter                                     ContainerType = "VCenter"
	ContainerTypeVMAppContainer                              ContainerType = "VMAppContainer"
	ContainerTypeWindows                                     ContainerType = "Windows"
)

func PossibleValuesForContainerType() []string {
	return []string{
		string(ContainerTypeAzureBackupServerContainer),
		string(ContainerTypeAzureSqlContainer),
		string(ContainerTypeAzureWorkloadContainer),
		string(ContainerTypeCluster),
		string(ContainerTypeDPMContainer),
		string(ContainerTypeGenericContainer),
		string(ContainerTypeIaasVMContainer),
		string(ContainerTypeIaasVMServiceContainer),
		string(ContainerTypeInvalid),
		string(ContainerTypeMABContainer),
		string(ContainerTypeMicrosoftPointClassicComputeVirtualMachines),
		string(ContainerTypeMicrosoftPointComputeVirtualMachines),
		string(ContainerTypeSQLAGWorkLoadContainer),
		string(ContainerTypeStorageContainer),
		string(ContainerTypeUnknown),
		string(ContainerTypeVCenter),
		string(ContainerTypeVMAppContainer),
		string(ContainerTypeWindows),
	}
}

func parseContainerType(input string) (*ContainerType, error) {
	vals := map[string]ContainerType{
		"azurebackupservercontainer": ContainerTypeAzureBackupServerContainer,
		"azuresqlcontainer":          ContainerTypeAzureSqlContainer,
		"azureworkloadcontainer":     ContainerTypeAzureWorkloadContainer,
		"cluster":                    ContainerTypeCluster,
		"dpmcontainer":               ContainerTypeDPMContainer,
		"genericcontainer":           ContainerTypeGenericContainer,
		"iaasvmcontainer":            ContainerTypeIaasVMContainer,
		"iaasvmservicecontainer":     ContainerTypeIaasVMServiceContainer,
		"invalid":                    ContainerTypeInvalid,
		"mabcontainer":               ContainerTypeMABContainer,
		"microsoft.classiccompute/virtualmachines": ContainerTypeMicrosoftPointClassicComputeVirtualMachines,
		"microsoft.compute/virtualmachines":        ContainerTypeMicrosoftPointComputeVirtualMachines,
		"sqlagworkloadcontainer":                   ContainerTypeSQLAGWorkLoadContainer,
		"storagecontainer":                         ContainerTypeStorageContainer,
		"unknown":                                  ContainerTypeUnknown,
		"vcenter":                                  ContainerTypeVCenter,
		"vmappcontainer":                           ContainerTypeVMAppContainer,
		"windows":                                  ContainerTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerType(input)
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
