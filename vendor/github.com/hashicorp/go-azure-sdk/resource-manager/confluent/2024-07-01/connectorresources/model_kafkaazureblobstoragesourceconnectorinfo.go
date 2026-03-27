package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PartnerInfoBase = KafkaAzureBlobStorageSourceConnectorInfo{}

type KafkaAzureBlobStorageSourceConnectorInfo struct {
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

func (s KafkaAzureBlobStorageSourceConnectorInfo) PartnerInfoBase() BasePartnerInfoBaseImpl {
	return BasePartnerInfoBaseImpl{
		PartnerConnectorType: s.PartnerConnectorType,
	}
}

var _ json.Marshaler = KafkaAzureBlobStorageSourceConnectorInfo{}

func (s KafkaAzureBlobStorageSourceConnectorInfo) MarshalJSON() ([]byte, error) {
	type wrapper KafkaAzureBlobStorageSourceConnectorInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KafkaAzureBlobStorageSourceConnectorInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KafkaAzureBlobStorageSourceConnectorInfo: %+v", err)
	}

	decoded["partnerConnectorType"] = "KafkaAzureBlobStorageSource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KafkaAzureBlobStorageSourceConnectorInfo: %+v", err)
	}

	return encoded, nil
}
