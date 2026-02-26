package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectorServiceTypeInfoBase = AzureCosmosDBSourceConnectorServiceInfo{}

type AzureCosmosDBSourceConnectorServiceInfo struct {
	CosmosConnectionEndpoint     *string `json:"cosmosConnectionEndpoint,omitempty"`
	CosmosContainersTopicMapping *string `json:"cosmosContainersTopicMapping,omitempty"`
	CosmosDatabaseName           *string `json:"cosmosDatabaseName,omitempty"`
	CosmosMasterKey              *string `json:"cosmosMasterKey,omitempty"`
	CosmosMessageKeyEnabled      *bool   `json:"cosmosMessageKeyEnabled,omitempty"`
	CosmosMessageKeyField        *string `json:"cosmosMessageKeyField,omitempty"`

	// Fields inherited from ConnectorServiceTypeInfoBase

	ConnectorServiceType ConnectorServiceType `json:"connectorServiceType"`
}

func (s AzureCosmosDBSourceConnectorServiceInfo) ConnectorServiceTypeInfoBase() BaseConnectorServiceTypeInfoBaseImpl {
	return BaseConnectorServiceTypeInfoBaseImpl{
		ConnectorServiceType: s.ConnectorServiceType,
	}
}

var _ json.Marshaler = AzureCosmosDBSourceConnectorServiceInfo{}

func (s AzureCosmosDBSourceConnectorServiceInfo) MarshalJSON() ([]byte, error) {
	type wrapper AzureCosmosDBSourceConnectorServiceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureCosmosDBSourceConnectorServiceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureCosmosDBSourceConnectorServiceInfo: %+v", err)
	}

	decoded["connectorServiceType"] = "AzureCosmosDBSourceConnector"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureCosmosDBSourceConnectorServiceInfo: %+v", err)
	}

	return encoded, nil
}
