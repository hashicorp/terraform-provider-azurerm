package protectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionContainer = MabContainer{}

type MabContainer struct {
	AgentVersion              *string                      `json:"agentVersion,omitempty"`
	CanReRegister             *bool                        `json:"canReRegister,omitempty"`
	ContainerHealthState      *string                      `json:"containerHealthState,omitempty"`
	ContainerId               *int64                       `json:"containerId,omitempty"`
	ExtendedInfo              *MabContainerExtendedInfo    `json:"extendedInfo,omitempty"`
	MabContainerHealthDetails *[]MABContainerHealthDetails `json:"mabContainerHealthDetails,omitempty"`
	ProtectedItemCount        *int64                       `json:"protectedItemCount,omitempty"`

	// Fields inherited from ProtectionContainer

	BackupManagementType  *BackupManagementType    `json:"backupManagementType,omitempty"`
	ContainerType         ProtectableContainerType `json:"containerType"`
	FriendlyName          *string                  `json:"friendlyName,omitempty"`
	HealthStatus          *string                  `json:"healthStatus,omitempty"`
	ProtectableObjectType *string                  `json:"protectableObjectType,omitempty"`
	RegistrationStatus    *string                  `json:"registrationStatus,omitempty"`
}

func (s MabContainer) ProtectionContainer() BaseProtectionContainerImpl {
	return BaseProtectionContainerImpl{
		BackupManagementType:  s.BackupManagementType,
		ContainerType:         s.ContainerType,
		FriendlyName:          s.FriendlyName,
		HealthStatus:          s.HealthStatus,
		ProtectableObjectType: s.ProtectableObjectType,
		RegistrationStatus:    s.RegistrationStatus,
	}
}

var _ json.Marshaler = MabContainer{}

func (s MabContainer) MarshalJSON() ([]byte, error) {
	type wrapper MabContainer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MabContainer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MabContainer: %+v", err)
	}

	decoded["containerType"] = "Windows"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MabContainer: %+v", err)
	}

	return encoded, nil
}
