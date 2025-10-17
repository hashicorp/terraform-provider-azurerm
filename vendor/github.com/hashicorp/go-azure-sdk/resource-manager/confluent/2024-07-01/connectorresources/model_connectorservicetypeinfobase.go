package connectorresources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectorServiceTypeInfoBase interface {
	ConnectorServiceTypeInfoBase() BaseConnectorServiceTypeInfoBaseImpl
}

var _ ConnectorServiceTypeInfoBase = BaseConnectorServiceTypeInfoBaseImpl{}

type BaseConnectorServiceTypeInfoBaseImpl struct {
	ConnectorServiceType ConnectorServiceType `json:"connectorServiceType"`
}

func (s BaseConnectorServiceTypeInfoBaseImpl) ConnectorServiceTypeInfoBase() BaseConnectorServiceTypeInfoBaseImpl {
	return s
}

var _ ConnectorServiceTypeInfoBase = RawConnectorServiceTypeInfoBaseImpl{}

// RawConnectorServiceTypeInfoBaseImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawConnectorServiceTypeInfoBaseImpl struct {
	connectorServiceTypeInfoBase BaseConnectorServiceTypeInfoBaseImpl
	Type                         string
	Values                       map[string]interface{}
}

func (s RawConnectorServiceTypeInfoBaseImpl) ConnectorServiceTypeInfoBase() BaseConnectorServiceTypeInfoBaseImpl {
	return s.connectorServiceTypeInfoBase
}

func UnmarshalConnectorServiceTypeInfoBaseImplementation(input []byte) (ConnectorServiceTypeInfoBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ConnectorServiceTypeInfoBase into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["connectorServiceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureBlobStorageSinkConnector") {
		var out AzureBlobStorageSinkConnectorServiceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobStorageSinkConnectorServiceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobStorageSourceConnector") {
		var out AzureBlobStorageSourceConnectorServiceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobStorageSourceConnectorServiceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureCosmosDBSinkConnector") {
		var out AzureCosmosDBSinkConnectorServiceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureCosmosDBSinkConnectorServiceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureCosmosDBSourceConnector") {
		var out AzureCosmosDBSourceConnectorServiceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureCosmosDBSourceConnectorServiceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSynapseAnalyticsSinkConnector") {
		var out AzureSynapseAnalyticsSinkConnectorServiceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSynapseAnalyticsSinkConnectorServiceInfo: %+v", err)
		}
		return out, nil
	}

	var parent BaseConnectorServiceTypeInfoBaseImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseConnectorServiceTypeInfoBaseImpl: %+v", err)
	}

	return RawConnectorServiceTypeInfoBaseImpl{
		connectorServiceTypeInfoBase: parent,
		Type:                         value,
		Values:                       temp,
	}, nil

}
