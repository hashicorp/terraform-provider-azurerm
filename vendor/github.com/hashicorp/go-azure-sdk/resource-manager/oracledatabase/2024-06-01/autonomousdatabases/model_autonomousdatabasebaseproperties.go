package autonomousdatabases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseBaseProperties interface {
	AutonomousDatabaseBaseProperties() BaseAutonomousDatabaseBasePropertiesImpl
}

var _ AutonomousDatabaseBaseProperties = BaseAutonomousDatabaseBasePropertiesImpl{}

type BaseAutonomousDatabaseBasePropertiesImpl struct {
	ActualUsedDataStorageSizeInTbs           *float64                           `json:"actualUsedDataStorageSizeInTbs,omitempty"`
	AdminPassword                            *string                            `json:"adminPassword,omitempty"`
	AllocatedStorageSizeInTbs                *float64                           `json:"allocatedStorageSizeInTbs,omitempty"`
	ApexDetails                              *ApexDetailsType                   `json:"apexDetails,omitempty"`
	AutonomousDatabaseId                     *string                            `json:"autonomousDatabaseId,omitempty"`
	AutonomousMaintenanceScheduleType        *AutonomousMaintenanceScheduleType `json:"autonomousMaintenanceScheduleType,omitempty"`
	AvailableUpgradeVersions                 *[]string                          `json:"availableUpgradeVersions,omitempty"`
	BackupRetentionPeriodInDays              *int64                             `json:"backupRetentionPeriodInDays,omitempty"`
	CharacterSet                             *string                            `json:"characterSet,omitempty"`
	ComputeCount                             *float64                           `json:"computeCount,omitempty"`
	ComputeModel                             *ComputeModel                      `json:"computeModel,omitempty"`
	ConnectionStrings                        *ConnectionStringType              `json:"connectionStrings,omitempty"`
	ConnectionURLs                           *ConnectionURLType                 `json:"connectionUrls,omitempty"`
	CpuCoreCount                             *int64                             `json:"cpuCoreCount,omitempty"`
	CustomerContacts                         *[]CustomerContact                 `json:"customerContacts,omitempty"`
	DataBaseType                             DataBaseType                       `json:"dataBaseType"`
	DataSafeStatus                           *DataSafeStatusType                `json:"dataSafeStatus,omitempty"`
	DataStorageSizeInGbs                     *int64                             `json:"dataStorageSizeInGbs,omitempty"`
	DataStorageSizeInTbs                     *int64                             `json:"dataStorageSizeInTbs,omitempty"`
	DatabaseEdition                          *DatabaseEditionType               `json:"databaseEdition,omitempty"`
	DbVersion                                *string                            `json:"dbVersion,omitempty"`
	DbWorkload                               *WorkloadType                      `json:"dbWorkload,omitempty"`
	DisplayName                              *string                            `json:"displayName,omitempty"`
	FailedDataRecoveryInSeconds              *int64                             `json:"failedDataRecoveryInSeconds,omitempty"`
	InMemoryAreaInGbs                        *int64                             `json:"inMemoryAreaInGbs,omitempty"`
	IsAutoScalingEnabled                     *bool                              `json:"isAutoScalingEnabled,omitempty"`
	IsAutoScalingForStorageEnabled           *bool                              `json:"isAutoScalingForStorageEnabled,omitempty"`
	IsLocalDataGuardEnabled                  *bool                              `json:"isLocalDataGuardEnabled,omitempty"`
	IsMtlsConnectionRequired                 *bool                              `json:"isMtlsConnectionRequired,omitempty"`
	IsPreview                                *bool                              `json:"isPreview,omitempty"`
	IsPreviewVersionWithServiceTermsAccepted *bool                              `json:"isPreviewVersionWithServiceTermsAccepted,omitempty"`
	IsRemoteDataGuardEnabled                 *bool                              `json:"isRemoteDataGuardEnabled,omitempty"`
	LicenseModel                             *LicenseModel                      `json:"licenseModel,omitempty"`
	LifecycleDetails                         *string                            `json:"lifecycleDetails,omitempty"`
	LifecycleState                           *AutonomousDatabaseLifecycleState  `json:"lifecycleState,omitempty"`
	LocalAdgAutoFailoverMaxDataLossLimit     *int64                             `json:"localAdgAutoFailoverMaxDataLossLimit,omitempty"`
	LocalDisasterRecoveryType                *DisasterRecoveryType              `json:"localDisasterRecoveryType,omitempty"`
	LocalStandbyDb                           *AutonomousDatabaseStandbySummary  `json:"localStandbyDb,omitempty"`
	LongTermBackupSchedule                   *LongTermBackUpScheduleDetails     `json:"longTermBackupSchedule,omitempty"`
	MemoryPerOracleComputeUnitInGbs          *int64                             `json:"memoryPerOracleComputeUnitInGbs,omitempty"`
	NcharacterSet                            *string                            `json:"ncharacterSet,omitempty"`
	NextLongTermBackupTimeStamp              *string                            `json:"nextLongTermBackupTimeStamp,omitempty"`
	OciURL                                   *string                            `json:"ociUrl,omitempty"`
	Ocid                                     *string                            `json:"ocid,omitempty"`
	OpenMode                                 *OpenModeType                      `json:"openMode,omitempty"`
	OperationsInsightsStatus                 *OperationsInsightsStatusType      `json:"operationsInsightsStatus,omitempty"`
	PeerDbId                                 *string                            `json:"peerDbId,omitempty"`
	PeerDbIds                                *[]string                          `json:"peerDbIds,omitempty"`
	PermissionLevel                          *PermissionLevelType               `json:"permissionLevel,omitempty"`
	PrivateEndpoint                          *string                            `json:"privateEndpoint,omitempty"`
	PrivateEndpointIP                        *string                            `json:"privateEndpointIp,omitempty"`
	PrivateEndpointLabel                     *string                            `json:"privateEndpointLabel,omitempty"`
	ProvisionableCPUs                        *[]int64                           `json:"provisionableCpus,omitempty"`
	ProvisioningState                        *AzureResourceProvisioningState    `json:"provisioningState,omitempty"`
	Role                                     *RoleType                          `json:"role,omitempty"`
	ScheduledOperations                      *ScheduledOperationsType           `json:"scheduledOperations,omitempty"`
	ServiceConsoleURL                        *string                            `json:"serviceConsoleUrl,omitempty"`
	SqlWebDeveloperURL                       *string                            `json:"sqlWebDeveloperUrl,omitempty"`
	SubnetId                                 *string                            `json:"subnetId,omitempty"`
	SupportedRegionsToCloneTo                *[]string                          `json:"supportedRegionsToCloneTo,omitempty"`
	TimeCreated                              *string                            `json:"timeCreated,omitempty"`
	TimeDataGuardRoleChanged                 *string                            `json:"timeDataGuardRoleChanged,omitempty"`
	TimeDeletionOfFreeAutonomousDatabase     *string                            `json:"timeDeletionOfFreeAutonomousDatabase,omitempty"`
	TimeLocalDataGuardEnabled                *string                            `json:"timeLocalDataGuardEnabled,omitempty"`
	TimeMaintenanceBegin                     *string                            `json:"timeMaintenanceBegin,omitempty"`
	TimeMaintenanceEnd                       *string                            `json:"timeMaintenanceEnd,omitempty"`
	TimeOfLastFailover                       *string                            `json:"timeOfLastFailover,omitempty"`
	TimeOfLastRefresh                        *string                            `json:"timeOfLastRefresh,omitempty"`
	TimeOfLastRefreshPoint                   *string                            `json:"timeOfLastRefreshPoint,omitempty"`
	TimeOfLastSwitchover                     *string                            `json:"timeOfLastSwitchover,omitempty"`
	TimeReclamationOfFreeAutonomousDatabase  *string                            `json:"timeReclamationOfFreeAutonomousDatabase,omitempty"`
	UsedDataStorageSizeInGbs                 *int64                             `json:"usedDataStorageSizeInGbs,omitempty"`
	UsedDataStorageSizeInTbs                 *int64                             `json:"usedDataStorageSizeInTbs,omitempty"`
	VnetId                                   *string                            `json:"vnetId,omitempty"`
	WhitelistedIPs                           *[]string                          `json:"whitelistedIps,omitempty"`
}

