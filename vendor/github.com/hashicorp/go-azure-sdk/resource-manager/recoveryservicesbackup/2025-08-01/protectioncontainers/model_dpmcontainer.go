package protectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionContainer = DpmContainer{}

type DpmContainer struct {
	CanReRegister      *bool                     `json:"canReRegister,omitempty"`
	ContainerId        *string                   `json:"containerId,omitempty"`
	DpmAgentVersion    *string                   `json:"dpmAgentVersion,omitempty"`
	DpmServers         *[]string                 `json:"dpmServers,omitempty"`
	ExtendedInfo       *DPMContainerExtendedInfo `json:"extendedInfo,omitempty"`
	ProtectedItemCount *int64                    `json:"protectedItemCount,omitempty"`
	ProtectionStatus   *string                   `json:"protectionStatus,omitempty"`
	UpgradeAvailable   *bool                     `json:"upgradeAvailable,omitempty"`

	// Fields inherited from ProtectionContainer

	BackupManagementType  *BackupManagementType    `json:"backupManagementType,omitempty"`
	ContainerType         ProtectableContainerType `json:"containerType"`
	FriendlyName          *string                  `json:"friendlyName,omitempty"`
	HealthStatus          *string                  `json:"healthStatus,omitempty"`
	ProtectableObjectType *string                  `json:"protectableObjectType,omitempty"`
	RegistrationStatus    *string                  `json:"registrationStatus,omitempty"`
}

func (s DpmContainer) ProtectionContainer() BaseProtectionContainerImpl {
	return BaseProtectionContainerImpl{
		BackupManagementType:  s.BackupManagementType,
		ContainerType:         s.ContainerType,
		FriendlyName:          s.FriendlyName,
		HealthStatus:          s.HealthStatus,
		ProtectableObjectType: s.ProtectableObjectType,
		RegistrationStatus:    s.RegistrationStatus,
	}
}

var _ json.Marshaler = DpmContainer{}

func (s DpmContainer) MarshalJSON() ([]byte, error) {
	type wrapper DpmContainer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DpmContainer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DpmContainer: %+v", err)
	}

	decoded["containerType"] = "DPMContainer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DpmContainer: %+v", err)
	}

	return encoded, nil
}
