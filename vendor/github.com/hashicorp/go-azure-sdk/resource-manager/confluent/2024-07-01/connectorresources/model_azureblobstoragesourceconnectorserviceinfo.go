package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectorServiceTypeInfoBase = AzureBlobStorageSourceConnectorServiceInfo{}

type AzureBlobStorageSourceConnectorServiceInfo struct {
	StorageAccountKey    *string `json:"storageAccountKey,omitempty"`
	StorageAccountName   *string `json:"storageAccountName,omitempty"`
	StorageContainerName *string `json:"storageContainerName,omitempty"`

	// Fields inherited from ConnectorServiceTypeInfoBase

	ConnectorServiceType ConnectorServiceType `json:"connectorServiceType"`
}

func (s AzureBlobStorageSourceConnectorServiceInfo) ConnectorServiceTypeInfoBase() BaseConnectorServiceTypeInfoBaseImpl {
	return BaseConnectorServiceTypeInfoBaseImpl{
		ConnectorServiceType: s.ConnectorServiceType,
	}
}

var _ json.Marshaler = AzureBlobStorageSourceConnectorServiceInfo{}

func (s AzureBlobStorageSourceConnectorServiceInfo) MarshalJSON() ([]byte, error) {
	type wrapper AzureBlobStorageSourceConnectorServiceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureBlobStorageSourceConnectorServiceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBlobStorageSourceConnectorServiceInfo: %+v", err)
	}

	decoded["connectorServiceType"] = "AzureBlobStorageSourceConnector"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureBlobStorageSourceConnectorServiceInfo: %+v", err)
	}

	return encoded, nil
}
