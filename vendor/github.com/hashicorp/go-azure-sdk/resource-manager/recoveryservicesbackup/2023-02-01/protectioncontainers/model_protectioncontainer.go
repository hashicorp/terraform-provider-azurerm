package protectioncontainers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectionContainer interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawProtectionContainerImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalProtectionContainerImplementation(input []byte) (ProtectionContainer, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProtectionContainer into map[string]interface: %+v", err)
	}

	value, ok := temp["containerType"].(string)
	if !ok {
		return nil, nil
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

	out := RawProtectionContainerImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
