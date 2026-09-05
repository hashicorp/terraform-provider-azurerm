package protectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionContainer = GenericContainer{}

type GenericContainer struct {
	ExtendedInformation *GenericContainerExtendedInfo `json:"extendedInformation,omitempty"`
	FabricName          *string                       `json:"fabricName,omitempty"`

	// Fields inherited from ProtectionContainer

	BackupManagementType  *BackupManagementType    `json:"backupManagementType,omitempty"`
	ContainerType         ProtectableContainerType `json:"containerType"`
	FriendlyName          *string                  `json:"friendlyName,omitempty"`
	HealthStatus          *string                  `json:"healthStatus,omitempty"`
	ProtectableObjectType *string                  `json:"protectableObjectType,omitempty"`
	RegistrationStatus    *string                  `json:"registrationStatus,omitempty"`
}

func (s GenericContainer) ProtectionContainer() BaseProtectionContainerImpl {
	return BaseProtectionContainerImpl{
		BackupManagementType:  s.BackupManagementType,
		ContainerType:         s.ContainerType,
		FriendlyName:          s.FriendlyName,
		HealthStatus:          s.HealthStatus,
		ProtectableObjectType: s.ProtectableObjectType,
		RegistrationStatus:    s.RegistrationStatus,
	}
}

var _ json.Marshaler = GenericContainer{}

func (s GenericContainer) MarshalJSON() ([]byte, error) {
	type wrapper GenericContainer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GenericContainer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GenericContainer: %+v", err)
	}

	decoded["containerType"] = "GenericContainer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GenericContainer: %+v", err)
	}

	return encoded, nil
}
