package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnableProtectionProviderSpecificInput interface {
	EnableProtectionProviderSpecificInput() BaseEnableProtectionProviderSpecificInputImpl
}

var _ EnableProtectionProviderSpecificInput = BaseEnableProtectionProviderSpecificInputImpl{}

type BaseEnableProtectionProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseEnableProtectionProviderSpecificInputImpl) EnableProtectionProviderSpecificInput() BaseEnableProtectionProviderSpecificInputImpl {
	return s
}

var _ EnableProtectionProviderSpecificInput = RawEnableProtectionProviderSpecificInputImpl{}

// RawEnableProtectionProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEnableProtectionProviderSpecificInputImpl struct {
	enableProtectionProviderSpecificInput BaseEnableProtectionProviderSpecificInputImpl
	Type                                  string
	Values                                map[string]interface{}
}

func (s RawEnableProtectionProviderSpecificInputImpl) EnableProtectionProviderSpecificInput() BaseEnableProtectionProviderSpecificInputImpl {
	return s.enableProtectionProviderSpecificInput
}

func UnmarshalEnableProtectionProviderSpecificInputImplementation(input []byte) (EnableProtectionProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EnableProtectionProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2ACrossClusterMigration") {
		var out A2ACrossClusterMigrationEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2ACrossClusterMigrationEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzureEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzureEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2EnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2EnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMage") {
		var out InMageEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseEnableProtectionProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEnableProtectionProviderSpecificInputImpl: %+v", err)
	}

	return RawEnableProtectionProviderSpecificInputImpl{
		enableProtectionProviderSpecificInput: parent,
		Type:                                  value,
		Values:                                temp,
	}, nil

}
