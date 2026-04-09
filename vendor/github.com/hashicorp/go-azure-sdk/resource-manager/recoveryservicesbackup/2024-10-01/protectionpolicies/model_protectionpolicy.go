package protectionpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectionPolicy interface {
	ProtectionPolicy() BaseProtectionPolicyImpl
}

var _ ProtectionPolicy = BaseProtectionPolicyImpl{}

type BaseProtectionPolicyImpl struct {
	BackupManagementType           string    `json:"backupManagementType"`
	ProtectedItemsCount            *int64    `json:"protectedItemsCount,omitempty"`
	ResourceGuardOperationRequests *[]string `json:"resourceGuardOperationRequests,omitempty"`
}

func (s BaseProtectionPolicyImpl) ProtectionPolicy() BaseProtectionPolicyImpl {
	return s
}

var _ ProtectionPolicy = RawProtectionPolicyImpl{}

// RawProtectionPolicyImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawProtectionPolicyImpl struct {
	protectionPolicy BaseProtectionPolicyImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawProtectionPolicyImpl) ProtectionPolicy() BaseProtectionPolicyImpl {
	return s.protectionPolicy
}

func UnmarshalProtectionPolicyImplementation(input []byte) (ProtectionPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProtectionPolicy into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["backupManagementType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseProtectionPolicyImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseProtectionPolicyImpl: %+v", err)
	}

	return RawProtectionPolicyImpl{
		protectionPolicy: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
