package connectorresources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerInfoBase interface {
	PartnerInfoBase() BasePartnerInfoBaseImpl
}

var _ PartnerInfoBase = BasePartnerInfoBaseImpl{}

type BasePartnerInfoBaseImpl struct {
	PartnerConnectorType PartnerConnectorType `json:"partnerConnectorType"`
}

func (s BasePartnerInfoBaseImpl) PartnerInfoBase() BasePartnerInfoBaseImpl {
	return s
}

var _ PartnerInfoBase = RawPartnerInfoBaseImpl{}

// RawPartnerInfoBaseImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawPartnerInfoBaseImpl struct {
	partnerInfoBase BasePartnerInfoBaseImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawPartnerInfoBaseImpl) PartnerInfoBase() BasePartnerInfoBaseImpl {
	return s.partnerInfoBase
}

func UnmarshalPartnerInfoBaseImplementation(input []byte) (PartnerInfoBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling PartnerInfoBase into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["partnerConnectorType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "KafkaAzureBlobStorageSink") {
		var out KafkaAzureBlobStorageSinkConnectorInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KafkaAzureBlobStorageSinkConnectorInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KafkaAzureBlobStorageSource") {
		var out KafkaAzureBlobStorageSourceConnectorInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KafkaAzureBlobStorageSourceConnectorInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KafkaAzureCosmosDBSink") {
		var out KafkaAzureCosmosDBSinkConnectorInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KafkaAzureCosmosDBSinkConnectorInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KafkaAzureCosmosDBSource") {
		var out KafkaAzureCosmosDBSourceConnectorInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KafkaAzureCosmosDBSourceConnectorInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KafkaAzureSynapseAnalyticsSink") {
		var out KafkaAzureSynapseAnalyticsSinkConnectorInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KafkaAzureSynapseAnalyticsSinkConnectorInfo: %+v", err)
		}
		return out, nil
	}

	var parent BasePartnerInfoBaseImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BasePartnerInfoBaseImpl: %+v", err)
	}

	return RawPartnerInfoBaseImpl{
		partnerInfoBase: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
