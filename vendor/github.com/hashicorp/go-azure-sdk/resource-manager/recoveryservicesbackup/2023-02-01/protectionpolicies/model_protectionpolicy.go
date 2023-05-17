package protectionpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectionPolicy interface {
}

func unmarshalProtectionPolicyImplementation(input []byte) (ProtectionPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProtectionPolicy into map[string]interface: %+v", err)
	}

	value, ok := temp["backupManagementType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureStorage") {
		var out AzureFileShareProtectionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFileShareProtectionPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureIaasVM") {
		var out AzureIaaSVMProtectionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureIaaSVMProtectionPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSql") {
		var out AzureSqlProtectionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlProtectionPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureWorkload") {
		var out AzureVMWorkloadProtectionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureVMWorkloadProtectionPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GenericProtectionPolicy") {
		var out GenericProtectionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GenericProtectionPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MAB") {
		var out MabProtectionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MabProtectionPolicy: %+v", err)
		}
		return out, nil
	}

	type RawProtectionPolicyImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawProtectionPolicyImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
