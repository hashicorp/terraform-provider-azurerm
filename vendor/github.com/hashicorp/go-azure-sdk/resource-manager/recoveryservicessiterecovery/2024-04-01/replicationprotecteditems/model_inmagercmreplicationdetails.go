package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = InMageRcmReplicationDetails{}

type InMageRcmReplicationDetails struct {
	AgentUpgradeAttemptToVersion               *string                                      `json:"agentUpgradeAttemptToVersion,omitempty"`
	AgentUpgradeBlockingErrorDetails           *[]InMageRcmAgentUpgradeBlockingErrorDetails `json:"agentUpgradeBlockingErrorDetails,omitempty"`
	AgentUpgradeJobId                          *string                                      `json:"agentUpgradeJobId,omitempty"`
	AgentUpgradeState                          *MobilityAgentUpgradeState                   `json:"agentUpgradeState,omitempty"`
	AllocatedMemoryInMB                        *float64                                     `json:"allocatedMemoryInMB,omitempty"`
	DiscoveredVMDetails                        *InMageRcmDiscoveredProtectedVMDetails       `json:"discoveredVmDetails,omitempty"`
	DiscoveryType                              *string                                      `json:"discoveryType,omitempty"`
	FabricDiscoveryMachineId                   *string                                      `json:"fabricDiscoveryMachineId,omitempty"`
	FailoverRecoveryPointId                    *string                                      `json:"failoverRecoveryPointId,omitempty"`
	FirmwareType                               *string                                      `json:"firmwareType,omitempty"`
	InitialReplicationProcessedBytes           *int64                                       `json:"initialReplicationProcessedBytes,omitempty"`
	InitialReplicationProgressHealth           *VMReplicationProgressHealth                 `json:"initialReplicationProgressHealth,omitempty"`
	InitialReplicationProgressPercentage       *int64                                       `json:"initialReplicationProgressPercentage,omitempty"`
	InitialReplicationTransferredBytes         *int64                                       `json:"initialReplicationTransferredBytes,omitempty"`
	InternalIdentifier                         *string                                      `json:"internalIdentifier,omitempty"`
	IsAgentRegistrationSuccessfulAfterFailover *bool                                        `json:"isAgentRegistrationSuccessfulAfterFailover,omitempty"`
	IsLastUpgradeSuccessful                    *string                                      `json:"isLastUpgradeSuccessful,omitempty"`
	LastAgentUpgradeErrorDetails               *[]InMageRcmLastAgentUpgradeErrorDetails     `json:"lastAgentUpgradeErrorDetails,omitempty"`
	LastAgentUpgradeType                       *string                                      `json:"lastAgentUpgradeType,omitempty"`
	LastRecoveryPointId                        *string                                      `json:"lastRecoveryPointId,omitempty"`
	LastRecoveryPointReceived                  *string                                      `json:"lastRecoveryPointReceived,omitempty"`
	LastRpoCalculatedTime                      *string                                      `json:"lastRpoCalculatedTime,omitempty"`
	LastRpoInSeconds                           *int64                                       `json:"lastRpoInSeconds,omitempty"`
	LicenseType                                *string                                      `json:"licenseType,omitempty"`
	MobilityAgentDetails                       *InMageRcmMobilityAgentDetails               `json:"mobilityAgentDetails,omitempty"`
	MultiVMGroupName                           *string                                      `json:"multiVmGroupName,omitempty"`
	OsName                                     *string                                      `json:"osName,omitempty"`
	OsType                                     *string                                      `json:"osType,omitempty"`
	PrimaryNicIPAddress                        *string                                      `json:"primaryNicIpAddress,omitempty"`
	ProcessServerId                            *string                                      `json:"processServerId,omitempty"`
	ProcessServerName                          *string                                      `json:"processServerName,omitempty"`
	ProcessorCoreCount                         *int64                                       `json:"processorCoreCount,omitempty"`
	ProtectedDisks                             *[]InMageRcmProtectedDiskDetails             `json:"protectedDisks,omitempty"`
	ResyncProcessedBytes                       *int64                                       `json:"resyncProcessedBytes,omitempty"`
	ResyncProgressHealth                       *VMReplicationProgressHealth                 `json:"resyncProgressHealth,omitempty"`
	ResyncProgressPercentage                   *int64                                       `json:"resyncProgressPercentage,omitempty"`
	ResyncRequired                             *string                                      `json:"resyncRequired,omitempty"`
	ResyncState                                *ResyncState                                 `json:"resyncState,omitempty"`
	ResyncTransferredBytes                     *int64                                       `json:"resyncTransferredBytes,omitempty"`
	RunAsAccountId                             *string                                      `json:"runAsAccountId,omitempty"`
	SeedManagedDiskTags                        *[]UserCreatedResourceTag                    `json:"seedManagedDiskTags,omitempty"`
	SqlServerLicenseType                       *string                                      `json:"sqlServerLicenseType,omitempty"`
	StorageAccountId                           *string                                      `json:"storageAccountId,omitempty"`
	SupportedOSVersions                        *[]string                                    `json:"supportedOSVersions,omitempty"`
	TargetAvailabilitySetId                    *string                                      `json:"targetAvailabilitySetId,omitempty"`
	TargetAvailabilityZone                     *string                                      `json:"targetAvailabilityZone,omitempty"`
	TargetBootDiagnosticsStorageAccountId      *string                                      `json:"targetBootDiagnosticsStorageAccountId,omitempty"`
	TargetGeneration                           *string                                      `json:"targetGeneration,omitempty"`
	TargetLocation                             *string                                      `json:"targetLocation,omitempty"`
	TargetManagedDiskTags                      *[]UserCreatedResourceTag                    `json:"targetManagedDiskTags,omitempty"`
	TargetNetworkId                            *string                                      `json:"targetNetworkId,omitempty"`
	TargetNicTags                              *[]UserCreatedResourceTag                    `json:"targetNicTags,omitempty"`
	TargetProximityPlacementGroupId            *string                                      `json:"targetProximityPlacementGroupId,omitempty"`
	TargetResourceGroupId                      *string                                      `json:"targetResourceGroupId,omitempty"`
	TargetVMName                               *string                                      `json:"targetVmName,omitempty"`
	TargetVMSecurityProfile                    *SecurityProfileProperties                   `json:"targetVmSecurityProfile,omitempty"`
	TargetVMSize                               *string                                      `json:"targetVmSize,omitempty"`
	TargetVMTags                               *[]UserCreatedResourceTag                    `json:"targetVmTags,omitempty"`
	TestNetworkId                              *string                                      `json:"testNetworkId,omitempty"`
	UnprotectedDisks                           *[]InMageRcmUnProtectedDiskDetails           `json:"unprotectedDisks,omitempty"`
	VMNics                                     *[]InMageRcmNicDetails                       `json:"vmNics,omitempty"`

	// Fields inherited from ReplicationProviderSpecificSettings

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmReplicationDetails) ReplicationProviderSpecificSettings() BaseReplicationProviderSpecificSettingsImpl {
	return BaseReplicationProviderSpecificSettingsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmReplicationDetails{}

func (s InMageRcmReplicationDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmReplicationDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmReplicationDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmReplicationDetails: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmReplicationDetails: %+v", err)
	}

	return encoded, nil
}
