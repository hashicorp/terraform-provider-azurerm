package protectioncontainers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectionContainer interface {
	ProtectionContainer() BaseProtectionContainerImpl
}

var _ ProtectionContainer = BaseProtectionContainerImpl{}

type BaseProtectionContainerImpl struct {
	BackupManagementType  *BackupManagementType    `json:"backupManagementType,omitempty"`
	ContainerType         ProtectableContainerType `json:"containerType"`
	FriendlyName          *string                  `json:"friendlyName,omitempty"`
	HealthStatus          *string                  `json:"healthStatus,omitempty"`
	ProtectableObjectType *string                  `json:"protectableObjectType,omitempty"`
	RegistrationStatus    *string                  `json:"registrationStatus,omitempty"`
}

func (s BaseProtectionContainerImpl) ProtectionContainer() BaseProtectionContainerImpl {
	return s
}

var _ ProtectionContainer = RawProtectionContainerImpl{}

// RawProtectionContainerImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawProtectionContainerImpl struct {
	protectionContainer BaseProtectionContainerImpl
	Type                string
	Values              map[string]interface{}
}

func (s RawProtectionContainerImpl) ProtectionContainer() BaseProtectionContainerImpl {
	return s.protectionContainer
}

func (s RawProtectionContainerImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalProtectionContainerImplementation(input []byte) (ProtectionContainer, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProtectionContainer into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["containerType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureBackupServerContainer") {
		var out AzureBackupServerContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupServerContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.ClassicCompute/virtualMachines") {
		var out AzureIaaSClassicComputeVMContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureIaaSClassicComputeVMContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Compute/virtualMachines") {
		var out AzureIaaSComputeVMContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureIaaSComputeVMContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SQLAGWorkLoadContainer") {
		var out AzureSQLAGWorkloadContainerProtectionContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSQLAGWorkloadContainerProtectionContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSqlContainer") {
		var out AzureSqlContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StorageContainer") {
		var out AzureStorageContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureStorageContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMAppContainer") {
		var out AzureVMAppContainerProtectionContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMAppContainerProtectionContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureWorkloadContainer") {
		var out AzureWorkloadContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureWorkloadContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DPMContainer") {
		var out DpmContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DpmContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GenericContainer") {
		var out GenericContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GenericContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IaasVMContainer") {
		var out IaaSVMContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IaaSVMContainer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Windows") {
		var out MabContainer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MabContainer: %+v", err)
		}
		return out, nil
	}

	var parent BaseProtectionContainerImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseProtectionContainerImpl: %+v", err)
	}

	return RawProtectionContainerImpl{
		protectionContainer: parent,
		Type:                value,
		Values:              temp,
	}, nil

}
