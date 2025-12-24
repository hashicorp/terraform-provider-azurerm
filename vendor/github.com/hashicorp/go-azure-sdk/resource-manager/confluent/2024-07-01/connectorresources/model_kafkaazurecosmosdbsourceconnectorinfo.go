package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PartnerInfoBase = KafkaAzureCosmosDBSourceConnectorInfo{}

type KafkaAzureCosmosDBSourceConnectorInfo struct {
	ApiKey           *string         `json:"apiKey,omitempty"`
	ApiSecret        *string         `json:"apiSecret,omitempty"`
	AuthType         *AuthType       `json:"authType,omitempty"`
	InputFormat      *DataFormatType `json:"inputFormat,omitempty"`
	MaxTasks         *string         `json:"maxTasks,omitempty"`
	OutputFormat     *DataFormatType `json:"outputFormat,omitempty"`
	ServiceAccountId *string         `json:"serviceAccountId,omitempty"`
	TopicRegex       *string         `json:"topicRegex,omitempty"`
	TopicsDir        *string         `json:"topicsDir,omitempty"`

	// Fields inherited from PartnerInfoBase

	PartnerConnectorType PartnerConnectorType `json:"partnerConnectorType"`
}

func (s KafkaAzureCosmosDBSourceConnectorInfo) PartnerInfoBase() BasePartnerInfoBaseImpl {
	return BasePartnerInfoBaseImpl{
		PartnerConnectorType: s.PartnerConnectorType,
	}
}

var _ json.Marshaler = KafkaAzureCosmosDBSourceConnectorInfo{}

func (s KafkaAzureCosmosDBSourceConnectorInfo) MarshalJSON() ([]byte, error) {
	type wrapper KafkaAzureCosmosDBSourceConnectorInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KafkaAzureCosmosDBSourceConnectorInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KafkaAzureCosmosDBSourceConnectorInfo: %+v", err)
	}

	decoded["partnerConnectorType"] = "KafkaAzureCosmosDBSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KafkaAzureCosmosDBSourceConnectorInfo: %+v", err)
	}

	return encoded, nil
}
