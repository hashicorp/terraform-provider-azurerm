package providerinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderSpecificProperties interface {
	ProviderSpecificProperties() BaseProviderSpecificPropertiesImpl
}

var _ ProviderSpecificProperties = BaseProviderSpecificPropertiesImpl{}

type BaseProviderSpecificPropertiesImpl struct {
	ProviderType string `json:"providerType"`
}

func (s BaseProviderSpecificPropertiesImpl) ProviderSpecificProperties() BaseProviderSpecificPropertiesImpl {
	return s
}

var _ ProviderSpecificProperties = RawProviderSpecificPropertiesImpl{}

// RawProviderSpecificPropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawProviderSpecificPropertiesImpl struct {
	providerSpecificProperties BaseProviderSpecificPropertiesImpl
	Type                       string
	Values                     map[string]interface{}
}

func (s RawProviderSpecificPropertiesImpl) ProviderSpecificProperties() BaseProviderSpecificPropertiesImpl {
	return s.providerSpecificProperties
}

func UnmarshalProviderSpecificPropertiesImplementation(input []byte) (ProviderSpecificProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ProviderSpecificProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["providerType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Db2") {
		var out DB2ProviderInstanceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DB2ProviderInstanceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapHana") {
		var out HanaDbProviderInstanceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HanaDbProviderInstanceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MsSqlServer") {
		var out MsSqlServerProviderInstanceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MsSqlServerProviderInstanceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PrometheusHaCluster") {
		var out PrometheusHaClusterProviderInstanceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PrometheusHaClusterProviderInstanceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PrometheusOS") {
		var out PrometheusOSProviderInstanceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PrometheusOSProviderInstanceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapNetWeaver") {
		var out SapNetWeaverProviderInstanceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapNetWeaverProviderInstanceProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseProviderSpecificPropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseProviderSpecificPropertiesImpl: %+v", err)
	}

	return RawProviderSpecificPropertiesImpl{
		providerSpecificProperties: parent,
		Type:                       value,
		Values:                     temp,
	}, nil

}
