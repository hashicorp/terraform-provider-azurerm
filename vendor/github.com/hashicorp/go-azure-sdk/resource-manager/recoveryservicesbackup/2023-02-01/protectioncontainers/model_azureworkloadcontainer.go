package protectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionContainer = AzureWorkloadContainer{}

type AzureWorkloadContainer struct {
	ExtendedInfo     *AzureWorkloadContainerExtendedInfo `json:"extendedInfo,omitempty"`
	LastUpdatedTime  *string                             `json:"lastUpdatedTime,omitempty"`
	OperationType    *OperationType                      `json:"operationType,omitempty"`
	SourceResourceId *string                             `json:"sourceResourceId,omitempty"`
	WorkloadType     *WorkloadType                       `json:"workloadType,omitempty"`

	// Fields inherited from ProtectionContainer

	BackupManagementType  *BackupManagementType    `json:"backupManagementType,omitempty"`
	ContainerType         ProtectableContainerType `json:"containerType"`
	FriendlyName          *string                  `json:"friendlyName,omitempty"`
	HealthStatus          *string                  `json:"healthStatus,omitempty"`
	ProtectableObjectType *string                  `json:"protectableObjectType,omitempty"`
	RegistrationStatus    *string                  `json:"registrationStatus,omitempty"`
}

func (s AzureWorkloadContainer) ProtectionContainer() BaseProtectionContainerImpl {
	return BaseProtectionContainerImpl{
		BackupManagementType:  s.BackupManagementType,
		ContainerType:         s.ContainerType,
		FriendlyName:          s.FriendlyName,
		HealthStatus:          s.HealthStatus,
		ProtectableObjectType: s.ProtectableObjectType,
		RegistrationStatus:    s.RegistrationStatus,
	}
}

var _ json.Marshaler = AzureWorkloadContainer{}

func (s AzureWorkloadContainer) MarshalJSON() ([]byte, error) {
	type wrapper AzureWorkloadContainer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureWorkloadContainer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureWorkloadContainer: %+v", err)
	}

	decoded["containerType"] = "AzureWorkloadContainer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureWorkloadContainer: %+v", err)
	}

	return encoded, nil
}