func (s BaseAutonomousDatabaseBasePropertiesImpl) AutonomousDatabaseBaseProperties() BaseAutonomousDatabaseBasePropertiesImpl {
	return s
}

var _ AutonomousDatabaseBaseProperties = RawAutonomousDatabaseBasePropertiesImpl{}

// RawAutonomousDatabaseBasePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAutonomousDatabaseBasePropertiesImpl struct {
	autonomousDatabaseBaseProperties BaseAutonomousDatabaseBasePropertiesImpl
	Type                             string
	Values                           map[string]interface{}
}

func (s RawAutonomousDatabaseBasePropertiesImpl) AutonomousDatabaseBaseProperties() BaseAutonomousDatabaseBasePropertiesImpl {
	return s.autonomousDatabaseBaseProperties
}

func UnmarshalAutonomousDatabaseBasePropertiesImplementation(input []byte) (AutonomousDatabaseBaseProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AutonomousDatabaseBaseProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["dataBaseType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Clone") {
		var out AutonomousDatabaseCloneProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutonomousDatabaseCloneProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Regular") {
		var out AutonomousDatabaseProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutonomousDatabaseProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseAutonomousDatabaseBasePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAutonomousDatabaseBasePropertiesImpl: %+v", err)
	}

	return RawAutonomousDatabaseBasePropertiesImpl{
		autonomousDatabaseBaseProperties: parent,
		Type:                             value,
		Values:                           temp,
	}, nil

}
