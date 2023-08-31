package backupprotectableitems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkloadProtectableItem interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawWorkloadProtectableItemImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalWorkloadProtectableItemImplementation(input []byte) (WorkloadProtectableItem, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling WorkloadProtectableItem into map[string]interface: %+v", err)
	}

	value, ok := temp["protectableItemType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureFileShare") {
		var out AzureFileShareProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFileShareProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.ClassicCompute/virtualMachines") {
		var out AzureIaaSClassicComputeVMProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureIaaSClassicComputeVMProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Compute/virtualMachines") {
		var out AzureIaaSComputeVMProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureIaaSComputeVMProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureVmWorkloadProtectableItem") {
		var out AzureVMWorkloadProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SAPAseSystem") {
		var out AzureVMWorkloadSAPAseSystemProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSAPAseSystemProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SAPHanaDBInstance") {
		var out AzureVMWorkloadSAPHanaDBInstance
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSAPHanaDBInstance: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SAPHanaDatabase") {
		var out AzureVMWorkloadSAPHanaDatabaseProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSAPHanaDatabaseProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HanaHSRContainer") {
		var out AzureVMWorkloadSAPHanaHSRProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSAPHanaHSRProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SAPHanaSystem") {
		var out AzureVMWorkloadSAPHanaSystemProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSAPHanaSystemProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SQLAvailabilityGroupContainer") {
		var out AzureVMWorkloadSQLAvailabilityGroupProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSQLAvailabilityGroupProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SQLDataBase") {
		var out AzureVMWorkloadSQLDatabaseProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSQLDatabaseProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SQLInstance") {
		var out AzureVMWorkloadSQLInstanceProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSQLInstanceProtectableItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IaaSVMProtectableItem") {
		var out IaaSVMProtectableItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IaaSVMProtectableItem: %+v", err)
		}
		return out, nil
	}

	out := RawWorkloadProtectableItemImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
