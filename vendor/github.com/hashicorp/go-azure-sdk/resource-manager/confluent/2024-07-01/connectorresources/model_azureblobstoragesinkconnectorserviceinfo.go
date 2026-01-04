package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectorServiceTypeInfoBase = AzureBlobStorageSinkConnectorServiceInfo{}

type AzureBlobStorageSinkConnectorServiceInfo struct {
	StorageAccountKey    *string `json:"storageAccountKey,omitempty"`
	StorageAccountName   *string `json:"storageAccountName,omitempty"`
	StorageContainerName *string `json:"storageContainerName,omitempty"`

	// Fields inherited from ConnectorServiceTypeInfoBase

	ConnectorServiceType ConnectorServiceType `json:"connectorServiceType"`
}

func (s AzureBlobStorageSinkConnectorServiceInfo) ConnectorServiceTypeInfoBase() BaseConnectorServiceTypeInfoBaseImpl {
	return BaseConnectorServiceTypeInfoBaseImpl{
		ConnectorServiceType: s.ConnectorServiceType,
	}
}

var _ json.Marshaler = AzureBlobStorageSinkConnectorServiceInfo{}

func (s AzureBlobStorageSinkConnectorServiceInfo) MarshalJSON() ([]byte, error) {
	type wrapper AzureBlobStorageSinkConnectorServiceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureBlobStorageSinkConnectorServiceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBlobStorageSinkConnectorServiceInfo: %+v", err)
	}

	decoded["connectorServiceType"] = "AzureBlobStorageSinkConnector"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureBlobStorageSinkConnectorServiceInfo: %+v", err)
	}

	return encoded, nil
}
