package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PartnerInfoBase = KafkaAzureBlobStorageSinkConnectorInfo{}

type KafkaAzureBlobStorageSinkConnectorInfo struct {
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

func (s KafkaAzureBlobStorageSinkConnectorInfo) PartnerInfoBase() BasePartnerInfoBaseImpl {
	return BasePartnerInfoBaseImpl{
		PartnerConnectorType: s.PartnerConnectorType,
	}
}

var _ json.Marshaler = KafkaAzureBlobStorageSinkConnectorInfo{}

func (s KafkaAzureBlobStorageSinkConnectorInfo) MarshalJSON() ([]byte, error) {
	type wrapper KafkaAzureBlobStorageSinkConnectorInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KafkaAzureBlobStorageSinkConnectorInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KafkaAzureBlobStorageSinkConnectorInfo: %+v", err)
	}

	decoded["partnerConnectorType"] = "KafkaAzureBlobStorageSink"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KafkaAzureBlobStorageSinkConnectorInfo: %+v", err)
	}

	return encoded, nil
}
