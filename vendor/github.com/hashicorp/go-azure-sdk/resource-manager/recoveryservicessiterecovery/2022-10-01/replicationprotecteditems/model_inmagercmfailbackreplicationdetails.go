package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = InMageRcmFailbackReplicationDetails{}

type InMageRcmFailbackReplicationDetails struct {
	AzureVirtualMachineId                      *string                                        `json:"azureVirtualMachineId,omitempty"`
	DiscoveredVMDetails                        *InMageRcmFailbackDiscoveredProtectedVMDetails `json:"discoveredVmDetails,omitempty"`
	InitialReplicationProcessedBytes           *int64                                         `json:"initialReplicationProcessedBytes,omitempty"`
	InitialReplicationProgressHealth           *VMReplicationProgressHealth                   `json:"initialReplicationProgressHealth,omitempty"`
	InitialReplicationProgressPercentage       *int64                                         `json:"initialReplicationProgressPercentage,omitempty"`
	InitialReplicationTransferredBytes         *int64                                         `json:"initialReplicationTransferredBytes,omitempty"`
	InternalIdentifier                         *string                                        `json:"internalIdentifier,omitempty"`
	IsAgentRegistrationSuccessfulAfterFailover *bool                                          `json:"isAgentRegistrationSuccessfulAfterFailover,omitempty"`
	LastPlannedFailoverStartTime               *string                                        `json:"lastPlannedFailoverStartTime,omitempty"`
	LastPlannedFailoverStatus                  *PlannedFailoverStatus                         `json:"lastPlannedFailoverStatus,omitempty"`
	LastUsedPolicyFriendlyName                 *string                                        `json:"lastUsedPolicyFriendlyName,omitempty"`
	LastUsedPolicyId                           *string                                        `json:"lastUsedPolicyId,omitempty"`
	LogStorageAccountId                        *string                                        `json:"logStorageAccountId,omitempty"`
	MobilityAgentDetails                       *InMageRcmFailbackMobilityAgentDetails         `json:"mobilityAgentDetails,omitempty"`
	MultiVMGroupName                           *string                                        `json:"multiVmGroupName,omitempty"`
	OsType                                     *string                                        `json:"osType,omitempty"`
	ProtectedDisks                             *[]InMageRcmFailbackProtectedDiskDetails       `json:"protectedDisks,omitempty"`
	ReprotectAgentId                           *string                                        `json:"reprotectAgentId,omitempty"`
	ReprotectAgentName                         *string                                        `json:"reprotectAgentName,omitempty"`
	ResyncProcessedBytes                       *int64                                         `json:"resyncProcessedBytes,omitempty"`
	ResyncProgressHealth                       *VMReplicationProgressHealth                   `json:"resyncProgressHealth,omitempty"`
	ResyncProgressPercentage                   *int64                                         `json:"resyncProgressPercentage,omitempty"`
	ResyncRequired                             *string                                        `json:"resyncRequired,omitempty"`
	ResyncState                                *ResyncState                                   `json:"resyncState,omitempty"`
	ResyncTransferredBytes                     *int64                                         `json:"resyncTransferredBytes,omitempty"`
	TargetDataStoreName                        *string                                        `json:"targetDataStoreName,omitempty"`
	TargetVMName                               *string                                        `json:"targetVmName,omitempty"`
	TargetvCenterId                            *string                                        `json:"targetvCenterId,omitempty"`
	VMNics                                     *[]InMageRcmFailbackNicDetails                 `json:"vmNics,omitempty"`

	// Fields inherited from ReplicationProviderSpecificSettings
}

var _ json.Marshaler = InMageRcmFailbackReplicationDetails{}

func (s InMageRcmFailbackReplicationDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmFailbackReplicationDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmFailbackReplicationDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmFailbackReplicationDetails: %+v", err)
	}
	decoded["instanceType"] = "InMageRcmFailback"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmFailbackReplicationDetails: %+v", err)
	}

	return encoded, nil
}
