package protecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectedItem interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawProtectedItemImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalProtectedItemImplementation(input []byte) (ProtectedItem, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProtectedItem into map[string]interface: %+v", err)
	}

	value, ok := temp["protectedItemType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureFileShareProtectedItem") {
		var out AzureFileshareProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFileshareProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.ClassicCompute/virtualMachines") {
		var out AzureIaaSClassicComputeVMProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureIaaSClassicComputeVMProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Compute/virtualMachines") {
		var out AzureIaaSComputeVMProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureIaaSComputeVMProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureIaaSVMProtectedItem") {
		var out AzureIaaSVMProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureIaaSVMProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Sql/servers/databases") {
		var out AzureSqlProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureVmWorkloadProtectedItem") {
		var out AzureVMWorkloadProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureVmWorkloadSAPAseDatabase") {
		var out AzureVMWorkloadSAPAseDatabaseProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSAPAseDatabaseProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureVmWorkloadSAPHanaDBInstance") {
		var out AzureVMWorkloadSAPHanaDBInstanceProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSAPHanaDBInstanceProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureVmWorkloadSAPHanaDatabase") {
		var out AzureVMWorkloadSAPHanaDatabaseProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSAPHanaDatabaseProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureVmWorkloadSQLDatabase") {
		var out AzureVMWorkloadSQLDatabaseProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadSQLDatabaseProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DPMProtectedItem") {
		var out DPMProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DPMProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GenericProtectedItem") {
		var out GenericProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GenericProtectedItem: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MabFileFolderProtectedItem") {
		var out MabFileFolderProtectedItem
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MabFileFolderProtectedItem: %+v", err)
		}
		return out, nil
	}

	out := RawProtectedItemImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
