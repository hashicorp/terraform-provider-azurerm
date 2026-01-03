package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PartnerInfoBase = KafkaAzureSynapseAnalyticsSinkConnectorInfo{}

type KafkaAzureSynapseAnalyticsSinkConnectorInfo struct {
	ApiKey           *string         `json:"apiKey,omitempty"`
	ApiSecret        *string         `json:"apiSecret,omitempty"`
	AuthType         *AuthType       `json:"authType,omitempty"`
	FlushSize        *string         `json:"flushSize,omitempty"`
	InputFormat      *DataFormatType `json:"inputFormat,omitempty"`
	MaxTasks         *string         `json:"maxTasks,omitempty"`
	OutputFormat     *DataFormatType `json:"outputFormat,omitempty"`
	ServiceAccountId *string         `json:"serviceAccountId,omitempty"`
	TimeInterval     *string         `json:"timeInterval,omitempty"`
	Topics           *[]string       `json:"topics,omitempty"`
	TopicsDir        *string         `json:"topicsDir,omitempty"`

	// Fields inherited from PartnerInfoBase

	PartnerConnectorType PartnerConnectorType `json:"partnerConnectorType"`
}

func (s KafkaAzureSynapseAnalyticsSinkConnectorInfo) PartnerInfoBase() BasePartnerInfoBaseImpl {
	return BasePartnerInfoBaseImpl{
		PartnerConnectorType: s.PartnerConnectorType,
	}
}

var _ json.Marshaler = KafkaAzureSynapseAnalyticsSinkConnectorInfo{}

func (s KafkaAzureSynapseAnalyticsSinkConnectorInfo) MarshalJSON() ([]byte, error) {
	type wrapper KafkaAzureSynapseAnalyticsSinkConnectorInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KafkaAzureSynapseAnalyticsSinkConnectorInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KafkaAzureSynapseAnalyticsSinkConnectorInfo: %+v", err)
	}

	decoded["partnerConnectorType"] = "KafkaAzureSynapseAnalyticsSink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KafkaAzureSynapseAnalyticsSinkConnectorInfo: %+v", err)
	}

	return encoded, nil
}
